package elmgen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateTwirp(m *Module, g *protogen.GeneratedFile) {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}

	g.P("-- Generated by protoc-gen-elmgen. DO NOT EDIT!")
	gFP("module %sTwirp exposing (..)", m.Name)

	gFP("import %s as Data", m.Name)
	gFP("import Http")
	gFP("import Protobuf.Decode as PD")
	gFP("import Protobuf.Encode as PE")

	for _, rpc := range m.RPCs {
		// TODO: api will need to be replaced with options
		gFP("%s : (Result Http.Error Data.%s -> msg) -> String -> Data.%s -> Cmd msg",
			rpc.MethodID, rpc.Out, rpc.In)
		gFP("%s msg api data =", rpc.MethodID)
		gFP("    Http.post")
		gFP(`        { url = api ++ "/%s/%s"`,
			rpc.Service, rpc.Method)
		gFP("        , body =")
		gFP("            Data.%s data", rpc.InEncoder)
		gFP("                |> PE.encode")
		gFP(`                |> Http.bytesBody "application/protobuf"`)
		gFP("        , expect = PD.expectBytes msg Data.%s", rpc.OutDecoder)
		gFP("        }")
	}
}