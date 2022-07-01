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

	filePrefix = flag.String("file_prefix", "",
		"Prefix on where to place generated files. Should be your `elm-project/src` directory.")
	modulePrefix = flag.String("module_prefix", "",
		"Literal prefix for generated Elm module. For example `Gen.` becomes `Gen.My.Module`.")
)

type Generator func(*elmgen.Module, *protogen.GeneratedFile)

func RunGenerator(suffix string, generator Generator) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		// If filePrefix is set, it must always end in a /
		if *filePrefix != "" {
			*filePrefix += "/"
		}

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			// Map Proto to Elm types
			elm := elmgen.NewModule(*modulePrefix, file.Desc)
			// Write to file
			path := *filePrefix + elm.Path + suffix
			genFile := plugin.NewGeneratedFile(path, "")
			generator(elm, genFile)
			// Format file?
			if *format {
				elmgen.FormatFile(plugin, path, genFile)
			}
		}
		return nil
	}
}
