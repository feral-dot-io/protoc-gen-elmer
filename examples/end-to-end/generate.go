package e2e

// Instead of another Makefile, this is how you might run the codegen with Go

// Go server codegen (protobuf-go and Twirp)
//go:generate protoc --go_out=go-server --twirp_out=go-server api.proto

// Elm client codegen (output of this project)
//go:generate protoc --elmer_out=elm-client/src api.proto

// Builds an RPC client
//go:generate protoc --elmer-twirp_out=elm-client/src api.proto

// For completions sake, builds test cases for decoders and encoders
//go:generate protoc --elmer-fuzzer_out=elm-client/tests api.proto
