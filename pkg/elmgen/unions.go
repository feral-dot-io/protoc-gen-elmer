package elmgen

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addUnions(enums []*protogen.Enum) {
	for _, enum := range enums {
		union := m.newUnion(enum)
		m.Unions = append(m.Unions, union)
	}
	sort.Sort(m.Unions)
}

func (m *Module) newUnion(enum *protogen.Enum) *Union {
	ed := enum.Desc
	union := new(Union)
	union.Type = NewElmType(ed.ParentFile(), ed)
	union.Comments = NewCommentSet(enum.Comments)
	// Add variants
	aliases := make(map[protoreflect.EnumNumber]*Variant)
	for i, value := range enum.Values {
		vd := value.Desc
		num := vd.Number()
		// Variant (type) or alias (value)?
		if original := aliases[num]; original != nil {
			alias := NewElmValue(vd.ParentFile(), vd)
			union.Aliases = append(union.Aliases,
				&VariantAlias{original, alias})
		} else {
			// Create
			id := NewElmType(vd.ParentFile(), vd).ElmRef
			v := &Variant{&id, num, NewCommentSet(value.Comments)}
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

func (m *Module) newOneof(protoOneof *protogen.Oneof) *Oneof {
	od := protoOneof.Desc
	oneof := new(Oneof)
	oneof.Comments = NewCommentSet(protoOneof.Comments)
	oneof.IsSynthetic = od.IsSynthetic()
	if oneof.IsSynthetic {
		firstField := od.Fields().Get(0)
		oneof.Type = NewElmType(firstField.ParentFile(), firstField)
	} else {
		oneof.Type = NewElmType(od.ParentFile(), od)
	}
	// Add field types
	for _, field := range protoOneof.Fields {
		fd := field.Desc
		v := &OneofVariant{
			&NewElmType(fd.ParentFile(), fd).ElmRef,
			m.newField(field)}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof
}
