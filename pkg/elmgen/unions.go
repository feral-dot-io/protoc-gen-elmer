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
	union.Type = m.NewElmType(ed.ParentFile(), ed)
	union.Comments = NewCommentSet(enum.Comments)
	// Add variants
	aliases := make(map[protoreflect.EnumNumber]*Variant)
	for _, value := range enum.Values {
		vd := value.Desc
		num := vd.Number()
		// Variant (type) or alias (value)?
		if original := aliases[num]; original != nil {
			alias := &VariantAlias{
				m.NewElmValue(vd.ParentFile(), vd),
				original,
				NewCommentSet(value.Comments)}
			union.Aliases = append(union.Aliases, alias)
		} else {
			// Create
			id := m.NewElmType(vd.ParentFile(), vd).ElmRef
			v := &Variant{id, num, NewCommentSet(value.Comments)}
			// Add
			union.Variants = append(union.Variants, v)
			aliases[v.Number] = v
		}
	}
	return union
}

// Returns a union's default variant. Never returns nil.
// All unions have at least one variant. The first is our default (zero)
func (u *Union) Default() *Variant {
	return u.Variants[0]
}

func (m *Module) newOneof(protoOneof *protogen.Oneof) *Oneof {
	od := protoOneof.Desc
	oneof := new(Oneof)
	oneof.Comments = NewCommentSet(protoOneof.Comments)
	oneof.IsSynthetic = od.IsSynthetic()
	if oneof.IsSynthetic {
		firstField := od.Fields().Get(0)
		oneof.Type = m.NewElmType(firstField.ParentFile(), firstField)
	} else {
		oneof.Type = m.NewElmType(od.ParentFile(), od)
	}
	// Add field types
	for _, field := range protoOneof.Fields {
		fd := field.Desc
		v := &OneofVariant{
			m.NewElmType(fd.ParentFile(), fd).ElmRef,
			m.newField(field)}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof
}
