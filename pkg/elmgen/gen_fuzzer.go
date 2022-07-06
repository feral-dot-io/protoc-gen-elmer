package elmgen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func runElmTest(projDir, globs string, fuzz int) error {
	cmd := exec.Command("elm-test")
	if fuzz > 0 {
		cmd.Args = append(cmd.Args, "--fuzz", strconv.Itoa(fuzz))
	}
	if globs != "" {
		cmd.Args = append(cmd.Args, globs)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func GenerateFuzzTests(m *Module, g *protogen.GeneratedFile) {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}

	g.P("module ", m.Name, " exposing (..)")
	printDoNotEdit(g)

	g.P("import Expect")
	g.P("import Fuzz exposing (Fuzzer)")
	g.P("import Test exposing (Test, fuzz, test)")
	printImports(g, m, false)

	// Union fuzzers
	for _, u := range m.Unions {
		t := u.Type
		gFP("%s : Fuzzer %s", t.Fuzzer.ID, t)
		gFP("%s =", t.Fuzzer.ID)
		gFP("    Fuzz.oneOf")
		gFP("        [ Fuzz.map %s %s.fuzzInt32", u.DefaultVariant.ID, importElmerTests)
		for _, v := range u.Variants {
			gFP("        , Fuzz.constant %s", v.ID)
		}
		gFP("        ]")
	}

	// Record fuzzers
	for _, r := range m.Records {
		gFP("%s : Fuzzer %s", r.Type.Fuzzer.ID, r.Type)
		gFP("%s =", r.Type.Fuzzer.ID)
		if len(r.Oneofs) > 0 {
			gFP("    let")
			for _, o := range r.Oneofs {
				gFP("        %s =", o.Type.Fuzzer.ID)
				gFP("            Fuzz.oneOf")
				for j, v := range o.Variants {
					prefix := "                "
					if j == 0 {
						prefix += "["
					} else {
						prefix += ","
					}
					fuzzer := fieldFuzzer(m, v.Field.Desc)
					if o.IsSynthetic { // Optional field
						gFP("%s %s", prefix, fuzzer)
					} else {
						gFP("%s Fuzz.map %s %s", prefix, v.ID, fuzzer)
					}
				}
				g.P("                ]")
			}
			gFP("    in")
		}

		if len(r.Fields) == 0 {
			gFP("    Fuzz.constant %s", r.Type)
		} else {
			gFP("    Fuzz.map %s", r.Type)
		}
		for i, f := range r.Fields {
			spacer := "        "
			prefix := spacer
			if i != 0 {
				prefix += "|> Fuzz.andMap "
			}
			if f.Oneof != nil {
				gFP("%s(Fuzz.maybe %s)", prefix, f.Oneof.Type.Fuzzer)
			} else {
				gFP("%s%s", prefix, fieldFuzzer(m, f.Desc))
			}
		}
	}

	// Test cases
	var types []*ElmType
	for _, r := range m.Records {
		types = append(types, r.Type)
	}
	for _, u := range m.Unions {
		types = append(types, u.Type)
	}
	for _, t := range types {
		gFP("test%s : Test", t.ID)
		gFP("test%s =", t.ID)
		gFP("    let")
		gFP("        run = %s.runTest %s %s", importElmerTests, t.Decoder, t.Encoder)
		gFP("    in")
		gFP(`    Test.describe "encode then decode %s"`, t.ID)
		gFP(`        [ test "empty" (\_ -> run %s)`, t.Zero)
		gFP(`        , fuzz %s "fuzzer" run`, t.Fuzzer.ID)
		gFP("        ]")
	}
}

func fieldFuzzer(m *Module, fd protoreflect.FieldDescriptor) string {
	if fd.IsMap() {
		key := fieldFuzzer(m, fd.MapKey())
		val := fieldFuzzer(m, fd.MapValue())
		return "(Fuzz.map Dict.fromList (Fuzz.list (Fuzz.tuple (" + key + ", " + val + "))))"
	} else if fd.IsList() {
		return "(Fuzz.list " + fieldFuzzerKind(m, fd) + ")"
	}
	return fieldFuzzerKind(m, fd)
}

func fieldFuzzerKind(m *Module, fd protoreflect.FieldDescriptor) string {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return "Fuzz.bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return importElmerTests + ".fuzzInt32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return importElmerTests + ".fuzzUInt32"

	/* Unsupported by Elm / JS
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64", nil
	*/

	case protoreflect.FloatKind:
		return importElmerTests + ".fuzzFloat32"
	case protoreflect.DoubleKind:
		return "Fuzz.float"

	case protoreflect.StringKind:
		return "Fuzz.string"
	case protoreflect.BytesKind:
		return importElmerTests + ".fuzzBytes"

	case protoreflect.EnumKind:
		ed := fd.Enum()
		return m.NewElmType(ed.ParentFile(), ed).Fuzzer.String()

	case protoreflect.MessageKind, protoreflect.GroupKind:
		md := fd.Message()
		return m.NewElmType(md.ParentFile(), md).Fuzzer.String()
	}

	log.Panicf("kindFuzzer: unknown protoreflect.Kind: %s", fd.Kind())
	return ""
}
