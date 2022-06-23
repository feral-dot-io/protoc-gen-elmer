package main

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
		"Prefix to place generated files.")
	modulePrefix = flag.String("module_prefix", "",
		"Literal prefix for generated Elm module. For example `Gen.` becomes `Gen.My.Module`.")
	moduleName = flag.String("module", "",
		"Overrides the module name derived from the Proto package. Used in lieu of a google.protobuf.FileOptions entry. Ignores the module_prefix option. Should be avoided where possible.")

	qualifyNested = flag.Bool("qualify", elmgen.DefaultConfig.QualifyNested,
		"When dealing with nested Protobuf we can choose to fully qualify them or not. For example `Message.Enum` becomes `MessageEnum` or just `Enum`.")
	qualifiedSeparator = flag.String("separator", elmgen.DefaultConfig.QualifiedSeparator,
		"Use a separator when transforming nested Protobuf to an Elm ID. For example `Nested.Message` becomes `Nested_Message` in Elm with `_`.")
	variantSuffixes = flag.Bool("variant_suffix", elmgen.DefaultConfig.VariantSuffixes,
		"Suffixes Elm union variants with their parent's name. For example `enum Role { ... Actor ... }` becomes `ActorRole")
	collisionSuffix = flag.String("collision", elmgen.DefaultConfig.CollisionSuffix,
		"Suffix applied to an Elm ID to resolve collision. If empty, returns an error instead.")
)

func main() {
	opts := protogen.Options{
		ParamFunc: flag.CommandLine.Set}
	opts.Run(run)
}

func run(gen *protogen.Plugin) error {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	// If filePrefix is set, it must always end in a /
	if *filePrefix != "" {
		*filePrefix += "/"
	}

	config := elmgen.Config{
		ModulePrefix: *modulePrefix,
		ModuleName:   *moduleName,

		QualifyNested:      *qualifyNested,
		QualifiedSeparator: *qualifiedSeparator,
		VariantSuffixes:    *variantSuffixes,
		CollisionSuffix:    *collisionSuffix,
	}

	for _, file := range gen.Files {
		if !file.Generate {
			continue
		}
		// Map Proto to Elm types
		elm, err := config.NewModule(file)
		gen.Error(err)
		// Write to file
		path := *filePrefix + elm.Path + ".elm"
		genFile := gen.NewGeneratedFile(path, "")
		elmgen.GenerateCodec(elm, genFile)
		// Format file?
		if *format {
			elmgen.FormatFile(gen, path, genFile)
		}
	}
	return nil
}
