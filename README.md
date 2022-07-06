# protoc-gen-elmer

A protoc code generator that produces decoders, encoders and RPC clients for Elm. A solution for [types without borders] (https://www.youtube.com/watch?v=memIRXFSNkU).

Under heavy / active development. Speak to @Joshua on Slack if needed.

## Motivation
- Elm decoders and encoders are a pain:
    - Easy to write,
    - usually long and repetitive,
    - subtleties like nullable leads to accidental footguns,
    - writing them is simply not valuable:
        - just take time
        - usually when you least want to e.g. making changes elsewhere, head full of ideas of your actual problem
- RPC clients
    - we have to define a protocol: a form of client <> server
    - the semantics of communicating with an external service are known
        - message passing to TEA
    - so ultimately: trigger (probably init) -> cmd -> ... -> msg
    - Also a boring problem
- Server side
    - Not unique to Elm, this pain is everywhere
    - unless you have the same language on front and back plus a transport
    - This sucks too!
- backwards / future compatability
    - dealing with change
    - remarkably little opinion on how to deal with this, I just want instruction

Why now? There's a good protobuf library that deals with the wire format. A couple of codegen projects but they're not feature complete and don't go far as they could. For example, well-known type support is missing and no RPC clients. This is in part because the PB browser support is poor. I think this is a missed opportunity for Elm and missing puzzle piece.

In summary: let's minimise this problem and not think about it ever again.

### What for real
- Types without borders
    - Specify an API
    - Clients use the API, servers implement it
    - and watch codegen / well-established libraries propagate from this
- Protobuf
    - IDL
    - Lots of codegen, probably covering your language (now including Elm!)
    - Specifies compatability layer and how to deal with change
    - they've thought about it! They have an answer! They have opinions on it!
- Elm has a good base protobuf codec library
    - It's binary (presumably faster),
    - it's feature complete / stays true to Protobuf
- Provide answers and clear guidelines for server-side support.

Bring all this together and write a client-side generator for:
    - codecs
    - RPC
    - tests
    - minimal, simple server side

### Examples

Showcase the ease and value-add. Given a Protobuf example:

```
show example of small Protobuf: message+enum with a few values
```

Translate it for the reader.
- messages = records
- enum = sum type
- show example Elm output

Example of RPC:

```
show Twirp hat service example
```

Again, translate: service / rpc -> elm interface

This section is TODO. For now, check the examples folder. Also see ##Development testing.

### counter-what
- RPC moves away from REST which existing tooling is focussed on. It limits observability.
- Protobufs have _opinions_. It borrows a lot from Go's semantics. Notably Go's codegen results sometimes look out of place; e.g. naming with underscores.

### So counter-what
- Codecs are the /real/ value add. Get rid of your boilerplate!
- Don't use RPC if you don't want to or it doesn't add much for you.

### Trade-offs, downsides, and limitations

- Protobuf enums are open, Elm wants them to be closed. Fundamental to Protobuf's opinions on API extensibility.
- Protobuf oneofs may be nil.

### Ecosystem

- https://github.com/eriktim/elm-protocol-buffers Underlying Protobuf transport library. It's a good foundation.
- https://github.com/andreasewering/protoc-gen-elm An alternative codegen solving the same or similar problems.
- https://github.com/tiziano88/elm-protobuf
- #elm-protobuf
- TODO: alternatives to IDL+RPC approach. Graphql?

## Installing

TODO provide GitHub binary download to $PATH

For now, see ##Development

### Using `protoc-gen-elmer`

These commands output to `examples/`.

The following generates code for our example.proto:
```
protoc --elmer_out=examples --elmer_opt="" examples/example.proto
protoc --elmer-fuzzer_out=examples --elmer-fuzzer_opt="" examples/example.proto
protoc --elmer-twirp_out=examples --elmer-twirp_opt="" examples/example.proto
```

TODO comment on how to organise .proto. Best practices, etc

### Options

TODO expand this section with options and commentary on each option

## Server-side

- Describe why PB on the server is a natural fit.
    - For Go, kind of looks like a HTTP handler anyway: `handler(ctx, inputs) (output, error)` where inputs is a mapping from req. Difficult to reduce this.
    - Great for inter-server comms (hopefully you don't have this issue)
- Mention pitfalls
    - transcoding brings caveats, the absolute behemouths known as gRPC+envoy
    - no streaming
- Twirp is just a simple answer. buf/connect is a candidate. would like streams which would probably involved websockets (this is an idea for future work).

Provide examples. Add hat making service to `examples/` to show end-to-end usage as a template / starting point. This aim is the holy grail of this library!

## Development notes

You don't need to read this section to use `protoc-gen-elmer`. Steps to set up a development environment:
- Install https://grpc.io/docs/protoc-installation/ (note just the compiler, we're explicitly avoiding GRPC)
- Install Go (1.16 min)
- Install `elm`, `elm-format` and `elm-test`
- Put `bin/` in your $PATH and see ##installing
- Run `go generate ./...` (for tests)

Tests can then be run with: `go test ./...` Most test cases specify a Protobuf file, run it through codec and fuzzer codegen, `elm-format` and finally `elm-test`. If you run an individual test case you can see the code generated in `pkg/elmgen/testdata/gen-elm/src/`

Build the `protoc-gen-elmer` binaries:
```
go build -o bin/protoc-gen-elmer cmd/protoc-gen-elmer/main.go
go build -o bin/protoc-gen-elmer-fuzzer cmd/protoc-gen-elmer-fuzzer/main.go
go build -o bin/protoc-gen-elmer-twirp cmd/protoc-gen-elmer-twirp/main.go
```

An approximate, high-level view: `.proto` (stdin) -> `protogen` (PB library used by cmd/) -> `elmgen` (core pkg in this repo) -> `gen_*.go` -> `*.elm` (stdout).

We rely on `protocolbuffers/protobuf-go` to read Protobuf. The `protogen` pkg provides helper structs to translate to Go code. So the heart of this library is `elmgen` that has a similar goal in organising ingested PBs to be consumed by Elm codegen tools. The types are specified in `elmgen.go` and the entry point is `NewModule` at the bottom. It's goal is to always produce valid Elm code (except for name collisions). You'll also find the code generators in this folder as they're coupled to the internal elmgen.

You'll find `printgen` is just a simple, dumb way to dump the structures of `protogen` to better understand inputs. It aided initial development.

Finally we also have `cmdgen`. A small pkg that holds common options and helpers for the `protoc-gen-elm*` commands.

- Understanding PB comments: https://pkg.go.dev/google.golang.org/protobuf/types/descriptorpb#SourceCodeInfo_Location

### TODO
This is my dev scratchpad of ideas and in-progress notes.

Major goals to complete:
- Existing options have caveats (e.g., partial feature support). Avoid this.
- Twirp client options (URL prefix, auth, etc)

Smaller steps:
- any TODO comments
- lazy handling on recursive structures?
- elm_package= in comments?
- add Makefile to examples/e2e
- review name collisions
- codegen gives warnings on some code:
    - gen_codec: messages with one field while the rest are optional generates a warning "consider using cons instead"
    - gen_twirp: has unused imports
- comments:
    - trailing comment (and only comment) prefixes with a blank "-- "
    - comments with a leading space translates to a double spaced comment "--  asdf"
    - Large codegen files are harder to generate. Perhaps create a comment documenting all structures?
    - Go code needs more
- structure:
    - default variant of enums has the wire number stored. This is protobuf semantics leaking and makes the variants non-uniform. Don't add an unrecognised option as this should default to the default variant
    - Twirp methods are long. Move types to separate lines for readability
    - enums need a mapper to / from string

Explore options:
- Nested messages could be nil (incl. oneof)
- Nested separator should be configurable. For example using `ꓸ` https://www.compart.com/en/unicode/U+A4F8 to achieve near-_zero_ naming collisions
- Enums have a zero field that acts as an embedded option. We might want to remove this for closed enums and rejecting the whole payload (losing compatibility).
- Revert well-known types to erk/proto's "raw" interpretation

release checklist
- Review README
- examples folder

## Bugs, other

For now, talk to @Joshua on Slack
