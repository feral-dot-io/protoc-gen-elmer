package e2e

// Go server codegen (protobuf-go and Twirp)
//go:generate protoc --go_out=go-server --twirp_out=go-server api.proto

// Elm client codegen (this project)
//go:generate protoc --elmer_out=elm-client/src --elmer_opt "qualify=t,separator=_,collision=_" api.proto
// Builds an RPC client
//go:generate protoc --elm-twirp_out=elm-client/src --elm-twirp_opt "qualify=t,separator=_,collision=_" api.proto
// For completions sake, builds test cases for decoders and encoders
//go:generate protoc --elm-fuzzer_out=elm-client/tests --elm-fuzzer_opt "qualify=t,separator=_,collision=_" api.proto