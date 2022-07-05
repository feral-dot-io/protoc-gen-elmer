package elmgen

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	importBytes     = "Bytes"
	importDict      = "Dict"
	importGooglePB  = "Google.Protobuf"
	importElmer     = "Protobuf.Elmer"
	importElmerTest = "Protobuf.ElmerTest"
)

func (m *Module) addImport(mod string) {
	if mod != "" {
		m.importsSeen[mod] = true
	}
}

func (m *Module) findImports() {
	if len(m.Unions) > 0 {
		m.addImport(importElmerTest)
	}
	// Iterate over fields since they hold non-ref values which can trigger imports
	for _, r := range m.Records {
		for _, f := range r.Fields {
			if f.Oneof != nil {
				for _, v := range f.Oneof.Variants {
					m.fieldImports(v.Field.Desc)
				}
			}
			if f.Desc != nil { // Including oneof optional
				m.fieldImports(f.Desc)
			}
		}
	}
}

func (m *Module) fieldImports(fd protoreflect.FieldDescriptor) {
	if fd.IsMap() { // Dict
		m.addImport(importDict)
		m.fieldImports(fd.MapKey())
		m.fieldImports(fd.MapValue())
	} else {
		switch fd.Kind() {
		case protoreflect.BytesKind: // Bytes
			m.addImport(importBytes)
			m.addImport(importElmer)
			m.addImport(importElmerTest)

		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind:
			m.addImport(importElmerTest)

		case protoreflect.EnumKind:
			ed := fd.Enum()
			m.NewElmType(ed.ParentFile(), ed)

		case protoreflect.MessageKind, protoreflect.GroupKind:
			md := fd.Message()
			m.NewElmType(md.ParentFile(), md)
		}
	}
}
