package elmgen

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addRecords() error {
	for _, proto := range m.protoMessages {
		// Defer map handling to fields
		if proto.Desc.IsMapEntry() {
			continue
		}

		record, err := m.newRecord(proto)
		if err != nil {
			return err
		}
		m.Records = append(m.Records, record)
	}
	sort.Sort(m.Records)
	sort.Sort(m.Oneofs)
	return nil
}

func (m *Module) newRecord(proto *protogen.Message) (*Record, error) {
	var record Record
	err := record.CodecIDs.register(m, proto.Desc.FullName())
	if err != nil {
		return nil, err
	}
	oneofsSeen := make(map[protoreflect.FullName]bool)
	for _, proto := range proto.Fields {
		// Part of a oneof field?
		if po := proto.Oneof; po != nil {
			// Only add once
			if oneofsSeen[po.Desc.FullName()] {
				continue
			}
			oneofsSeen[po.Desc.FullName()] = true

			oneof, field, err := m.newOneofField(po.Desc)
			if err != nil {
				return nil, err
			}
			m.Oneofs = append(m.Oneofs, oneof)
			record.Oneofs = append(record.Oneofs, oneof)
			record.Fields = append(record.Fields, field)
		} else {
			// Regular field
			field, err := m.newField(proto.Desc)
			if err != nil {
				return nil, err
			}
			record.Fields = append(record.Fields, field)
		}
	}
	return &record, nil
}

func (m *Module) newField(pd protoreflect.FieldDescriptor) (*Field, error) {
	var typ, decoder, encoder, fuzzer string
	var zero interface{}
	var err error
	if typ, err = fieldType(m, pd); err != nil {
		return nil, err
	}
	if zero, err = fieldZero(m, pd); err != nil {
		return nil, err
	}
	if decoder, err = fieldDecoder(m, pd); err != nil {
		return nil, err
	}
	if encoder, err = fieldEncoder(m, pd); err != nil {
		return nil, err
	}
	if fuzzer, err = fieldFuzzer(m, pd); err != nil {
		return nil, err
	}
	return &Field{
		m.getElmValue(protoreflect.FullName(pd.Name())),
		false, pd.Number(), pd.Cardinality(),
		typ, zero, decoder, encoder, fuzzer,
	}, nil
}

func (m *Module) newOneofField(po protoreflect.OneofDescriptor) (*Oneof, *Field, error) {
	oneof, err := m.newOneof(po)
	if err != nil {
		return nil, nil, err
	}
	field := &Field{
		m.getElmValue(protoreflect.FullName(po.Name())),
		true, 0, 0,
		"(Maybe " + string(oneof.ID) + ")",
		"Nothing",
		oneof.DecodeID,
		oneof.EncodeID,
		oneof.FuzzerID}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Label = string(po.Fields().Get(0).Name())
		field.Type = "(Maybe " + oneof.Variants[0].Field.Type + ")"
	}
	return oneof, field, nil
}

func fieldType(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	if pd.IsMap() {
		key, err := fieldType(m, pd.MapKey())
		if err != nil {
			return "", err
		}
		val, err := fieldType(m, pd.MapValue())
		m.Imports.Dict = true
		return "(Dict " + key + " " + val + ")", err
	}

	if pd.IsList() {
		val, err := fieldTypeFromKind(m, pd)
		return "(List " + val + ")", err
	}

	return fieldTypeFromKind(m, pd)
}

func fieldTypeFromKind(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return "Bool", nil

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind:
		return "Int", nil

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "Float", nil

	case protoreflect.StringKind:
		return "String", nil

	case protoreflect.BytesKind:
		m.Imports.Bytes = true
		return "Bytes", nil

	case protoreflect.EnumKind:
		id, err := m.getElmType(pd.Enum().FullName())
		return string(id), err

	case protoreflect.MessageKind, protoreflect.GroupKind:
		id, err := m.getElmType(pd.Message().FullName())
		return string(id), err
	}

	return "", fmt.Errorf("fieldType: unknown protoreflect.Kind: %s", pd.Kind())
}

func fieldZero(m *Module, pd protoreflect.FieldDescriptor) (interface{}, error) {
	if pd.IsMap() {
		return "Dict.empty", nil
	}
	if pd.IsList() {
		return "[]", nil
	}

	switch pd.Kind() {
	// No formatting difference for these fields
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return pd.Default(), nil

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.BoolKind:
		/* If we were to support proto2:
		if proto.Desc.Default().Bool() {
			return "True", nil
		}
		*/
		return "False", nil

	case protoreflect.StringKind:
		return `""`, nil

	case protoreflect.BytesKind:
		return "(BE.encode (BE.sequence []))", nil

	case protoreflect.EnumKind:
		id, err := m.getElmType(pd.Enum().FullName())
		return "empty" + id, err

	case protoreflect.MessageKind, protoreflect.GroupKind:
		id, err := m.getElmType(pd.Message().FullName())
		return "empty" + id, err
	}

	return "", fmt.Errorf("fieldZero: unknown protoreflect.Kind: %s", pd.Kind())
}

func fieldDecoder(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	return fieldCodec(m, "PD.", "Decoder", pd)
}

func fieldEncoder(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	return fieldCodec(m, "PE.", "Encoder", pd)
}

func fieldCodec(m *Module, lib, dir string, pd protoreflect.FieldDescriptor) (string, error) {
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return lib + "bool", nil

	case protoreflect.Int32Kind:
		return lib + "int32", nil
	case protoreflect.Sint32Kind:
		return lib + "sint32", nil
	case protoreflect.Uint32Kind:
		return lib + "uint32", nil
	case protoreflect.Sfixed32Kind:
		return lib + "sfixed32", nil
	case protoreflect.Fixed32Kind:
		return lib + "fixed32", nil

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind:
		return lib + "float", nil
	case protoreflect.DoubleKind:
		return lib + "double", nil

	case protoreflect.StringKind:
		return lib + "string", nil

	case protoreflect.BytesKind:
		return lib + "bytes", nil

	case protoreflect.EnumKind:
		return m.getElmValue(pd.Enum().FullName()) + dir, nil

	case protoreflect.MessageKind, protoreflect.GroupKind:
		return m.getElmValue(pd.Message().FullName()) + dir, nil
	}
	return "", fmt.Errorf("fieldCodec: unknown protoreflect.Kind: %s", pd.Kind())
}

func fieldFuzzer(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	switch pd.Kind() {
	case protoreflect.BoolKind:
		return "Fuzz.bool", nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		m.Fuzzers.Int32 = true
		return "fuzzInt32", nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		m.Fuzzers.Uint32 = true
		return "fuzzUint32", nil

	/* Unsupported by Elm / JS
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64", nil
	*/

	case protoreflect.FloatKind:
		m.Fuzzers.Float32 = true
		return "fuzzFloat32", nil
	case protoreflect.DoubleKind:
		return "Fuzz.float", nil

	case protoreflect.StringKind:
		return "Fuzz.string", nil
	case protoreflect.BytesKind:
		m.Imports.Bytes = true
		return "fuzzBytes", nil

	case protoreflect.EnumKind:
		return m.getElmValue(pd.Enum().FullName()) + "Fuzzer", nil

	case protoreflect.MessageKind, protoreflect.GroupKind:
		return m.getElmValue(pd.Message().FullName()) + "Fuzzer", nil
	}

	return "", fmt.Errorf("kindFuzzer: unknown protoreflect.Kind: %s", pd.Kind())
}
