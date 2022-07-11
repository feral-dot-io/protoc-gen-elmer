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
	gFP("{-| Protobuf library for executing RPC methods defined in " + m.Proto + ". This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit. -}")
	printDoNotEdit(g)

	gFP("import Http")
	printImports(g, m, true)

	for _, s := range m.Services {
		s.Comments.printDashDash(g)
		for _, rpc := range s.Methods {
			// TODO: api will need to be replaced with options
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
