// This file is part of protoc-gen-elmer.
//
// Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Protoc-gen-elmer. If not, see <https://www.gnu.org/licenses/>.
package elmgen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateTwirp(m *Module, g *protogen.GeneratedFile) bool {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}

	gFP("module %s exposing (..)", m.Name)
	gFP("{-| Protobuf library for executing RPC methods defined in package `" + m.ProtoPackage + "`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit. -}")
	printDoNotEdit(g)

	gFP("import Http")
	printImports(g, m, true)

	for _, s := range m.Services {
		s.Comments.printDashDash(g)
		for _, rpc := range s.Methods {
			rpc.Comments.printBlock(g)
			gFP("%s : (Result Http.Error %s -> msg)\n -> String -> %s -> Cmd msg",
				rpc.ID.ID, rpc.Out, rpc.In)
			gFP("%s msg api data =", rpc.ID.ID)
			gFP("    Http.riskyRequest")
			gFP(`        { method = "POST"`)
			gFP(`        , headers = []`)
			gFP(`        , url = api ++ "/%s/%s"`, rpc.Service, rpc.Method)
			gFP("        , body =")
			gFP("            %s data", rpc.In.Encoder)
			gFP("                |> PE.encode")
			gFP(`                |> Http.bytesBody "application/protobuf"`)
			gFP("        , expect = PD.expectBytes msg %s", rpc.Out.Decoder)
			gFP(`        , timeout = Nothing`)
			gFP(`        , tracker = Nothing`)
			gFP("        }")
			rpc.Comments.printBlockTrailing(g)
		}
		g.P(s.Comments.Trailing)
	}

	return len(m.Services) > 0
}
