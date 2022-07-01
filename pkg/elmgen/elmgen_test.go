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

func testPlugin(t *testing.T, specs ...string) *protogen.Plugin {
	t.Helper()
	// Due to https://github.com/protocolbuffers/protobuf/issues/4163
	// We need to pass files to protoc instead of connecting stdin / stdout
	// Write all specs to files
	tmpDir := t.TempDir()
	var stdin, stdout, genReqparams string
	for i, spec := range specs {
		proto := "test" + strconv.Itoa(i) + ".proto"
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
		FileToGenerate: []string{"test0.proto"},
		Parameter:      &genReqparams,
		ProtoFile:      protoFiles.File,
	}
	plugin, err := (protogen.Options{}).New(req)
	assert.NoError(t, err)
	// Build Module from File
	assert.Len(t, plugin.Files, len(specs))
	return plugin
}

//go:generate testdata/gen-elm-test-proj

var changedTestDir bool
var testFileContents map[string][]byte // For comment testing

func testModule(t *testing.T, specs ...string) *Module {
	t.Helper()
	plugin := testPlugin(t, specs...)
	var elm *Module
	// One file should have Generate=true. Run NewModule on it
	for _, f := range plugin.Files {
		if f.Generate {
			if elm != nil {
				t.Fail()
			}
			var err error
			elm = NewModule("", f)
			assert.NoError(t, err)
		}
	}

	testProjectDir := "./testdata/gen-elm"
	// Remove old tests
	err := os.RemoveAll(testProjectDir + "/src")
	assert.NoError(t, err)
	err = os.MkdirAll(testProjectDir+"/src", 0755)
	assert.NoError(t, err)

	testFileContents = make(map[string][]byte)
	assertCodec := func(suffix string, gen func(m *Module, g *protogen.GeneratedFile)) {
		t.Helper()
		file := elm.Path + suffix + ".elm"
		genFile := plugin.NewGeneratedFile(file, "")
		gen(elm, genFile)
		// Always format (checks Elm syntax)
		formatted := FormatFile(plugin, file, genFile)
		content, _ := formatted.Content()
		assert.NotEmpty(t, content)
		// We generated badly formatted Elm code, write unformatted instead
		if len(content) == 0 {
			content, _ = genFile.Content()
		}
		// Ensure folder path exists
		fullFile := testProjectDir + "/src/" + file
		err = os.MkdirAll(filepath.Dir(fullFile), 0755)
		assert.NoError(t, err)
		// Copy to testdata for inspection / tests
		err = os.WriteFile(fullFile, content, 0644)
		assert.NoError(t, err)
		// Make available
		testFileContents[file] = content
	}
	// Sanity check: always run through code gen
	assertCodec("", GenerateCodec)
	assertCodec("Tests", GenerateFuzzTests)
	if len(elm.Services) > 0 {
		assertCodec("Twirp", GenerateTwirp)
	}
	// Finally, run tests
	if !changedTestDir {
		err = os.Chdir(testProjectDir)
		assert.NoError(t, err)
		changedTestDir = true
	}
	err = runElmTest(testProjectDir, "src/**/*Tests.elm", 10)
	assert.NoError(t, err)
	return elm
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

func TestProto2(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto2";
		message SearchRequest {
		  required string query = 1;
		  optional bool yes = 2 [default = true];
		  optional bool no = 3;
		}`)
	r := elm.Records[0]
	assert.Equal(t, "SearchRequest", r.Type.Local())
	assert.Equal(t, protoreflect.Required, r.Fields[0].Cardinality)
	assert.Equal(t, protoreflect.Optional, r.Fields[1].Cardinality)
	assert.Equal(t, protoreflect.Optional, r.Fields[2].Cardinality)
	assert.Equal(t, "True", r.Fields[1].Zero)
	assert.Equal(t, "False", r.Fields[2].Zero)
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
				optional bool maybe = 5;
			}
			Inner inner = 1;
		}
	`)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 2)
	// Union
	u := elm.Unions[0]
	assert.Equal(t, "Outer_Option", u.Type.Local())
	assert.Len(t, u.Variants, 2)
	assert.Equal(t, "Outer_Hero", u.DefaultVariant.ID.Local())
	assert.Equal(t, "Outer_Worst", u.Variants[0].ID.Local())
	assert.Equal(t, "Outer_Best", u.Variants[1].ID.Local())
	// Records
	assert.Equal(t, "Outer", elm.Records[0].Type.Local())
	assert.Equal(t, "Outer_Inner", elm.Records[1].Type.Local())
	// Oneof
	o := elm.Oneofs[0]
	assert.Equal(t, "Outer_Inner_Conundrum", o.Type.Local())
	assert.Equal(t, "Outer_Inner_Or", o.Variants[0].ID.Local())
	assert.Equal(t, "Outer_Inner_And", o.Variants[1].ID.Local())
	assert.Equal(t, "Outer_Inner_Maybe", elm.Oneofs[1].Type.Local())
	assert.Equal(t, "outer_Inner_MaybeDecoder", elm.Oneofs[1].Type.Decoder().Local())
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

func _TestImports(t *testing.T) {
	elm := testModule(t, `
		syntax = "proto3";
		import "test1.proto";
		message MyMessage {
			test1.Other out_of_this_world = 1;
		}`, `
		syntax = "proto3";
		package test1;
		message Other {
			int32 a = 1;
			int32 b = 2;
			int32 c = 3;
		}`)
	assert.Len(t, elm.Records, 1)
	assert.Equal(t, "MyMessage", elm.Records[0].Type.Local())
	f := elm.Records[0].Fields[0]
	assert.Equal(t, "out_of_this_world", f.Label)
	assert.Equal(t, "Other", f.Type)
	// TODO: this requires a lot of work on naming
}
