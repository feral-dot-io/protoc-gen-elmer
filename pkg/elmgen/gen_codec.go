// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Formats our Elm code and returns the replacement File. Errors are passed to plugin.Error. Assumes `elm-format` is in `$PATH`
func FormatFile(plugin *protogen.Plugin, path string, file *protogen.GeneratedFile) *protogen.GeneratedFile {
	// Fetch unformatted content
	content, err := file.Content()
	plugin.Error(err)
	unformatted := bytes.NewBuffer(content)
	file.Skip()
	// Run through elm-format
	formatted := plugin.NewGeneratedFile(path, "")
	err = runElmFormat(unformatted, formatted)
	plugin.Error(err)
	return formatted
}

func runElmFormat(in io.Reader, out io.Writer) error {
	cmd := exec.Command("elm-format", "--yes", "--stdin")
	cmd.Stdin = in
	cmd.Stderr = os.Stderr
	cmd.Stdout = out
	return cmd.Run()
}

// Prints a do not edit string that editors may recognise as likely from codegen. It's Go specific and doesn't quite match the expected regex so we do our best to make it look similar
func printDoNotEdit(g *protogen.GeneratedFile) {
	// Ideally this would be left-aligned but elm-format will indent it
	g.P("-- // Code generated protoc-gen-elmer DO NOT EDIT \\\\")
}

// Formats an Elm comment by trimming white space and prefixing it with double dashes (--). String may contain new lines.
func (c Comments) String() (block string) {
	if c == "" {
		return ""
	}
	dropTrailingNl := strings.TrimSuffix(string(c), "\n")
	for _, line := range strings.Split(dropTrailingNl, "\n") {
		block += "-- " + strings.TrimSpace(line) + "\n"
	}
	return
}

// Prints the leading comments from a set. The main comments are printed as a document block comment ({-|})
func (set *CommentSet) printBlock(g *protogen.GeneratedFile) {
	for _, c := range set.LeadingDetached {
		g.P(c)
	}
	if set.Leading != "" {
		g.P("{-| ", string(set.Leading), " -}")
	}
}

// Prints the trailing comments of a block.
func (set *CommentSet) printBlockTrailing(g *protogen.GeneratedFile) {
	g.P(set.Trailing)
}

// Prints all comments with a `-- ` prefix. Includes trailing as the formatter will get rid of excess new lines.
func (set *CommentSet) printDashDash(g *protogen.GeneratedFile) {
	if len(set.LeadingDetached) > 0 {
		for _, c := range set.LeadingDetached {
			g.P(c)
		}
		if set.Leading != "" {
			g.P("-- ")
		}
	}
	g.P(set.Leading)
}

// Prints the imports of a module. Since a `Module` holds references to tests via types, these should be skipped for non-test generators.
func printImports(g *protogen.GeneratedFile, m *Module, skipTests bool) {
	g.P("import Protobuf.Decode as PD")
	g.P("import Protobuf.Encode as PE")
	for _, i := range m.Imports {
		// Skip tests? Since our Elm types always generate a reference to Tests, we need to be able to skip them
		if skipTests && strings.HasSuffix(i, "Tests") {
			continue
		}
		switch i {
		case "Bytes":
			g.P("import Bytes exposing (Bytes)")
			g.P("import Bytes.Encode as BE")
		case "Dict":
			g.P("import Dict exposing (Dict)")
		default:
			g.P("import ", i)
		}
	}
}

