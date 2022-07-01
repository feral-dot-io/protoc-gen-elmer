package elmgen

import (
	"log"
	"sort"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addRecords(msgs protoreflect.MessageDescriptors) {
	for i := 0; i < msgs.Len(); i++ {
		md := msgs.Get(i)

		// Defer map handling to fields
		if md.IsMapEntry() {
			continue
		}

		m.Records = append(m.Records, m.newRecord(md))
		// Add nested
		m.addUnions(md.Enums())
		m.addRecords(md.Messages())
	}
	sort.Sort(m.Records)
	sort.Sort(m.Oneofs)
}

func (m *Module) newRecord(md protoreflect.MessageDescriptor) *Record {
	var record Record
	record.Type = NewElmType(md.ParentFile(), md)
	oneofsSeen := make(map[protoreflect.FullName]bool)

	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		// Oneof field?
		if od := fd.ContainingOneof(); od != nil {
			// Only add once
			if oneofsSeen[od.FullName()] {
				continue
			}
			oneofsSeen[od.FullName()] = true

			oneof, field := m.newOneofField(od)
			m.Oneofs = append(m.Oneofs, oneof)
			record.Oneofs = append(record.Oneofs, oneof)
			record.Fields = append(record.Fields, field)
		} else { // Regular field
			field := m.newField(fd)
			record.Fields = append(record.Fields, field)
		}
	}
	return &record
}

func (m *Module) newField(pd protoreflect.FieldDescriptor) *Field {
	var typ, zero, decoder, encoder, fuzzer string
	typ = fieldType(m, pd)
	zero = fieldZero(m, pd)
	decoder = fieldDecoder(m, pd)
	encoder = fieldEncoder(m, pd)
	fuzzer = fieldFuzzer(m, pd)
	var key *MapKey
	if pd.IsMap() {
		pdKey := pd.MapKey()
		key = new(MapKey)
		key.Zero = fieldZero(m, pdKey)
		key.Decoder = fieldDecoder(m, pdKey)
		key.Encoder = fieldEncoder(m, pdKey)
		key.Fuzzer = fieldFuzzer(m, pdKey)
	}
	return &Field{
		protoFullIdentToElmCasing(string(pd.Name()), "", false),
		false, pd.IsMap(), pd.Number(), pd.Cardinality(),
		typ, zero, decoder, encoder, fuzzer, key,
	}
}

func (m *Module) newOneofField(po protoreflect.OneofDescriptor) (*Oneof, *Field) {
	oneof := m.newOneof(po)
	field := &Field{
		protoFullIdentToElmCasing(string(po.Name()), "", false),
		true, false, 0, 0,
		"(Maybe " + oneof.Type.Local() + ")",
		"Nothing",
		oneof.Type.Decoder().Local(),
		oneof.Type.Encoder().Local(),
		oneof.Type.Fuzzer().Local(),
		nil}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Label = string(po.Fields().Get(0).Name())
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
		return NewElmType(ed.ParentFile(), ed).Local()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return NewElmType(md.ParentFile(), md).Local()
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
		if pd.Default().Bool() { // TODO proto2 test
			return "True"
		}
		return "False"

	case protoreflect.StringKind:
		return `""`

	case protoreflect.BytesKind:
		return "(BE.encode (BE.sequence []))"

	case protoreflect.EnumKind:
		ed := pd.Enum()
		return NewElmType(ed.ParentFile(), ed).Zero().Local()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return NewElmType(md.ParentFile(), md).Zero().Local()
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
		return NewElmValue(ed.ParentFile(), ed).Local() + dir

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return NewElmValue(md.ParentFile(), md).Local() + dir
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
		return NewElmType(ed.ParentFile(), ed).Fuzzer().Local()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := pd.Message()
		return NewElmType(md.ParentFile(), md).Fuzzer().Local()
	}

	log.Panicf("kindFuzzer: unknown protoreflect.Kind: %s", pd.Kind())
	return ""
}
