package elmgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:generate testdata/gen-elm-test-proj

func testPlugin(t *testing.T, specs ...string) *protogen.Plugin {
	// Due to https://github.com/protocolbuffers/protobuf/issues/4163
	// We need to pass files to protoc instead of connecting stdin / stdout
	// Write all specs to files
	tmpDir := t.TempDir()
	var stdin, stdout, genReqparams string
	var filesToGen []string
	for i, spec := range specs {
		proto := "test" + strconv.Itoa(i) + ".proto"
		filesToGen = append(filesToGen, proto)
		fullProto := tmpDir + "/" + proto
		// Keep first, our main test file
		if i == 0 {
			stdin = fullProto
			stdout = fullProto + ".desc"
		}
		genReqparams += "M" + proto + "=" + proto + "/" + proto + ","
		// Write file
		err := os.WriteFile(fullProto, []byte(spec), 0644)
		assert.NoError(t, err)
	}
	// Invoke protoc's parser
	cmd := exec.Command(
		"protoc",
		"--proto_path="+tmpDir,
		"--include_imports",
		"--include_source_info",
		"--descriptor_set_out="+stdout,
		stdin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	assert.NoError(t, err)
	// Read contents of file stdout
	descriptorBytes, err := os.ReadFile(stdout)
	assert.NoError(t, err)

	// Read into pb
	protoFiles := new(descriptorpb.FileDescriptorSet)
	err = proto.Unmarshal(descriptorBytes, protoFiles)
	assert.NoError(t, err)
	// Create new codegen request
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: filesToGen,
		Parameter:      &genReqparams,
		ProtoFile:      protoFiles.File,
	}
	plugin, err := (protogen.Options{}).New(req)
	assert.NoError(t, err)
	// Build Module from File
	//assert.Len(t, plugin.Files, len(specs))
	var expFiles int
	for _, f := range protoFiles.File {
		if f.Package == nil || !strings.HasPrefix(*f.Package, importGooglePB) {
			expFiles += 1
		}
	}
	assert.Len(t, plugin.Files, expFiles)
	return plugin
}

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

var testFileContents map[string][]byte // For comment testing

func testModule(t *testing.T, specs ...string) *Module {
	plugin := testPlugin(t, specs...)
	testProjectDir := "./testdata/gen-elm"
	testFileContents = make(map[string][]byte)

	// Remove old tests
	err := os.RemoveAll(testProjectDir + "/src")
	assert.NoError(t, err)
	err = os.MkdirAll(testProjectDir+"/src", 0755)
	assert.NoError(t, err)

	var lastCodec *Module
	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}

		var elm *Module
		runGenerator := func(suffix string, gen func(m *Module, g *protogen.GeneratedFile)) {
			elm = NewModule(suffix, f)
			// Generate file
			genFile := plugin.NewGeneratedFile(elm.Path, "")
			gen(elm, genFile)
			// Always format (checks Elm syntax)
			formatted := FormatFile(plugin, elm.Path, genFile)
			content, _ := formatted.Content()
			assert.NotEmpty(t, content)
			// We generated badly formatted Elm code, write unformatted instead
			if len(content) == 0 {
				content, _ = genFile.Content()
			}
			// Ensure folder path exists
			fullFile := testProjectDir + "/src/" + elm.Path
			err = os.MkdirAll(filepath.Dir(fullFile), 0755)
			assert.NoError(t, err)
			// Copy to testdata for inspection / tests
			err = os.WriteFile(fullFile, content, 0644)
			assert.NoError(t, err)
			// Make available
			testFileContents[elm.Path] = content
		}

		// Run through all of our codegen
		runGenerator("", GenerateCodec)
		lastCodec = elm
		runGenerator("Tests", GenerateFuzzTests)
		if len(lastCodec.Services) > 0 {
			runGenerator("Twirp", GenerateTwirp)
		}
	}
	// Change pwd to tests
	wd, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(testProjectDir)
	assert.NoError(t, err)
	// Finally, run tests
	err = runElmTest(testProjectDir, "src/**/*Tests.elm", 10)
	assert.NoError(t, err)
	// Reset wd
	err = os.Chdir(wd)
	assert.NoError(t, err)
	// Run tests as local codec
	return lastCodec
}

