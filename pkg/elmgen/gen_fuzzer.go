package elmgen

import (
	"fmt"
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

	g.P("module ", m.Name, "Tests exposing (..)")
	printDoNotEdit(g)

	gFP("import %s", m.Name)
	if m.Imports.Bytes {
		g.P("import Bytes exposing (Bytes)")
		g.P("import Bytes.Encode as BE")
	}
	if m.Imports.Dict {
		g.P("import Dict")
	}
	g.P("import Expect")
	g.P("import Fuzz exposing (Fuzzer)")
	g.P("import Protobuf.Decode as PD")
	g.P("import Protobuf.Encode as PE")
	g.P("import Test exposing (Test, fuzz, test)")

	// Helpers
	if m.Fuzzers.Int32 || m.Fuzzers.Float32 || len(m.Unions) > 0 {
		g.P("fuzzInt32 : Fuzzer Int")
		g.P("fuzzInt32 =")
		g.P("    Fuzz.intRange -2147483648 2147483647")
	}
	if m.Fuzzers.Uint32 {
		g.P("fuzzUint32 : Fuzzer Int")
		g.P("fuzzUint32 =")
		g.P("    Fuzz.intRange 0 4294967295")
	}
	if m.Fuzzers.Float32 {
		// Avoid trying to robusly map float64 (JS) -> float32
		// Only tests exponent (float32 has 8 bits)
		g.P("fuzzFloat32 : Fuzzer Float")
		g.P("fuzzFloat32 =")
		g.P("    Fuzz.map (\\i -> 2 ^ toFloat i) fuzzInt32")
	}
	if m.Imports.Bytes {
		g.P("fuzzBytes : Fuzzer Bytes")
		g.P("fuzzBytes =")
		g.P("    Fuzz.intRange 0 255")
		g.P("        |> Fuzz.map BE.unsignedInt8")
		g.P("        |> Fuzz.list")
		g.P("        |> Fuzz.map (BE.sequence >> BE.encode)")
	}

	// Union fuzzers
	for _, u := range m.Unions {
		t := u.Type
		gFP("%s : Fuzzer %s", t.Fuzzer().Local(), t)
		gFP("%s =", t.Fuzzer().Local())
		gFP("    Fuzz.oneOf")
		gFP("        [ Fuzz.map %s fuzzInt32", u.DefaultVariant.ID)
		for _, v := range u.Variants {
			gFP("        , Fuzz.constant %s", v.ID)
		}
		gFP("        ]")
	}

	// Record fuzzers
	for _, r := range m.Records {
		gFP("%s : Fuzzer %s", r.Type.Fuzzer().Local(), r.Type)
		gFP("%s =", r.Type.Fuzzer().Local())
		if len(r.Oneofs) > 0 {
			gFP("    let")
			for _, o := range r.Oneofs {
				gFP("        %s =", o.Type.Fuzzer().Local())
				gFP("            Fuzz.oneOf")
				for j, v := range o.Variants {
					prefix := "                "
					if j == 0 {
						prefix += "["
					} else {
						prefix += ","
					}
					if o.IsSynthetic { // Optional field
						gFP("%s %s", prefix, v.Field.Fuzzer)
					} else {
						gFP("%s Fuzz.map %s %s", prefix, v.ID, v.Field.Fuzzer)
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
			if f.IsOneof {
				gFP("%s(Fuzz.maybe %s)", prefix, f.Fuzzer)
			} else if f.IsMap {
				gFP("%s(Fuzz.map Dict.fromList", prefix)
				gFP("%s    (Fuzz.list (Fuzz.tuple (%s, %s))))",
					spacer, f.Key.Fuzzer, f.Fuzzer)
			} else if f.Cardinality == protoreflect.Repeated {
				gFP("%s(Fuzz.list %s)", prefix, f.Fuzzer)
			} else { // No special treatment
				gFP("%s%s", prefix, f.Fuzzer)
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
		gFP("test%s : Test", t.Local())
		gFP("test%s =", t.Local())
		gFP("    let")
		// TODO move `run` to top-level
		gFP("        run data =")
		gFP("            PE.encode (%s data)", t.Encoder())
		gFP("                |> PD.decode %s", t.Decoder())
		gFP("                |> Expect.equal (Just data)")
		gFP("    in")
		gFP(`    Test.describe "encode then decode %s"`, t.ID)
		gFP(`        [ test "empty" (\_ -> run %s)`, t.Zero())
		gFP(`        , fuzz %s "fuzzer" run`, t.Fuzzer().Local())
		gFP("        ]")
	}
}
