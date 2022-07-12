# protoc-gen-elmer

A protoc code generator that produces decoders, encoders and a minimal RPC client in Elm. A [Protobuf](https://developers.google.com/protocol-buffers/docs/overview) solution for [types without borders](https://www.youtube.com/watch?v=memIRXFSNkU).

* [Documentation](#documentation)
    + [Motivation](#motivation)
    + [What is this project?](#what-is-this-project)
    + [Examples](#examples)
    + [Features](#features)
    + [Trade-offs, downsides, and limitations](#trade-offs-downsides-and-limitations)
    + [Ecosystem](#ecosystem)
* [Install](#install)
* [Usage](#usage)
* [Server side](#server-side)
* [Development](#development)
* [Questions, feedback and bugs](#questions-feedback-and-bugs)

## Documentation
### Motivation
In Elm, writing decoders and encoders (codecs) are a pain. They're easy to write but are long and repetitive. As the primary way of ingesting data, they're also crucial to the integrity of most Elm apps. An error in a decoder usually means the app can't progress further. 

On top of that you _must_ write them when you least want to. When your head is full of ideas on your real problem, you have to stop all momentum and build a decoder. Describing what's happening on the network should help us think through our problem. Like how Elm's type system encourages us to think through state.

We want to avoid the subtleties of `null` vs `[]` when we actually care about the semantics of "is this an empty list?" We want to focus on the bigger picture rather than the details and avoid accidental foot guns.

In a nutshell, writing codecs is not valuable: they're rote and take time.

There is also the lesser problem of describing how we talk to the network. In Elm this results in one or more `elm/http` requests. The details of these requests matter but, like writing the codecs, it's of little interest. A uniform approach to using remote resources is preferable.

This pain point isn't unique to Elm. The server also needs to deal with decoding and encoding from the network. The implementation needs to be able to evolve while not breaking existing clients.

Clients and servers are separate so we need to handle independent updates. This makes backward and forward compatibility an important consideration. We need to be able to update our API without immediately changing our clients. We also need existing clients in use to continue working rather than disrupting the user. At the same time this shouldn't have a big impact on the way we write our Elm code.

In summary: let's never write a decoder or encoder again.

### What is this project?
Network apps can define an API that a client uses and a server implements. This allows independent evolution of all three components enabling compatibility. The API specification allows us to define the same types across different languages. This is the goal outlined in the talk [types without borders](https://www.youtube.com/watch?v=memIRXFSNkU).

[Protocol Buffers](https://developers.google.com/protocol-buffers/docs/overview) allows us to specify the API as a schema. Code generators then run on this schema to define types in your language of choice. This generated code is never edited and should be reran whenever our schema changes.

Protobuf also provides a data format so that clients and servers can talk. Elm already has a [good, feature-complete library](https://package.elm-lang.org/packages/eriktim/elm-protocol-buffers/latest/) for parsing the data format. So we rely on that in our generated codecs.

This project is the code generator that translates Protobuf schemas to Elm code. A small Elm package called `Protobuf.Elmer` provides helper functions for generated codecs.

Finally, we also try to generate code that we'd write by hand and want to use. We achieve this by embracing and integrating Protobuf semantics where possible. This distinguishes out attempt from others which offers a more direct one-to-one translation. For example the well-known type timestamp becomes a `Time.Posix` in our codegen.

We generate code for:
- Types: enums and messages.
- Decoders, encoders and empty (zero) values for those types.
- Conversion to and from strings for enums.
- Fuzz tests.
- A minimal [Twirp RPC client](https://github.com/twitchtv/twirp) for non-streaming services.

Right! That's enough theory 😶‍🌫️ Let's move onto the practical 🛠️

### Examples

Let's start off with the basics. Here's a complete `.proto` describing a single message:
```protobuf
syntax = "proto3";
package My.FirstExample;
option go_package = "./.";
// Our very first Protobuf!
message MyFirstMessage {
    double my_first_float = 2;
    int32 my_favourite_number = 1;
    bool on_or_off = 3;
}
```

Every `.proto` is self-describing and generates a single Elm module. This will create the file [My.FirstExample.elm](/examples/readme/Ex01Records.elm). Here's a snippet:
```elm
{-| Our very first Protobuf!
-}
type alias MyFirstMessage =
    { myFirstFloat : Float
    , myFavouriteNumber : Int
    , onOrOff : Bool
    }
```

Some observations:
- direct mapping from the package to the Elm module name,
- we must specify something in `go_package` since we rely on Go's codegen,
- comments pass through,
- types map to expected Elm types,
- field ordering is the same as the source (not the wire number),
- naming is in camel case as expected by Elm.

You'll also find an `emptyMyFirstMessage`, `decodeMyFirstMessage` and `encodeMyFirstMessage` functions in the module. These are your codecs.

Next up are enums:
```protobuf
syntax = "proto3";
package Ex02;
option go_package = "./.";
enum Answer {
    // Look out! Name collision!
    MAYBE = 0;
    YES = 1;
    NO = 2;
}
```

Take a look at the [generated Elm type](/examples/readme/Ex02Enums.elm):
```elm
type
    Answer
    -- Look out! Name collision!
    = XMaybe
    | Yes
    | No
    | PleaseRepeat
```

As you might expect, there is a direct mapping from enums to Elm's custom types. Inline comments also pass through. You might also be able to tell in this example that we pass all generated code through `elm-format`.

We treat wire numbers as a hidden detail of data format. This means there's no final "unrecognised" option giving the wire number. The semantics of [enumerations](https://developers.google.com/protocol-buffers/docs/proto3#enum) say that unrecognised options should use the default value. The [default value of enums](https://developers.google.com/protocol-buffers/docs/proto3#default) is also the first value. So you should see enums as a `Maybe` type with the first value meaning `Nothing` and other values being part of the `Just`.

Since "MAYBE" maps to "Maybe" which would collide with the [Maybe in Elm's core library](https://package.elm-lang.org/packages/elm/core/latest/Maybe) it's prefixed with an "X". While `protoc-gen-elmer` will work around Elm naming collisions it's best to try and minimise these conflicts.

Those two examples highlight a crucial point: we expect our `.proto` files to be _designed_. By this I mean we have to take into consideration how our schema translates to Elm code. There's no way to avoid this when using another language since there won't be a direct mapping.

Finally, lots more examples can be found under [/examples](/examples):
- A showcase of a more [complex.proto](/examples/readme/03-complex.proto) and the [resulting Elm module](/examples/readme/Ex03.elm).
- A minimal [Twirp RPC client](/examples/readme/04-twirp.proto) and the [resulting Elm module](/examples/readme/Ex04.elm).
- A [real world example](/examples/real-world).
- Finally, check out the [end-to-end hat making example](/examples/end-to-end) as a quick start template.

### Features

Mapping of Protobuf features to their corresponding Elm:

| Protobuf | Elm | Default / empty / zero value ||
|---|---|---|---|
| `package` | Module name and path | n/a |
| `double`, `float` | `Float` | `0.0` |
| `int32`, `uint32`, `sint32`, `fixed32`, `sfixed32` | `Int` |
| `int64`, `uint64`, `sint64`, `fixed64`, `sfixed64` | n/a | [Not supported by the parser library](https://package.elm-lang.org/packages/eriktim/elm-protocol-buffers/latest/#known-limitations). Elm doesn't have 64-bit integer support
| `bool` | `Bool` | `False` |
| `string` | `String` | `""` |
| `bytes` | `elm/Bytes` | `[]` |
| `optional` | `Maybe ...` | Nothing | Nilable type instead of taking the default value
| `repeated` | `List ...` | `[]` | Our list type
| `required` | n/a | n/a | Proto2 option for semantics default in proto3. All types are required and take the default value if missing
| `message` | Record | `emptyRecord` function | Protobuf requires every type to have a default value
| `enum` | Custom type | First defined value | Must be `= 0;`
| Comments | Location dependent `{-\|` and `--` | n/a | An Elm document string is generated for the whole module
| `oneof` | `Maybe ...` | `Nothing` | A special, data holding, kind of enum
| `map<key, val>` | `Dict Key Val` | `Dict.empty` | The key must be a scalar type
| `Timstamp` | `Time.Posix` | Zero (1970 epoch) | Well-known type from `google/protobuf/timestamp.proto`
| Well-known types | `Google.Protobuf.*` | `Protobuf.Elmer.empty*` | Pass through to the [raw type](https://package.elm-lang.org/packages/eriktim/elm-protocol-buffers/latest/Google-Protobuf).
| `service` | n/a | n/a | Use `protoc-gen-elmer-twirp` to generate a `*Twirp.elm` RPC client.

Proto3 relies on default values, but these can be overriden when using proto2 syntax. This will override the default values specified above.

### Trade-offs, downsides, and limitations

Elm, or rather, JavaScript doesn't support 64-bit integers. You will see an error if you try to use them.

Naming collisions are resolved by prefixing with an "x" or "X".

Protobuf schemas are hierarchical with many namespaces. We try to stick to Elm naming conventions but Protobuf namespaces show up as underscores `_`. Enums are also prefixed by their enum name. This will make our codegen look a little out of place. It is done to make the Protobuf to Elm mapping clear and follow Go's tried and tested approach.

Recursive Protobuf schemas are not supported. Generated code produces [recursive aliases](https://github.com/elm/compiler/blob/master/hints/recursive-alias.md).

Protobuf enums are open which is fundamental for backward and forward compatibility. The Proto3 language requires us to [rely on default values](https://developers.google.com/protocol-buffers/docs/proto3#default) when missing. But when deserialising enumerations it asks us to keep the unrecognised option around in some fashion. It says that languages with closed enums, like Elm, should have an extra option.

The code generated by this library foregoes this and overloads the unrecognised option with the default value option. Having an unrecognised option turns it into a kind of semi-open enum. Everywhere it's used we must take into consideration this extra variant, creating another code path and adding to our code complexity.

With extra options in our enums we're forced to deal with them every time they're used. This has to be meaningful in both the update and view functions. Ideally, this would be done avoiding a catch-all pattern to ignore the unused options as these lead to bugs. If we don't reject the payload as a whole then we're forced to create more enums to cover our reduced use cases properly.

But this is a lot more work. Work that might not be useful to your app. Work that we're actively trying to avoid with codegen. So we wrap the unrecognised option into the default value (which cannot go away) and design our `.proto` files with this in mind. We then get code that's easier to work directly with.

Your opinion on this probably depends on your use case. If you come up with a situation where this doesn't well, please open an issue and share the details. Other ideas include rejecting the payload instead (dropping compatibility) and an `--elmer_opt` for adding in unrecognised options to generated code.

Nested messages are not wrapped in a `Maybe` type representing a `null`. In languages where nulls are less explicit such as Go, this is normal. For Elm it makes dealing with the code much harder but doesn't appear essential to Protobuf semantics.

If you do need nullable types then the `optional` field type is available. This will wrap any field in a `Maybe`. So will `oneof` since it needs to handle the case of no field being passed on the wire. Finally there are the well-known wrapper types which were originally used for this optionality.

Another downside includes trying to integrate with server-side technology. If you try to integrate with REST APIs then you end up having to transcode GET queries. Translating from a query string to Protobuf leads to a _lot_ of restrictions on what you can represent. Since I control both client and server and I'm not writing an open API so I chose to take the path of setting up RPC endpoints using Twirp. This is why you see a minimal Twirp client integrated into this project.

The generated Twirp client is under-developed. It's the minimum implementation required over a trusted connection.

RPC moves away from REST which has less familiar tooling and may reduce your observability. You'll also find the Protobuf RPC ecosystem is dominated by gRPC. If you're trying to minimise the scope of your projects and reduce operational complexity then you need to be careful with the technology you pick in this area.

Despite being incredibly useful, Protobuf's streaming RPC methods are not an available. Browser options to do this over HTTP are limited so you would need to rely on technology such as WebSockets. This is a natural next step for this project.

One thing to remember when evaluating this project is that *not writing codecs is the goal*. Don't use RPC if it doesn't work for you. If it comes to writing your own mapping layer, while avoiding this is preferable, it's still a much nicer problem than writing your own decoding layer.

### Ecosystem

This project is a direct alternative to [`protoc-gen-elm`](https://github.com/andreasewering/protoc-gen-elm). This project was written to try and further Protobuf support in the Elm ecosystem.

Key differences:
- Nested messages aren't wrapped in a `Maybe` making it easier to use. Use "optional" to trigger this behaviour in `protoc-gen-elmer`
- Enums don't have an unrecognised option. Use the default value (first option) as a `Nothing` value instead
- We handle imports
- [Well-known type](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf) support:
    - Timestamp uses a `Time.Posix`
    - Wrappers wrap scalars in in a `Maybe`
- Minimal (Twirp client) RPC support
- `protoc-gen-elm` is older, more established and been in use longer

Other parts of the ecosystem:
- https://github.com/eriktim/elm-protocol-buffers Underlying Protobuf binary data format library. A good foundation
- https://github.com/tiziano88/elm-protobuf Protobuf JSON library, less feature-complete
- [#elm-protobuf on Elm's Slack](https://elmlang.slack.com/messages/elm-protobuf/details/)
- Alternatives to our codegen approach such as GraphQL or sticking to handwritten JSON decoders.

## Install

Requirements:
- [The Protobuf compiler, `protoc`](https://grpc.io/docs/protoc-installation/) (for the instructions only, gRPC is not required).
- [`elm-format`](https://github.com/avh4/elm-format) in your `$PATH` unless the option `format=f` is passed.
- `elm make` for running tests.

This project is made up of three binaries: `protoc-gen-elmer`, `protoc-gen-elmer-fuzzer`, and `protoc-gen-elmer-twirp`. They all need to be available on your `$PATH` for `protoc` to work.

Copy the binaries from the Github release to `~/bin`

## Usage

Run with `protoc` against a `.proto` file. Example usage:

```
protoc
    --elmer_out=src --elmer_opt='' \
    --elmer-fuzzer_out=src --elmer-fuzzer_opt='format=f' \
    --elmer-twirp_out=src --elmer-twirp_opt='' \
    rpc/sflow/api.proto
```

See more available under [/examples](/examples).

The `--elmer_out` options trigger the plugins. Set it to your Elm `src` directory so that generated code lands in the correct location. Set options if needed with `--elmer_opt`. You can specify multiple `.proto` files and you can specify an import path with `-I`.

Each `.proto` should be self-contained. For example if you want a separate `Gen.` namespace you'll need to change the internal package name. This is critical for referencing other imports while keeping the implementation simple.

Recommendations:
- Run `protoc` commands relative to your project / repository root.
- Avoid pre or post codegen commands. If you do, make sure they're scripted.
- Have a single command to rerun all codegen. This could be `make`, a `./script` or even language-specific tools like `go generate ./...` (but this pushes the limit).
- Document this command somewhere like your README.
- Check your generated code into source control.

Options for `--elmer_opt=`:

| Option | Default | |
|---|---|---|
| format | format=t | Runs `elm-format` on generated code.

You can then send and receive in Elm with something like:
```elm
Http.request
    { url = "https://example.com/path/to/api.pb"
    , body =
        Gen.Example.encodeRequest data
            |> PE.encode
            |> Http.bytesBody "application/protobuf"
    , expect = PD.expectBytes msg Gen.Example.decodeResponse
    }
```

## Server side

To use the decoders, you will need something that responds with Protobuf. To use the encoders, you will need something that accepts Protobuf.

It may be possible to retrofit an existing API response with the content type "application/protobuf". This is the simples and easiest way to 

Ideally you'd take advantage of Protobuf's extensive codegen availability to generate server stubs and build up from there. In terms of getting on with just solving your problem this is it: write a handler that takes your well-formed inputs (plus context like DB) and return the response (or an error).

This, however, means solving the "RPC" mechanism of how to talk with it. [Twirp is the easiest solution](https://github.com/twitchtv/twirp) which this project embraces. Another upcoming solution is [Buf's Connect](https://buf.build/blog/connect-a-better-grpc).

Don't forget the [/examples/end-to-end](/examples/end-to-end) to see a complete example.

## Development

You'll need Go and you'll need to run `go generate ./...` to prepare for tests. Tests can be run with `go test ./...` Most tests run the generators on a `.proto` file. The output of these tests can be found under `/pkg/elmgen/testdata/gen-elm`

Binaries can be built with:
```
go build -o bin/protoc-gen-elmer cmd/protoc-gen-elmer/main.go
go build -o bin/protoc-gen-elmer-fuzzer cmd/protoc-gen-elmer-fuzzer/main.go
go build -o bin/protoc-gen-elmer-twirp cmd/protoc-gen-elmer-twirp/main.go
# Optionally
cp bin/protoc-gen-elmer* ~/bin
```

See also the [Makefile](/Makefile).

release checklist
- Update #Install with releases

## Questions, feedback and bugs

If you have any questions or want to give feedback, please open an issue. If you run into a bug please provide a minimal `.proto` and open an issue.