// Generates Elm decoders and encoders (making a codec) to a file
func GenerateCodec(m *Module, g *protogen.GeneratedFile) bool {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}
	g.P("module ", m.Name, " exposing (..)")

	g.P("{-| Protobuf library for decoding and encoding structures found in package `" + m.ProtoPackage + "` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.")
	g.P("")

	// List relevant top-level types
	var allTypes []*ElmType
	if len(m.Records) > 0 {
		g.P("Records:")
		for _, r := range m.Records {
			g.P("- ", r.Type.ID)
			allTypes = append(allTypes, r.Type)
		}
	} else {
		g.P("Records: (none)")
	}
	g.P("")
	var docsValues, docsStr []string
	if len(m.Unions) > 0 {
		g.P("Unions:")
		for _, u := range m.Unions {
			t := u.Type
			g.P("- ", t.ID)
			allTypes = append(allTypes, t)
			docsValues = append(docsValues, "valuesOf"+t.ID)
			docsStr = append(docsStr, "from"+t.ID, "to"+t.ID)
		}
	} else {
		g.P("Unions: (none)")
	}
	g.P("")

	g.P("Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.")
	g.P("")
	g.P("Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.")
	g.P("")

	// Build lists of IDs for @docs
	var docsTypes, docsEmpty, docsDecs, docsEncs []string
	for _, t := range allTypes {
		docsTypes = append(docsTypes, t.ID)
		docsEmpty = append(docsEmpty, "empty"+t.ID)
		docsDecs = append(docsDecs, "decode"+t.ID)
		docsEncs = append(docsEncs, "encode"+t.ID)
	}
	g.P("# Types")
	g.P("@docs ", strings.Join(docsTypes, ", "))
	g.P("# Empty (zero values)")
	g.P("@docs ", strings.Join(docsEmpty, ", "))
	if len(docsStr) > 0 {
		g.P("# Enum valuesOf")
		g.P("@docs ", strings.Join(docsValues, ", "))
		g.P("# Enum and String converters")
		g.P("@docs ", strings.Join(docsStr, ", "))
	}
	g.P("# Decoders")
	g.P("@docs ", strings.Join(docsDecs, ", "))
	g.P("# Encoders")
	g.P("@docs ", strings.Join(docsEncs, ", "))
	g.P("-}")
	printDoNotEdit(g)
	printImports(g, m, true)

	// Unions
	for _, u := range m.Unions {
		u.Comments.printBlock(g)
		g.P("type ", u.Type)
		for i, v := range u.Variants {
			prefix := "|"
			if i == 0 {
				prefix = "="
			}
			v.Comments.printDashDash(g)
			gFP("    %s %s %s", prefix, v.ID, v.Comments.Trailing)
		}
		for _, a := range u.Aliases {
			a.Comments.printDashDash(g)
			gFP("%s : %s", a.Alias, u.Type)
			gFP("%s =", a.Alias)
			gFP("    %s", a.Comments.Trailing)
			gFP("    %s", a.Variant.ID)
		}
		u.Comments.printBlockTrailing(g)
	}

	// Records
	for _, r := range m.Records {
		r.Comments.printBlock(g)
		gFP("type alias %s =", r.Type)
		for i, f := range r.Fields {
			prefix := ","
			if i == 0 {
				prefix = "{"
			}
			f.Comments.printDashDash(g)
			gFP("    %s %s : %s %s",
				prefix, f.Label, fieldType(m, f), f.Comments.Trailing)
		}
		if len(r.Fields) == 0 {
			gFP("    {")
		}
		gFP("    }")
		r.Comments.printBlockTrailing(g)
	}

	// Oneofs (nested unions)
	for _, o := range m.Oneofs {
		// Skip optional
		if o.IsSynthetic {
			continue
		}
		gFP("type %s", o.Type)
		for i, v := range o.Variants {
			prefix := "|"
			if i == 0 {
				prefix = "="
			}
			v.Field.Comments.printDashDash(g)
			gFP("    %s %s %s %s",
				prefix, v.ID, fieldType(m, v.Field), v.Field.Comments.Trailing)
		}
	}

	// Zero records
	for _, r := range m.Records {
		gFP("%s : %s", r.Type.Zero, r.Type)
		gFP("%s =", r.Type.Zero)
		zeros := []interface{}{"    ", r.Type}
		for _, f := range r.Fields {
			var zero interface{}
			if f.Oneof != nil {
				zero = "Nothing"
			} else {
				zero = fieldZero(m, f.Desc)
			}
			zeros = append(zeros, " ", zero)
		}
		g.P(zeros...)
	}

	// Zero unions
	for _, u := range m.Unions {
		t := u.Type
		gFP("%s : %s", t.Zero, t)
		gFP("%s =", t.Zero)
		gFP("    %s", u.Default().ID)
	}

	// Extra union helpers: valuesOf, to, and from
	for _, u := range m.Unions {
		var values []string
		for _, v := range u.Variants {
			values = append(values, v.ID.ID)
		}
		gFP("valuesOf%s : List %s", u.Type.ID, u.Type.ID)
		gFP("valuesOf%s = [ %s ]", u.Type.ID, strings.Join(values, ", "))

		gFP("from%s : %s -> String", u.Type.ID, u.Type.ID)
		gFP("from%s u =", u.Type.ID)
		gFP("  case u of")
		for _, v := range u.Variants {
			gFP(`    %s ->`, v.ID)
			gFP(`      "%s"`, v.Label)
		}

		gFP("to%s : String -> %s", u.Type.ID, u.Type.ID)
		gFP("to%s str =", u.Type.ID)
		gFP("  case str of")
		for _, v := range u.Variants {
			gFP(`    "%s" ->`, v.Label)
			gFP(`      %s`, v.ID)
		}
		// No match? Use default
		gFP("    _ ->")
		gFP("      %s", u.Default().ID)
	}

	// Record decoders
	for _, r := range m.Records {
		gFP("%s : PD.Decoder %s", r.Type.Decoder, r.Type)
		gFP("%s =", r.Type.Decoder)
		// Build oneof decoders inline since they're unique to the message
		// Ideally they'd be inline here (and in the decoder + fuzzer)
		if oneofs := r.Oneofs(); len(oneofs) > 0 {
			g.P("    let")
			for _, f := range oneofs {
				o := f.Oneof
				gFP("        %s =", o.Type.Decoder)
				g.P("            [")
				for j, v := range o.Variants {
					prefix := "            "
					if j != 0 {
						prefix += ","
					}
					wire := v.Field.Desc.Number()
					decoder := fieldDecoder(m, v.Field.Desc)
					if o.IsSynthetic { // Only field, skip map
						gFP("%s( %d, %s )",
							prefix, wire, decoder)
					} else {
						gFP("%s( %d, PD.map %s %s )",
							prefix, wire, v.ID, decoder)
					}
				}
				g.P("                ]")
			}
			g.P("    in")
		}
		g.P("    PD.message ", r.Type.Zero)
		g.P("        [")
		for i, f := range r.Fields {
			prefix := "            "
			if i != 0 {
				prefix += ","
			}
			getter := "(\\v m -> { m | %s = v })"
			// Pick a FieldDecoder
			if f.Oneof != nil {
				gFP("%s PD.oneOf %s "+getter,
					prefix, f.Oneof.Type.Decoder, f.Label)
			} else {
				wire := f.Desc.Number()
				decoder := fieldDecoder(m, f.Desc)
				if f.Desc.IsMap() {
					key := f.Desc.MapKey()
					val := f.Desc.MapValue()
					gFP("%s PD.mapped %d ( %s , %s ) %s %s .%s "+getter,
						prefix, wire,
						fieldZero(m, key), fieldZero(m, val),
						fieldDecoder(m, key), fieldDecoder(m, val),
						f.Label, f.Label)
				} else {
					switch f.Desc.Cardinality() {
					case protoreflect.Optional:
						gFP("%s PD.optional %d %s "+getter,
							prefix, wire, decoder, f.Label)

					case protoreflect.Required:
						gFP("%s PD.required %d %s "+getter,
							prefix, wire, decoder, f.Label)

					case protoreflect.Repeated:
						gFP("%s PD.repeated %d %s .%s "+getter,
							prefix, wire, decoder, f.Label, f.Label)
					}
				}
			}
		}
		g.P("        ]")
	}

	// Union decoders
	for _, u := range m.Unions {
		t := u.Type
		gFP("%s : PD.Decoder %s", t.Decoder, t)
		gFP("%s =", t.Decoder)
		g.P("    let")
		g.P("        conv v =")
		g.P("            case v of")
		for _, v := range u.Variants {
			g.P("                ", v.Number, " ->")
			g.P("                    ", v.ID)
		}
		g.P("                _ ->")
		g.P("                    ", u.Default().ID)
		g.P("    in")
		g.P("    PD.map conv PD.int32")
	}

	// Record encoders
	for _, r := range m.Records {
		param := "v"
		if len(r.Fields) == 0 {
			param = "_"
		}
		gFP("%s : %s -> PE.Encoder", r.Type.Encoder, r.Type)
		gFP("%s %s =", r.Type.Encoder, param)
		oneofs := r.Oneofs()
		if len(oneofs) > 0 {
			g.P("    let")
			for _, f := range oneofs {
				o := f.Oneof
				ws := "        "
				gFP("%s%s o =", ws, o.Type.Encoder)
				gFP("%s    case o of", ws)
				ws += "        "
				for _, v := range o.Variants {
					f := v.Field
					id := v.ID.String()
					if o.IsSynthetic { // No sub-enum
						id = ""
					}
					gFP("%sJust (%s data) ->", ws, id)
					gFP("%s    [ ( %d, %s data ) ]",
						ws, f.Desc.Number(), fieldEncoder(m, f.Desc))
				}
				// Nil isn't encoded on the wire
				gFP("%sNothing ->", ws)
				gFP("%s    []", ws)
			}
			g.P("    in")
		}
		g.P("    PE.message <|")
		g.P("        [")
		// Regular (non-oneof) fields
		var written bool
		for _, f := range r.Fields {
			if f.Oneof != nil { // Skip
				continue
			}
			prefix := "            "
			if written { // Can't do i != 0 because of "continue"
				prefix += ","
			}
			encoder := fieldEncoder(m, f.Desc)
			// Special fields?
			if f.Desc.IsMap() {
				keyEnc := fieldEncoder(m, f.Desc.MapKey())
				valEnc := fieldEncoder(m, f.Desc.MapValue())
				gFP("%s ( %d, PE.dict %s %s v.%s )",
					prefix, f.Desc.Number(), keyEnc, valEnc, f.Label)
			} else if f.Desc.Cardinality() == protoreflect.Repeated {
				gFP("%s ( %d, PE.list %s v.%s )",
					prefix, f.Desc.Number(), encoder, f.Label)
			} else {
				gFP("%s ( %d, %s v.%s )",
					prefix, f.Desc.Number(), encoder, f.Label)
			}
			written = true
		}
		g.P("        ]")
		if len(oneofs) > 0 {
			// Oneof field handling
			for _, f := range oneofs {
				gFP("        ++ %s v.%s", f.Oneof.Type.Encoder, f.Label)
			}
		}
	}

	// Union encoders
	for _, u := range m.Unions {
		t := u.Type
		gFP("%s : %s -> PE.Encoder", t.Encoder, t)
		gFP("%s v =", t.Encoder)
		g.P("    let")
		g.P("        conv =")
		g.P("            case v of")
		for _, v := range u.Variants {
			g.P("                ", v.ID, " ->")
			g.P("                    ", v.Number)
		}
		g.P("    in")
		g.P("    PE.int32 conv")
	}

	return true
}

