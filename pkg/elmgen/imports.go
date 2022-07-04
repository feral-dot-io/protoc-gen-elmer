package elmgen

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	importBytes = "Bytes"
	importDict  = "Dict"
)

func (m *Module) addImport(mod string) {
	if mod != "" {
		m.importsSeen[mod] = true
	}
}

func (m *Module) findImports() {
	// Iterate over fields since they hold non-ref values which can trigger imports
	for _, r := range m.Records {
		for _, f := range r.Fields {
			if f.Oneof != nil {
				fields := f.Oneof.Fields()
				for i := 0; i < fields.Len(); i++ {
					m.fieldImports(fields.Get(i))
				}
			} else {
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
		}
	}
}
