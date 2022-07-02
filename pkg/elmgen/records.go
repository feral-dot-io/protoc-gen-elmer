package elmgen

import (
	"log"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addRecords(msgs []*protogen.Message) {
	for _, msg := range msgs {
		// Defer map handling to fields
		if msg.Desc.IsMapEntry() {
			continue
		}

		m.Records = append(m.Records, m.newRecord(msg))
		// Add nested
		m.addUnions(msg.Enums)
		m.addRecords(msg.Messages)
	}
	sort.Sort(m.Records)
	sort.Sort(m.Oneofs)
}

func (m *Module) newRecord(msg *protogen.Message) *Record {
	md := msg.Desc
	var record Record
	record.Type = m.NewElmType(md.ParentFile(), md)
	record.Comments = NewCommentSet(msg.Comments)
	oneofsSeen := make(map[protoreflect.FullName]bool)

	for _, field := range msg.Fields {
		// Oneof field?
		if field.Oneof != nil {
			od := field.Oneof.Desc
			// Only add once
			if oneofsSeen[od.FullName()] {
				continue
			}
			oneofsSeen[od.FullName()] = true

			oneof, field := m.newOneofField(field.Oneof)
			m.Oneofs = append(m.Oneofs, oneof)
			record.Oneofs = append(record.Oneofs, oneof)
			record.Fields = append(record.Fields, field)
		} else { // Regular field
			field := m.newField(field)
			record.Fields = append(record.Fields, field)
		}
	}
	return &record
}

func (m *Module) newField(field *protogen.Field) *Field {
	fd := field.Desc
	var typ, zero, decoder, encoder, fuzzer string
	typ = fieldType(m, fd)
	zero = fieldZero(m, fd)
	decoder = fieldDecoder(m, fd)
	encoder = fieldEncoder(m, fd)
	fuzzer = fieldFuzzer(m, fd)
	var key *MapKey
	if fd.IsMap() {
		pdKey := fd.MapKey()
		key = new(MapKey)
		key.Zero = fieldZero(m, pdKey)
		key.Decoder = fieldDecoder(m, pdKey)
		key.Encoder = fieldEncoder(m, pdKey)
		key.Fuzzer = fieldFuzzer(m, pdKey)
	}
	return &Field{
		protoFullIdentToElmCasing(string(fd.Name()), "", false),
		NewCommentSet(field.Comments),
		false, fd.IsMap(), fd.Number(), fd.Cardinality(),
		typ, zero, decoder, encoder, fuzzer, key,
	}
}

func (m *Module) newOneofField(protoOneof *protogen.Oneof) (*Oneof, *Field) {
	od := protoOneof.Desc
	oneof := m.newOneof(protoOneof)
	field := &Field{
		protoFullIdentToElmCasing(string(od.Name()), "", false),
		NewCommentSet(protoOneof.Comments),
		true, false, 0, 0,
		"(Maybe " + oneof.Type.String() + ")",
		"Nothing",
		oneof.Type.Decoder().String(),
		oneof.Type.Encoder().String(),
		oneof.Type.Fuzzer().String(),
		nil}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Label = string(od.Fields().Get(0).Name())
		field.Type = "(Maybe " + oneof.Variants[0].Field.Type + ")"
	}
	return oneof, field
}

func fieldType(m *Module, pd protoreflect.FieldDescriptor) string {
	if pd.IsMap() {
		key := fieldType(m, pd.MapKey())
		val := fieldType(m, pd.MapValue())
		m.Imports.Dict = true
		return "(Dict " + key + " " + val + ")"
	}

	if pd.IsList() {
		val := fieldTypeFromKind(m, pd)
		return "(List " + val + ")"
	}

	return fieldTypeFromKind(m, pd)
}

func fieldTypeFromKind(m *Module, pd protoreflect.FieldDescriptor) string {
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return "Bool"

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind:
		return "Int"

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "Float"

	case protoreflect.StringKind:
		return "String"

	case protoreflect.BytesKind:
		m.Imports.Bytes = true
		return "Bytes"

	case protoreflect.EnumKind:
		ed := pd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return m.NewElmType(md.ParentFile(), md).String()
	}

	log.Panicf("fieldType: unknown protoreflect.Kind: %s", pd.Kind())
	return ""
}

func fieldZero(m *Module, pd protoreflect.FieldDescriptor) string {
	if pd.IsMap() {
		pd = pd.MapValue()
	}
	if pd.IsList() {
		return "[]"
	}

	switch pd.Kind() {
	// No formatting difference for these fields
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return pd.Default().String()

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.BoolKind:
		if pd.Default().Bool() {
			return "True"
		}
		return "False"

	case protoreflect.StringKind:
		return `""`

	case protoreflect.BytesKind:
		return "(BE.encode (BE.sequence []))"

	case protoreflect.EnumKind:
		ed := pd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).Zero().String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return m.NewElmType(md.ParentFile(), md).Zero().String()
	}

	log.Panicf("fieldZero: unknown protoreflect.Kind: %s", pd.Kind())
	return ""
}

func fieldDecoder(m *Module, pd protoreflect.FieldDescriptor) string {
	return fieldKindCodec(m, "PD.", "Decoder", pd)
}

func fieldEncoder(m *Module, pd protoreflect.FieldDescriptor) string {
	return fieldKindCodec(m, "PE.", "Encoder", pd)
}

// Just the Kind. Does not take into account special features like lists.
func fieldKindCodec(m *Module, lib, dir string, pd protoreflect.FieldDescriptor) string {
	if pd.IsMap() {
		pd = pd.MapValue()
	}
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return lib + "bool"

	case protoreflect.Int32Kind:
		return lib + "int32"
	case protoreflect.Sint32Kind:
		return lib + "sint32"
	case protoreflect.Uint32Kind:
		return lib + "uint32"
	case protoreflect.Sfixed32Kind:
		return lib + "sfixed32"
	case protoreflect.Fixed32Kind:
		return lib + "fixed32"

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind:
		return lib + "float"
	case protoreflect.DoubleKind:
		return lib + "double"

	case protoreflect.StringKind:
		return lib + "string"

	case protoreflect.BytesKind:
		return lib + "bytes"

	case protoreflect.EnumKind:
		ed := pd.Enum()
		return m.NewElmValue(ed.ParentFile(), ed).String() + dir

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return m.NewElmValue(md.ParentFile(), md).String() + dir
	}

	log.Panicf("fieldCodec: unknown protoreflect.Kind: %s", pd.Kind())
	return ""
}

func fieldFuzzer(m *Module, pd protoreflect.FieldDescriptor) string {
	if pd.IsMap() {
		pd = pd.MapValue()
	}
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return "Fuzz.bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		m.Fuzzers.Int32 = true
		return "fuzzInt32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		m.Fuzzers.Uint32 = true
		return "fuzzUint32"

	/* Unsupported by Elm / JS
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64", nil
	*/

	case protoreflect.FloatKind:
		m.Fuzzers.Float32 = true
		return "fuzzFloat32"
	case protoreflect.DoubleKind:
		return "Fuzz.float"

	case protoreflect.StringKind:
		return "Fuzz.string"
	case protoreflect.BytesKind:
		m.Imports.Bytes = true
		return "fuzzBytes"

	case protoreflect.EnumKind:
		ed := pd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).Fuzzer().String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return m.NewElmType(md.ParentFile(), md).Fuzzer().String()
	}

	log.Panicf("kindFuzzer: unknown protoreflect.Kind: %s", pd.Kind())
	return ""
}
