package elmgen

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var TestConfig = Config{
	VariantSuffixes: true}

func testPlugin(t *testing.T, raw string) *protogen.Plugin {
	t.Helper()
	// Due to https://github.com/protocolbuffers/protobuf/issues/4163
	// We need to pass files to protoc instead of connecting stdin / stdout
	// Write file stdin
	tmpDir := t.TempDir()
	stdin, stdout := tmpDir+"/stdin", tmpDir+"/stdout"
	err := os.WriteFile(stdin, []byte(raw), 0644)
	assert.NoError(t, err)
	// Invoke protoc's parser
	cmd := exec.Command(
		"protoc",
		"--proto_path="+tmpDir,
		"--descriptor_set_out="+stdout,
		stdin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	assert.NoError(t, err)
	// Read contents of file stdout
	descriptorBytes, err := os.ReadFile(tmpDir + "/stdout")
	assert.NoError(t, err)

	// Read into pb
	protoFiles := new(descriptorpb.FileDescriptorSet)
	err = proto.Unmarshal(descriptorBytes, protoFiles)
	assert.NoError(t, err)
	// Create new codegen request
	params := "Mstdin=stdin/stdin"
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"stdin"},
		Parameter:      &params,
		ProtoFile:      protoFiles.File,
	}
	plugin, err := (protogen.Options{}).New(req)
	assert.NoError(t, err)
	// Build Module from File
	assert.Len(t, plugin.Files, 1)
	return plugin
}

//go:generate testdata/gen-elm-test-proj

func (config *Config) testModule(t *testing.T, raw string) *Module {
	t.Helper()
	config.ModuleName = "Codec" // Override to run tests

	plugin := testPlugin(t, raw)
	elm, err := config.NewModule(plugin.Files[0])
	assert.NoError(t, err)

	testProjectDir := "./testdata/gen-elm"
	assertCodec := func(file string, gen func(m *Module, g *protogen.GeneratedFile)) {
		t.Helper()
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
		// Copy to testdata for inspection / tests
		err = os.WriteFile(testProjectDir+"/src/"+file, content, 0644)
		assert.NoError(t, err)
	}
	// Sanity check: always run through code gen
	assertCodec("Codec.elm", GenerateCodec)
	assertCodec("CodecTests.elm", GenerateFuzzTests)
	// Finally, run tests
	err = runElmTest(testProjectDir, "src/**/*Tests.elm", 10)
	assert.NoError(t, err)
	return elm
}

func TestProtoUnderscores(t *testing.T) {
	config := &Config{QualifyNested: true}
	config.testModule(t, `
		syntax = "proto3";
		message _ {
			bool _ = 1;
		}
	`)
}

func TestNameFromPath(t *testing.T) {
	// Normally takes from pkg
	name, path := TestConfig.nameAndPath("My.Path", "file.elm")
	assert.Equal(t, "My.Path", name)
	assert.Equal(t, "My/Path", path)
	// Mising pkg (not in source) pulls from file
	name, path = TestConfig.nameAndPath("", "my/proto_file.elm")
	assert.Equal(t, "My.ProtoFile", name)
	assert.Equal(t, "My/ProtoFile", path)
	// Module override
	config := &Config{ModuleName: "My.Override"}
	name, path = config.nameAndPath("My.Path", "file.elm")
	assert.Equal(t, "My.Override", name)
	assert.Equal(t, "My/Override", path)
}

func TestQualified(t *testing.T) {
	nestedProto := `
		syntax = "proto3";
		message Outer {
			enum Option {
				Hero = 0;
				Worst = 1;
				Best = 2;
			}
			message Inner {
				bool a = 1;
				Option o = 2;
				oneof conundrum {
					bool or = 3;
					bool and = 4;
				};
				optional bool maybe = 5;
			}
			Inner inner = 1;
		}
	`
	config := TestConfig
	config.QualifyNested = true
	config.QualifiedSeparator = "_"
	elm := config.testModule(t, nestedProto)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 2)
	// Union
	u := elm.Unions[0]
	assert.Equal(t, ElmType("Outer_Option"), u.ID)
	assert.Len(t, u.Variants, 2)
	assert.Equal(t, ElmType("Hero_Outer_Option"), u.DefaultVariant.ID)
	assert.Equal(t, ElmType("Worst_Outer_Option"), u.Variants[0].ID)
	assert.Equal(t, ElmType("Best_Outer_Option"), u.Variants[1].ID)
	// Records
	assert.Equal(t, ElmType("Outer"), elm.Records[0].ID)
	assert.Equal(t, ElmType("Outer_Inner"), elm.Records[1].ID)
	// Oneof
	o := elm.Oneofs[0]
	assert.Equal(t, ElmType("Outer_Inner_Conundrum"), o.ID)
	assert.Equal(t, ElmType("Or_Outer_Inner_Conundrum"), o.Variants[0].ID)
	assert.Equal(t, ElmType("And_Outer_Inner_Conundrum"), o.Variants[1].ID)
	assert.Equal(t, ElmType("Outer_Inner_Maybe"), elm.Oneofs[1].ID)
	assert.Equal(t, "outer_Inner_MaybeDecoder", elm.Oneofs[1].DecodeID)

	// Again but without qualifying
	config.QualifyNested = false
	elm = config.testModule(t, nestedProto)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 2)
	// Union
	u = elm.Unions[0]
	assert.Equal(t, ElmType("Option"), u.ID)
	assert.Len(t, u.Variants, 2)
	assert.Equal(t, ElmType("Hero_Option"), u.DefaultVariant.ID)
	assert.Equal(t, ElmType("Worst_Option"), u.Variants[0].ID)
	assert.Equal(t, ElmType("Best_Option"), u.Variants[1].ID)
	// Records
	assert.Equal(t, ElmType("Inner"), elm.Records[0].ID)
	assert.Equal(t, ElmType("Outer"), elm.Records[1].ID)
	// Oneof
	o = elm.Oneofs[0]
	assert.Equal(t, ElmType("Inner_Conundrum"), o.ID)
	assert.Equal(t, ElmType("Or_Conundrum"), o.Variants[0].ID)
	assert.Equal(t, ElmType("And_Conundrum"), o.Variants[1].ID)

	// With invalid separator
	_, err := (&Config{QualifiedSeparator: " "}).NewModule(nil)
	assert.ErrorContains(t, err, "separator")
}
