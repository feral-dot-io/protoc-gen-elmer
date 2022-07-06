package elmgen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateTwirp(m *Module, g *protogen.GeneratedFile) {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}

	gFP("module %s exposing (..)", m.Name)
	printDoNotEdit(g)

	gFP("import Http")
	printImports(g, m, true)

	for _, s := range m.Services {
		s.Comments.printDashDash(g)
		for _, rpc := range s.Methods {
			// TODO: api will need to be replaced with options
			rpc.Comments.printBlock(g)
			gFP("%s : (Result Http.Error %s -> msg) -> String -> %s -> Cmd msg",
				rpc.ID.ID, rpc.Out, rpc.In)
			gFP("%s msg api data =", rpc.ID.ID)
			gFP("    Http.post")
			gFP(`        { url = api ++ "/%s/%s"`, rpc.Service, rpc.Method)
			gFP("        , body =")
			gFP("            %s data", rpc.In.Encoder)
			gFP("                |> PE.encode")
			gFP(`                |> Http.bytesBody "application/protobuf"`)
			gFP("        , expect = PD.expectBytes msg %s", rpc.Out.Decoder)
			gFP("        }")
			rpc.Comments.printBlockTrailing(g)
		}
		g.P(s.Comments.Trailing)
	}
}
