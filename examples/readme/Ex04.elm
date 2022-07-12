module Ex04 exposing (..)

{-| Protobuf library for decoding and encoding structures found in 04-twirp.proto along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - Response

Unions: (none)

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs Response


# Empty (zero values)

@docs emptyResponse


# Decoders

@docs decodeResponse


# Encoders

@docs encodeResponse

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Google.Protobuf
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE


type alias Response =
    { message : String
    }


emptyResponse : Response
emptyResponse =
    Response ""


decodeResponse : PD.Decoder Response
decodeResponse =
    PD.message emptyResponse
        [ PD.optional 1 PD.string (\v m -> { m | message = v })
        ]


encodeResponse : Response -> PE.Encoder
encodeResponse v =
    PE.message <|
        [ ( 1, PE.string v.message )
        ]
