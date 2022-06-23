package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addEnums() {
	for _, proto := range m.protoEnums {
		m.Unions = append(m.Unions, m.newUnion(proto))
	}
	sort.Sort(m.Unions)
}

func (m *Module) newUnion(proto *protogen.Enum) *Union {
	union := new(Union)
	union.CodecIDs.register(m, proto.Desc.FullName())
	// Add variants
	aliases := make(map[protoreflect.EnumNumber]*Variant)
	for i, protoVal := range proto.Values {
		v := new(Variant)
		v.Number = protoVal.Desc.Number()
		// Variant (type) or alias (value)?
		if original := aliases[v.Number]; original != nil {
			union.Aliases = append(union.Aliases,
				&VariantAlias{
					original,
					m.getElmValue(protoVal.Desc.FullName())})
		} else {
			v.ID = m.getElmType(protoVal.Desc.FullName())
			aliases[v.Number] = v
			// First is the default
			if i == 0 {
				union.DefaultVariant = v
			} else {
				union.Variants = append(union.Variants, v)
			}
		}
	}
	return union
}

func (m *Module) newOneof(proto protoreflect.OneofDescriptor) (*Oneof, error) {
	var err error
	oneof := new(Oneof)
	oneof.CodecIDs.register(m, proto.FullName())
	oneof.IsSynthetic = proto.IsSynthetic()
	fields := proto.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		v := new(OneofVariant)
		v.ID = m.getElmType(field.FullName())
		if v.Field, err = m.newField(field); err != nil {
			return nil, err
		}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof, nil
}
