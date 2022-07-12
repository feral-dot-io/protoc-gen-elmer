module Ex01Records exposing (..)

{-| Protobuf library for decoding and encoding structures found in 01-records.proto along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - MyFirstMessage

Unions: (none)

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs MyFirstMessage


# Empty (zero values)

@docs emptyMyFirstMessage


# Decoders

@docs decodeMyFirstMessage


# Encoders

@docs encodeMyFirstMessage

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Protobuf.Decode as PD
import Protobuf.Encode as PE


{-| Our very first Protobuf!
-}
type alias MyFirstMessage =
    { myFirstFloat : Float
    , myFavouriteNumber : Int
    , onOrOff : Bool
    }


emptyMyFirstMessage : MyFirstMessage
emptyMyFirstMessage =
    MyFirstMessage 0 0 False


decodeMyFirstMessage : PD.Decoder MyFirstMessage
decodeMyFirstMessage =
    PD.message emptyMyFirstMessage
        [ PD.optional 2 PD.double (\v m -> { m | myFirstFloat = v })
        , PD.optional 1 PD.int32 (\v m -> { m | myFavouriteNumber = v })
        , PD.optional 3 PD.bool (\v m -> { m | onOrOff = v })
        ]


encodeMyFirstMessage : MyFirstMessage -> PE.Encoder
encodeMyFirstMessage v =
    PE.message <|
        [ ( 2, PE.double v.myFirstFloat )
        , ( 1, PE.int32 v.myFavouriteNumber )
        , ( 3, PE.bool v.onOrOff )
        ]
