package elmgen

import (
	"strings"

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

/* Elm type constructors */

func (m *Module) newElmRef(mod, id string) *ElmRef {
	ref := &ElmRef{mod, id}
	if mod == m.Name { // Local?
		ref.Module = ""
	}
	m.addImport(ref.Module)
	return ref
}

func (m *Module) NewElmValue(p Packager, d FullNamer) *ElmRef {
	mod, _, asValue := protoReflectToElm(p, d)
	return m.newElmRef(mod, asValue)
}

func (m *Module) NewElmType(p Packager, d FullNamer) *ElmType {
	mod, asType, asValue := protoReflectToElm(p, d)
	// Well-known type handling
	if mod == importGooglePB {
		// Use our own library?
		if strings.HasSuffix(asType, "Value") &&
			asType != "Value" && asType != "EnumValue" &&
			asType != "NullValue" && asType != "ListValue" ||
			asType == "Timestamp" {
			return &ElmType{
				m.newElmRef(importElmer, asType),
				m.newElmRef(importElmer, "empty"+asType),
				m.newElmRef(importElmer, asValue+"Decoder"),
				m.newElmRef(importElmer, asValue+"Encoder"),
				m.newElmRef(importElmerTest, asValue+"Fuzzer")}
		} else {
			// Passthru to Google.Protobuf
			gpType, gpValue := asType, asValue
			switch asType {
			case "Field_Cardinality":
				gpType, gpValue = "Cardinality", "cardinality"
			case "Field_Kind":
				gpType, gpValue = "Kind", "kind"
			case "XType":
				gpType, gpValue = "Type", "type"
			}

			return &ElmType{
				m.newElmRef(importGooglePB, gpType),
				m.newElmRef(importElmer, "empty"+asType),
				m.newElmRef(importGooglePB, gpValue+"Decoder"),
				m.newElmRef(importGooglePB, "to"+gpType+"Encoder"),
				m.newElmRef(importElmerTest, asValue+"Fuzzer")}
		}
	}
	return &ElmType{
		m.newElmRef(mod, asType),
		m.newElmRef(mod, "empty"+asType),
		m.newElmRef(mod, asValue+"Decoder"),
		m.newElmRef(mod, asValue+"Encoder"),
		m.newElmRef(mod+"Tests", asValue+"Fuzzer")}
}

func (r *ElmRef) String() string {
	if r.Module == "" { // Local ref
		return r.ID
	}
	return r.Module + "." + r.ID
}
