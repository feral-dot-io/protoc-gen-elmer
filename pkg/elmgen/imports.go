// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	importBytes      = "Bytes"
	importDict       = "Dict"
	importGooglePB   = "Google.Protobuf"
	importElmer      = "Protobuf.Elmer"
	importElmerTests = "Protobuf.ElmerTests"
)

// Adds a new Module import. Must be an Elm Module reference e.g. "Protobuf.Decode"
func (m *Module) addImport(mod string) {
	if mod != "" {
		m.importsSeen[mod] = true
	}
}

// Finds extra imports once module has been filled with known data strucutres
func (m *Module) findImports() {
	m.addImport(importElmerTests) // Needed by all tests, removed by non-test modules
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

// Finds imports on a record field
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

type (
	packager interface {
		Package() protoreflect.FullName
	}
	fullNamer interface {
		FullName() protoreflect.FullName
	}
)

// Converts a protoreflect package and ident to an Elm module plus type or value
func protoReflectToElm(p packager, d fullNamer) (mod, asType, asValue string) {
	pkg, fullIdent := string(p.Package()), string(d.FullName())
	mod = protoPkgToElmModule(pkg)
	postPkg := strings.TrimPrefix(fullIdent, pkg+".")
	asType, asValue = protoIdentToElmID(postPkg)
	return
}

// Creates a new Elm value (lowercase first char) from a proto ident. Prefix must be non-empty
func (m *Module) NewElmValue(p packager, prefix string, d fullNamer) *ElmRef {
	mod, asType, _ := protoReflectToElm(p, d)
	return m.newElmRef(mod, prefix+asType)
}

// Creates a new Elm type reference (uppercase first char) from a proto ident
func (m *Module) NewElmType(p packager, d fullNamer) *ElmType {
	mod, asType, asValue := protoReflectToElm(p, d)
	// Well-known type handling
	if mod == importGooglePB {
		// Use our own library?
		if strings.HasSuffix(asType, "Value") &&
			asType != "Value" && asType != "EnumValue" &&
			asType != "NullValue" && asType != "ListValue" {
			return &ElmType{
				m.newElmRef(importElmer, asType),
				m.newElmRef(importElmer, "empty"+asType),
				m.newElmRef(importElmer, "decode"+asType),
				m.newElmRef(importElmer, "encode"+asType),
				m.newElmRef(importElmerTests, "fuzz"+asType)}
		} else if asType == "Timestamp" {
			return &ElmType{
				m.newElmRef("Time", "Posix"),
				m.newElmRef(importElmer, "empty"+asType),
				m.newElmRef(importElmer, "decode"+asType),
				m.newElmRef(importElmer, "encode"+asType),
				m.newElmRef(importElmerTests, "fuzz"+asType)}
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
				m.newElmRef(importElmerTests, "fuzz"+asType)}
		}
	}
	return &ElmType{
		m.newElmRef(mod, asType),
		m.newElmRef(mod, "empty"+asType),
		m.newElmRef(mod, "decode"+asType),
		m.newElmRef(mod, "encode"+asType),
		m.newElmRef(mod+"Tests", "fuzz"+asType)}
}

// Converts an Elm reference to Elm code. If local, drops the module.
func (r *ElmRef) String() string {
	if r.Module == "" { // Local ref
		return r.ID
	}
	return r.Module + "." + r.ID
}
