package elmgen

import (
	"sort"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (m *Module) addUnions(enums protoreflect.EnumDescriptors) {
	for i := 0; i < enums.Len(); i++ {
		ed := enums.Get(i)
		union := m.newUnion(ed)
		m.Unions = append(m.Unions, union)
	}
	sort.Sort(m.Unions)
}

func (m *Module) newUnion(ed protoreflect.EnumDescriptor) *Union {
	union := new(Union)
	union.Type = NewElmType(ed.ParentFile(), ed)
	// Add variants
	aliases := make(map[protoreflect.EnumNumber]*Variant)
	values := ed.Values()
	for i := 0; i < values.Len(); i++ {
		vd := values.Get(i)
		num := vd.Number()
		// Variant (type) or alias (value)?
		if original := aliases[num]; original != nil {
			alias := NewElmValue(vd.ParentFile(), vd)
			union.Aliases = append(union.Aliases,
				&VariantAlias{original, alias})
		} else {
			// Create
			id := NewElmType(vd.ParentFile(), vd).ElmRef
			v := &Variant{&id, num}
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

func (m *Module) newOneof(od protoreflect.OneofDescriptor) *Oneof {
	oneof := new(Oneof)
	oneof.IsSynthetic = od.IsSynthetic()
	if oneof.IsSynthetic {
		firstField := od.Fields().Get(0)
		oneof.Type = NewElmType(firstField.ParentFile(), firstField)
	} else {
		oneof.Type = NewElmType(od.ParentFile(), od)
	}
	// Add field types
	fields := od.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		v := &OneofVariant{
			&NewElmType(fd.ParentFile(), fd).ElmRef,
			m.newField(fd)}
		oneof.Variants = append(oneof.Variants, v)
	}
	return oneof
}
