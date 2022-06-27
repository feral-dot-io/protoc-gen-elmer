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

	flagFilePrefix = flag.String("file_prefix", "",
		"Prefix on where to place generated files. Should be your `elm-project/src` directory.")
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
	// TODO: turn this into a flagset and reuse across other cmds?
)

type Flags struct {
	Format     bool
	FilePrefix string
	Config     *elmgen.Config
}

func NewFlags() *Flags {
	// If filePrefix is set, it must always end in a /
	if *flagFilePrefix != "" {
		*flagFilePrefix += "/"
	}
	return &Flags{
		*format,
		*flagFilePrefix,
		&elmgen.Config{
			ModulePrefix: *modulePrefix,
			ModuleName:   *moduleName,

			QualifyNested:      *qualifyNested,
			QualifiedSeparator: *qualifiedSeparator,
			VariantSuffixes:    *variantSuffixes,
			CollisionSuffix:    *collisionSuffix}}
}

type Generator func(*elmgen.Module, *protogen.GeneratedFile)

func RunGenerator(suffix string, generator Generator) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		flags := NewFlags()

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			// Map Proto to Elm types
			elm, err := flags.Config.NewModule(file)
			plugin.Error(err)
			// Write to file
			path := flags.FilePrefix + elm.Path + suffix
			genFile := plugin.NewGeneratedFile(path, "")
			generator(elm, genFile)
			// Format file?
			if flags.Format {
				elmgen.FormatFile(plugin, path, genFile)
			}
		}
		return nil
	}
}
