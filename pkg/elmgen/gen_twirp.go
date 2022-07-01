package elmgen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateTwirp(m *Module, g *protogen.GeneratedFile) {
	gFP := func(formatter string, args ...interface{}) {
		g.P(fmt.Sprintf(formatter, args...))
	}

	gFP("module %sTwirp exposing (..)", m.Name)
	printDoNotEdit(g)

	gFP("import Http")
	gFP("import Protobuf.Decode as PD")
	gFP("import Protobuf.Encode as PE")
	// Import dependencies
	seen := make(map[string]bool)
	importer := func(ref *ElmRef) {
		if !seen[ref.Module] {
			seen[ref.Module] = true
			gFP("import %s", ref.Module)
		}
	}
	for _, s := range m.Services {
		for _, rpc := range s.Methods {
			importer(&rpc.In.ElmRef)
			importer(&rpc.Out.ElmRef)
		}
	}

	for _, s := range m.Services {
		s.Comments.printDashDash(g)
		for _, rpc := range s.Methods {
			// TODO: api will need to be replaced with options
			rpc.Comments.printBlock(g)
			gFP("%s : (Result Http.Error %s -> msg) -> String -> %s -> Cmd msg",
				rpc.ID.Local(), rpc.Out, rpc.In)
			gFP("%s msg api data =", rpc.ID.Local())
			gFP("    Http.post")
			gFP(`        { url = api ++ "/%s/%s"`, rpc.Service, rpc.Method)
			gFP("        , body =")
			gFP("            %s data", rpc.In.Encoder())
			gFP("                |> PE.encode")
			gFP(`                |> Http.bytesBody "application/protobuf"`)
			gFP("        , expect = PD.expectBytes msg %s", rpc.Out.Decoder())
			gFP("        }")
			rpc.Comments.printBlockTrailing(g)
		}
	}
}
