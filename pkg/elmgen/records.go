package elmgen

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) regMessages(protoMsgs []*protogen.Message) {
	for _, proto := range protoMsgs {
		// Recursive
		m.regEnums(proto.Enums)
		m.regMessages(proto.Messages)
		// Add this msg
		m.registerProtoName(proto.Desc.FullName(), "")
		m.protoMessages = append(m.protoMessages, proto)

		var prefix protoreflect.Name
		if m.config.QualifyNested {
			prefix = protoreflect.Name(proto.Desc.FullName())
		} else {
			prefix = proto.Desc.Name()
		}

		for _, oneof := range proto.Oneofs {
			od := oneof.Desc
			m.registerProtoName(od.FullName(), string(prefix+"."+od.Name()))
			// Build a qualified alias?
			var suffix protoreflect.Name
			if m.config.VariantSuffixes {
				suffix = "."
				if m.config.QualifyNested {
					suffix += protoreflect.Name(od.FullName())
				} else {
					suffix += od.Name()
				}
			}
			// Add fields
			for _, field := range oneof.Fields {
				fd := field.Desc
				m.registerProtoName(fd.FullName(), string(fd.Name()+suffix))
			}
		}
	}
}

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
	return nil
}

func (m *Module) newRecord(proto *protogen.Message) (*Record, error) {
	var record Record
	record.CodecIDs.register(m, proto.Desc.FullName())
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
	field := new(Field)
	field.Label = m.protoFullIdentToElmCasing(string(pd.Name()), false)
	var err error
	if field.Type, err = fieldType(m, pd); err != nil {
		return nil, err
	}
	if field.Zero, err = fieldZero(m, pd); err != nil {
		return nil, err
	}
	if field.Decoder, err = fieldDecoder(m, pd); err != nil {
		return nil, err
	}
	if field.Encoder, err = fieldEncoder(m, pd); err != nil {
		return nil, err
	}
	field.WireNumber = pd.Number()
	field.Cardinality = pd.Cardinality()
	return field, err
}

func (m *Module) newOneofField(po protoreflect.OneofDescriptor) (*Oneof, *Field, error) {
	oneof, err := m.newOneof(po)
	if err != nil {
		return nil, nil, err
	}
	field := &Field{
		m.protoFullIdentToElmCasing(string(po.Name()), false),
		true, 0, 0,
		"Maybe " + string(oneof.ID),
		"Nothing",
		oneof.DecodeID,
		oneof.EncodeID,
	}
	// Optional field?
	if oneof.IsSynthetic {
		// Unwrap type
		field.Type = "Maybe " + oneof.Variants[0].Field.Type
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
		// TODO: check if key is scalar?
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
		id := m.getElmType(pd.Enum().FullName())
		return string(id), nil

	case protoreflect.MessageKind, protoreflect.GroupKind:
		id := m.getElmType(pd.Message().FullName())
		return string(id), nil
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
		id := m.getElmType(pd.Enum().FullName())
		return "empty" + id, nil

	case protoreflect.MessageKind, protoreflect.GroupKind:
		id := m.getElmType(pd.Message().FullName())
		return "empty" + id, nil
	}

	return "", fmt.Errorf("fieldZero: unknown protoreflect.Kind: %s", pd.Kind())
}

func fieldDecoder(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	return fieldCodec(m, "PD.", "decode", pd)
}

func fieldEncoder(m *Module, pd protoreflect.FieldDescriptor) (string, error) {
	return fieldCodec(m, "PE.", "encode", pd)
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
		id := m.getElmType(pd.Enum().FullName())
		return dir + string(id), nil

	case protoreflect.MessageKind, protoreflect.GroupKind:
		id := m.getElmType(pd.Message().FullName())
		return dir + string(id), nil
	}
	return "", fmt.Errorf("fieldCodec: unknown protoreflect.Kind: %s",
		pd.Kind())
}
