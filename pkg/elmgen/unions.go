package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addEnums() {
	for _, proto := range m.protoEnums {
		m.Unions = append(m.Unions,
			m.newUnion(proto))
	}
	sort.Sort(m.Unions)
}

func (m *Module) newUnion(proto *protogen.Enum) *Union {
	union := new(Union)
	union.CodecIDs.register(m, proto.Desc.FullName())
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
			id := m.getElmType(protoVal.Desc.FullName())
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
	return union
}

func (m *Module) newOneof(proto protoreflect.OneofDescriptor) (*Oneof, error) {
	oneof := new(Oneof)
	oneof.IsSynthetic = proto.IsSynthetic()
	// Register codec IDs
	oneof.CodecIDs.register(m, proto.FullName())
	// Add field types
	fields := proto.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		v := new(OneofVariant)
		v.ID = m.getElmType(field.FullName())
		var err error
		if v.Field, err = m.newField(field); err != nil {
			return nil, err
		}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof, nil
}
