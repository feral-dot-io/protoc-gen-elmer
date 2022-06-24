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

func (config *Config) testModule(t *testing.T, raw string) *Module {
	t.Helper()
	plugin := testPlugin(t, raw)
	elm, err := config.NewModule(plugin.Files[0])
	assert.NoError(t, err)

	// Sanity check code: always run through code gen and formatter
	path := os.TempDir() + "/Codec.elm"
	genFile := plugin.NewGeneratedFile(path, "")
	GenerateCodec(elm, genFile)
	formatted := FormatFile(plugin, path, genFile)
	content, _ := formatted.Content()
	assert.NotEmpty(t, content)
	/** /
	unformatted, _ := genFile.Content()
	fmt.Printf("Generated code:\n\n%s\n\n", unformatted)
	fmt.Printf("Formatted code:\n\n%s\n\n", content)
	/**/
	return elm
}

func TestModuleName(t *testing.T) {
	basic := `
		syntax = "proto3";
		package test.something;
	`
	elm := TestConfig.testModule(t, basic)
	assert.Equal(t, "Test.Something", elm.Name)
	assert.Equal(t, "Test/Something", elm.Path)
	assert.False(t, elm.Imports.Bytes)
	assert.False(t, elm.Imports.Dict)
	assert.Empty(t, elm.Unions)
	assert.Empty(t, elm.Records)
	// With module name override
	config := TestConfig
	config.ModulePrefix = "Ignored."
	config.ModuleName = "My.Override"
	elm = config.testModule(t, basic)
	assert.Equal(t, "My.Override", elm.Name)
	assert.Equal(t, "My/Override", elm.Path)
}

func TestEmptyPackage(t *testing.T) {
	elm := TestConfig.testModule(t, `
		syntax = "proto3";
		message A {
			message B {
				bool inner = 1;
			}
			B b = 1;
		}
	`)
	assert.Equal(t, "Stdin.Stdin.Stdin", elm.Name)
	assert.Equal(t, "Stdin/Stdin/Stdin", elm.Path)
	assert.Empty(t, elm.Unions)
	assert.Len(t, elm.Records, 2)
	assert.Equal(t, ElmType("A"), elm.Records[0].ID)
	assert.Equal(t, ElmType("AB"), elm.Records[1].ID)
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
func TestQualified(t *testing.T) {
	config := TestConfig
	config.QualifyNested = true
	config.QualifiedSeparator = "_"
	elm := config.testModule(t, `
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
			}
			Inner inner = 1;
		}
	`)
	assert.Len(t, elm.Unions, 1)
	assert.Len(t, elm.Records, 2)
	assert.Len(t, elm.Oneofs, 1)
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
	// Invalid separator
	_, err := (&Config{QualifiedSeparator: " "}).NewModule(nil)
	assert.ErrorContains(t, err, "separator")
}
