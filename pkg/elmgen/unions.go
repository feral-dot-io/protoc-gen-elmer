package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addEnums() error {
	for _, proto := range m.protoEnums {
		union, err := m.newUnion(proto)
		if err != nil {
			return err
		}
		m.Unions = append(m.Unions, union)
	}
	sort.Sort(m.Unions)
	return nil
}

func (m *Module) newUnion(proto *protogen.Enum) (*Union, error) {
	union := new(Union)
	err := union.CodecIDs.register(m, proto.Desc.FullName())
	if err != nil {
		return nil, err
	}
	// Add variants
	aliases := make(map[protoreflect.EnumNumber]*Variant)
	for i, protoVal := range proto.Values {
		num := protoVal.Desc.Number()
		// Variant (type) or alias (value)?
		if original := aliases[num]; original != nil {
			elmID := m.getElmValue(protoVal.Desc.FullName())
			union.Aliases = append(union.Aliases,
				&VariantAlias{original, elmID})
		} else {
			// Create
			id, err := m.getElmType(protoVal.Desc.FullName())
			if err != nil {
				return nil, err
			}
			v := &Variant{id, num}
			// Add
			if i == 0 { // First is the default
				union.DefaultVariant = v
			} else {
				union.Variants = append(union.Variants, v)
			}
			aliases[v.Number] = v
		}
	}
	return union, nil
}

func (m *Module) newOneof(proto protoreflect.OneofDescriptor) (*Oneof, error) {
	oneof := new(Oneof)
	oneof.IsSynthetic = proto.IsSynthetic()
	// Register codec IDs
	err := oneof.CodecIDs.register(m, proto.FullName())
	if err != nil {
		return nil, err
	}
	// Add field types
	fields := proto.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		v := new(OneofVariant)
		if v.ID, err = m.getElmType(field.FullName()); err != nil {
			return nil, err
		}
		if v.Field, err = m.newField(field); err != nil {
			return nil, err
		}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof, nil
}
