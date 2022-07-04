package elmgen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

func printDoNotEdit(g *protogen.GeneratedFile) {
	// Ideally this would be left-aligned but elm-format will indent it
	g.P("{-")
	g.P("// Code generated protoc-gen-elmer DO NOT EDIT \\\\")
	g.P("-}")
}

func (c Comments) String() (block string) {
	if c == "" {
		return ""
	}
	dropTrailingNl := strings.TrimSuffix(string(c), "\n")
	for _, line := range strings.Split(dropTrailingNl, "\n") {
		block += "-- " + line + "\n"
	}
	return
}

func (set *CommentSet) printBlock(g *protogen.GeneratedFile) {
	for _, c := range set.LeadingDetached {
		g.P(c)
	}
	if set.Leading != "" {
		g.P("{-| ", string(set.Leading), " -}")
	}
}

func (set *CommentSet) printBlockTrailing(g *protogen.GeneratedFile) {
	g.P(set.Trailing)
}

// Prints all comments with a `-- ` prefix. Includes trailing as the formatter will get rid of excess new lines.
func (set *CommentSet) printDashDash(g *protogen.GeneratedFile) {
	if len(set.LeadingDetached) > 0 {
		for _, c := range set.LeadingDetached {
			g.P(c)
		}
		g.P("-- ")
	}
	g.P(set.Leading)
	if set.Trailing != "" {
		g.P("-- ")
		g.P(set.Trailing)
	}
}

func printImports(g *protogen.GeneratedFile, m *Module, skipTests bool) {
	g.P("import Protobuf.Decode as PD")
	g.P("import Protobuf.Encode as PE")
	for _, i := range m.Imports {
		// Skip tests? TODO remove once Fuzzer pushed to gen_*
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

func GenerateCodec(m *Module, g *protogen.GeneratedFile) {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}
	g.P("module ", m.Name, " exposing (..)")
	printDoNotEdit(g)
	printImports(g, m, true)

	// Unions
	for _, u := range m.Unions {
		u.Comments.printBlock(g)
		g.P("type ", u.Type)
		u.DefaultVariant.Comments.printDashDash(g)
		g.P("    = ", u.DefaultVariant.ID, " Int")
		for _, v := range u.Variants {
			v.Comments.printDashDash(g)
			g.P("    | ", v.ID)
		}
		for _, a := range u.Aliases {
			gFP("%s : %s", a.Alias, u.Type)
			gFP("%s = %s", a.Alias, a.ID)
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
			gFP("    %s %s : %s", prefix, f.Label, f.Type)
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
			gFP("    %s %s %s", prefix, v.ID, v.Field.Type)
		}
	}

	// Zero records
	for _, r := range m.Records {
		gFP("%s : %s", r.Type.Zero, r.Type)
		gFP("%s =", r.Type.Zero)
		zeros := []interface{}{"    ", r.Type}
		for _, f := range r.Fields {
			zero := f.Zero
			if f.IsMap {
				zero = "Dict.empty"
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
		gFP("    %s 0", u.DefaultVariant.ID)
	}

	// Record decoders
	for _, r := range m.Records {
		gFP("%s : PD.Decoder %s", r.Type.Decoder, r.Type)
		gFP("%s =", r.Type.Decoder)
		// Build oneof decoders inline since they're unique to the message
		// Ideally they'd be inline here (and in the decoder + fuzzer)
		if len(r.Oneofs) > 0 {
			g.P("    let")
			for _, o := range r.Oneofs {
				gFP("        %s =", o.Type.Decoder)
				g.P("            [")
				for j, v := range o.Variants {
					prefix := "            "
					if j != 0 {
						prefix += ","
					}
					if o.IsSynthetic { // Only field, skip map
						gFP("%s( %d, %s )",
							prefix, v.Field.WireNumber, v.Field.Decoder)
					} else {
						gFP("%s( %d, PD.map %s %s )",
							prefix, v.Field.WireNumber, v.ID, v.Field.Decoder)
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
			if f.IsOneof {
				gFP("%s PD.oneOf %s "+getter, prefix, f.Decoder, f.Label)
			} else if f.IsMap {
				gFP("%s PD.mapped %d ( %s , %s ) %s %s .%s "+getter,
					prefix, f.WireNumber,
					f.Key.Zero, f.Zero, f.Key.Decoder, f.Decoder,
					f.Label, f.Label)
			} else {
				switch f.Cardinality {
				case protoreflect.Optional:
					gFP("%s PD.optional %d %s "+getter,
						prefix, f.WireNumber, f.Decoder, f.Label)

				case protoreflect.Required:
					gFP("%s PD.required %d %s "+getter,
						prefix, f.WireNumber, f.Decoder, f.Label)

				case protoreflect.Repeated:
					gFP("%s PD.repeated %d %s .%s "+getter,
						prefix, f.WireNumber, f.Decoder, f.Label, f.Label)
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
		g.P("                wire ->")
		g.P("                    ", u.DefaultVariant.ID, " wire")
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
		if len(r.Oneofs) > 0 {
			g.P("    let")
			for _, o := range r.Oneofs {
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
					gFP("%s    [ ( %d, %s data ) ]", ws, f.WireNumber, f.Encoder)
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
			if f.IsOneof { // Skip
				continue
			}
			prefix := "            "
			if written { // Can't do i != 0 because of "continue"
				prefix += ","
			}
			// Special fields?
			if f.IsMap {
				gFP("%s ( %d, PE.dict %s %s v.%s )",
					prefix, f.WireNumber, f.Key.Encoder, f.Encoder, f.Label)
			} else if f.Cardinality == protoreflect.Repeated {
				gFP("%s ( %d, PE.list %s v.%s )",
					prefix, f.WireNumber, f.Encoder, f.Label)
			} else {
				gFP("%s ( %d, %s v.%s )",
					prefix, f.WireNumber, f.Encoder, f.Label)
			}
			written = true
		}
		g.P("        ]")
		if len(r.Oneofs) > 0 {
			// Oneof field handling
			written = false
			for _, f := range r.Fields {
				if !f.IsOneof { // Already handled
					continue
				}
				gFP("        ++ %s v.%s", f.Encoder, f.Label)
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
		g.P("                ", u.DefaultVariant.ID, " wire ->")
		g.P("                    wire")
		for _, v := range u.Variants {
			g.P("                ", v.ID, " ->")
			g.P("                    ", v.Number)
		}
		g.P("    in")
		g.P("    PE.int32 conv")
	}
}
