// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
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

type Generator func(*elmgen.Module, *protogen.GeneratedFile) bool

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
			valid := generator(elm, genFile)
			if valid {
				// Format file?
				if *format {
					elmgen.FormatFile(plugin, elm.Path, genFile)
				}
			} else {
				genFile.Skip()
			}
		}
		return nil
	}
}
