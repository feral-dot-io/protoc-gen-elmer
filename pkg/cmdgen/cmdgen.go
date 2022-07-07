package cmdgen

import (
	"flag"

	"github.com/feral-dot-io/protoc-gen-elmer/pkg/elmgen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	format = flag.Bool("format", true,
		"Runs generated source code through elm-format.")
)

type Generator func(*elmgen.Module, *protogen.GeneratedFile)

// Creates a function that runs the given generator over all of a plugin's files to be generated. Applies options from global flags. The suffix is intended to identify the outputted files from the generator.
func RunGenerator(suffix string, generator Generator) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			// Map Proto to Elm types
			elm := elmgen.NewModule(suffix, file)
			// Write to file
			genFile := plugin.NewGeneratedFile(elm.Path, "")
			generator(elm, genFile)
			// Format file?
			if *format {
				elmgen.FormatFile(plugin, elm.Path, genFile)
			}
		}
		return nil
	}
}
