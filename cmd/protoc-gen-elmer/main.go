package main

import (
	"flag"

	"github.com/feral-dot-io/protoc-gen-elmer/pkg/cmdgen"
	"github.com/feral-dot-io/protoc-gen-elmer/pkg/elmgen"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	opts := protogen.Options{
		ParamFunc: flag.CommandLine.Set}
	opts.Run(cmdgen.RunGenerator("", elmgen.GenerateCodec))
}