func fieldType(m *Module, f *Field) string {
	if f.Oneof != nil {
		var inner string
		if f.Desc != nil { // Optional
			inner = fieldTypeDesc(m, f.Desc)
		} else {
			inner = f.Oneof.Type.String()
		}
		return "(Maybe " + inner + ")"
	}
	return fieldTypeDesc(m, f.Desc)
}

func fieldTypeDesc(m *Module, fd protoreflect.FieldDescriptor) string {
	if fd.IsMap() {
		key := fieldTypeKind(m, fd.MapKey())
		val := fieldTypeDesc(m, fd.MapValue())
		return "(Dict " + key + " " + val + ")"
	} else if fd.IsList() {
		val := fieldTypeKind(m, fd)
		return "(List " + val + ")"
	}
	return fieldTypeKind(m, fd)
}

func fieldTypeKind(m *Module, fd protoreflect.FieldDescriptor) string {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return "Bool"

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind:
		return "Int"

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "Float"

	case protoreflect.StringKind:
		return "String"

	case protoreflect.BytesKind:
		return "Bytes"

	case protoreflect.EnumKind:
		ed := fd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := fd.Message()
		return m.NewElmType(md.ParentFile(), md).String()
	}

	log.Panicf("fieldType: unknown protoreflect.Kind: %s", fd.Kind())
	return ""
}