func TestSpecialProto(t *testing.T) {
	testModule(t, `
		syntax = "proto3";
		message Emptyish {}
		message _ {
			bool _ = 1;
		}
	`)
}

func TestLocality(t *testing.T) {
	m := new(Module)
	m.importsSeen = make(map[string]bool)
	m.Name = "OurMod"
	ref1 := m.newElmRef("NotOurs", "a")
	ref2 := m.newElmRef("OurMod", "b")
	assert.Equal(t, "NotOurs", ref1.Module)
	assert.Equal(t, "", ref2.Module)
}

func TestProto2(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto2";
		message SearchRequest {
		  required string query = 1;
		  optional bool yes = 2 [default = true];
		  optional bool no = 3;
		}`)
	r := elm.Records[0]
	assert.Equal(t, "SearchRequest", r.Type.ID)
	assert.Equal(t, protoreflect.Required, r.Fields[0].Desc.Cardinality())
	assert.Equal(t, protoreflect.Optional, r.Fields[1].Desc.Cardinality())
	assert.Equal(t, protoreflect.Optional, r.Fields[2].Desc.Cardinality())
	assert.Equal(t, "True", fieldZero(elm, r.Fields[1].Desc))
	assert.Equal(t, "False", fieldZero(elm, r.Fields[2].Desc))
}

func TestQualifiedWithComments(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		package comments;
		message Outer {

			// enum comment 0
			// enum comment 1

			// enum comment 2
			// enum comment 3
			enum Option {
				// variant comment 0
				Hero = 0; // variant comment 1
				// variant comment 2

				// variant comment 3
				Worst = 1;
				Best = 2;
			}

			// message comment 0
			// message comment 1
			message Inner {
				// field comment 0
				bool a = 1; // field comment 1

				// field comment 2
				Option o = 2;
				
				// oneof comment 0
				oneof conundrum {
					// oneof comment 1
					bool or = 3; // oneof comment 2

					bool and = 4; // oneof comment 3
				};
				optional bool mayhem = 5;
			}
			Inner inner = 1;
		}
	`)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 2)
	// Union
	u := elm.Unions[0]
	assert.Equal(t, "Outer_Option", u.Type.ID)
	assert.Len(t, u.Variants, 3)
	assert.Equal(t, "Outer_Hero", u.Default().ID.ID)
	assert.Equal(t, "Outer_Hero", u.Variants[0].ID.ID)
	assert.Equal(t, "Outer_Worst", u.Variants[1].ID.ID)
	assert.Equal(t, "Outer_Best", u.Variants[2].ID.ID)
	// Records
	assert.Equal(t, "Outer", elm.Records[0].Type.ID)
	assert.Equal(t, "Outer_Inner", elm.Records[1].Type.ID)
	// Oneof
	o := elm.Oneofs[0]
	assert.Equal(t, "Outer_Inner_Conundrum", o.Type.ID)
	assert.Equal(t, "Outer_Inner_Or", o.Variants[0].ID.ID)
	assert.Equal(t, "Outer_Inner_And", o.Variants[1].ID.ID)
	assert.Equal(t, "Outer_Inner_Mayhem", elm.Oneofs[1].Type.ID)
	assert.Equal(t, "outer_Inner_MayhemDecoder", elm.Oneofs[1].Type.Decoder.ID)
	// Check comments
	content := string(testFileContents["Comments.elm"])
	for placement, max := range map[string]int{
		"enum": 3, "variant": 3, "message": 1, "field": 2, "oneof": 4,
	} {
		for i := 0; i < max; i++ {
			check := fmt.Sprintf("%s comment %d", placement, i)
			assert.True(t, strings.Contains(content, check), check)
		}
	}
}
