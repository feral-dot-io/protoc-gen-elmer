module Ex02Enums exposing (..)

{-| Protobuf library for decoding and encoding structures found in package `Ex02Enums` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records: (none)

Unions:

  - Answer

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs Answer


# Empty (zero values)

@docs emptyAnswer


# Enum valuesOf

@docs valuesOfAnswer


# Enum and String converters

@docs fromAnswer, toAnswer


# Decoders

@docs decodeAnswer


# Encoders

@docs encodeAnswer

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Protobuf.Decode as PD
import Protobuf.Encode as PE


type
    Answer
    -- Look out! Name collision!
    = XMaybe
    | Yes
    | No
    | PleaseRepeat


emptyAnswer : Answer
emptyAnswer =
    XMaybe


valuesOfAnswer : List Answer
valuesOfAnswer =
    [ XMaybe, Yes, No, PleaseRepeat ]


fromAnswer : Answer -> String
fromAnswer u =
    case u of
        XMaybe ->
            "MAYBE"

        Yes ->
            "YES"

        No ->
            "NO"

        PleaseRepeat ->
            "PLEASE_REPEAT"


toAnswer : String -> Answer
toAnswer str =
    case str of
        "MAYBE" ->
            XMaybe

        "YES" ->
            Yes

        "NO" ->
            No

        "PLEASE_REPEAT" ->
            PleaseRepeat

        _ ->
            XMaybe


decodeAnswer : PD.Decoder Answer
decodeAnswer =
    let
        conv v =
            case v of
                0 ->
                    XMaybe

                1 ->
                    Yes

                2 ->
                    No

                10 ->
                    PleaseRepeat

                _ ->
                    XMaybe
    in
    PD.map conv PD.int32


encodeAnswer : Answer -> PE.Encoder
encodeAnswer v =
    let
        conv =
            case v of
                XMaybe ->
                    0

                Yes ->
                    1

                No ->
                    2

                PleaseRepeat ->
                    10
    in
    PE.int32 conv