func fieldZero(m *Module, fd protoreflect.FieldDescriptor) string {
	if fd.IsMap() { // Dict
		return "Dict.empty"
	} else if fd.IsList() { // List
		return "[]"
	}

	switch fd.Kind() {
	// No formatting difference for these fields
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return fd.Default().String()

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
	//	protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.BoolKind:
		if fd.Default().Bool() {
			return "True"
		}
		return "False"

	case protoreflect.StringKind:
		return `""`

	case protoreflect.BytesKind:
		return "Protobuf.Elmer.emptyBytes"

	case protoreflect.EnumKind:
		ed := fd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).Zero.String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := fd.Message()
		return m.NewElmType(md.ParentFile(), md).Zero.String()
	}

	log.Panicf("fieldZero: unknown protoreflect.Kind: %s", fd.Kind())
	return ""
}

func fieldDecoder(m *Module, fd protoreflect.FieldDescriptor) string {
	return fieldCodecKind(m, "PD.", fd)
}

func fieldEncoder(m *Module, fd protoreflect.FieldDescriptor) string {
	return fieldCodecKind(m, "PE.", fd)
}

// Just the Kind. Does not take into account special features like lists.
func fieldCodecKind(m *Module, lib string, fd protoreflect.FieldDescriptor) string {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return lib + "bool"

	case protoreflect.Int32Kind:
		return lib + "int32"
	case protoreflect.Sint32Kind:
		return lib + "sint32"
	case protoreflect.Uint32Kind:
		return lib + "uint32"
	case protoreflect.Sfixed32Kind:
		return lib + "sfixed32"
	case protoreflect.Fixed32Kind:
		return lib + "fixed32"

	// Unsupported by Elm / JS
	//case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:

	case protoreflect.FloatKind:
		return lib + "float"
	case protoreflect.DoubleKind:
		return lib + "double"

	case protoreflect.StringKind:
		return lib + "string"

	case protoreflect.BytesKind:
		return lib + "bytes"

	case protoreflect.EnumKind:
		ed := fd.Enum()
		return m.fieldCodecElmType(lib, ed.ParentFile(), ed)

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := fd.Message()
		return m.fieldCodecElmType(lib, md.ParentFile(), md)
	}

	log.Panicf("fieldCodec: unknown protoreflect.Kind: %s", fd.Kind())
	return ""
}

func (m *Module) fieldCodecElmType(lib string, p packager, d fullNamer) string {
	t := m.NewElmType(p, d)
	if lib == "PD." {
		return t.Decoder.String()
	} else {
		return t.Encoder.String()
	}
}
