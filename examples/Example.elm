module Example exposing (..)

{-| Protobuf library for decoding and encoding structures found in package `example` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - AllTogether
  - AllTogether\_NestedAbc
  - Scalar

Unions:

  - AllTogether\_Answer

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs AllTogether, AllTogether_NestedAbc, Scalar, AllTogether_Answer


# Empty (zero values)

@docs emptyAllTogether, emptyAllTogether_NestedAbc, emptyScalar, emptyAllTogether_Answer


# Enum valuesOf

@docs valuesOfAllTogether_Answer


# Enum and String converters

@docs fromAllTogether_Answer, toAllTogether_Answer


# Decoders

@docs decodeAllTogether, decodeAllTogether_NestedAbc, decodeScalar, decodeAllTogether_Answer


# Encoders

@docs encodeAllTogether, encodeAllTogether_NestedAbc, encodeScalar, encodeAllTogether_Answer

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE


{-| Enums!
-}
type AllTogether_Answer
    = AllTogether_Maybe
    | AllTogether_Yes
    | AllTogether_No


{-| A complex record with lots of features
-}
type alias AllTogether =
    -- Lists
    { myList : List String

    -- Maps
    , myMap : Dict String Bool

    -- A nilable sum type
    , favourite : Maybe AllTogether_Favourite
    , myName : Maybe String
    , abc : AllTogether_NestedAbc
    , answer : AllTogether_Answer
    }


{-| We can nest messages
-}
type alias AllTogether_NestedAbc =
    { a : Int
    , b : Int
    , c : Int
    }


{-| Base types
-}
type alias Scalar =
    { myDouble : Float
    , myFloat : Float
    , myInt32 : Int
    , myUint32 : Int
    , mySint32 : Int
    , myFixed32 : Int
    , mySfixed32 : Int

    -- Underling eriktim/elm-protocol-buffers library does not support 64-bit
    -- int64 my_int64 = 4;
    -- uint64 my_uint64 = 6;
    -- sint64 my_sint64 = 8;
    -- fixed64 my_fixed64 = 10;
    -- sfixed64 my_sfixed64 = 12;
    , myBool : Bool
    , myString : String
    , myBytes : Bytes
    }


type AllTogether_Favourite
    = AllTogether_MyStr String
    | AllTogether_MyNum Int
    | AllTogether_Selection Scalar


emptyAllTogether : AllTogether
emptyAllTogether =
    AllTogether [] Dict.empty Nothing Nothing emptyAllTogether_NestedAbc emptyAllTogether_Answer


emptyAllTogether_NestedAbc : AllTogether_NestedAbc
emptyAllTogether_NestedAbc =
    AllTogether_NestedAbc 0 0 0


emptyScalar : Scalar
emptyScalar =
    Scalar 0 0 0 0 0 0 0 False "" Protobuf.Elmer.emptyBytes


emptyAllTogether_Answer : AllTogether_Answer
emptyAllTogether_Answer =
    AllTogether_Maybe


valuesOfAllTogether_Answer : List AllTogether_Answer
valuesOfAllTogether_Answer =
    [ AllTogether_Maybe, AllTogether_Yes, AllTogether_No ]


fromAllTogether_Answer : AllTogether_Answer -> String
fromAllTogether_Answer u =
    case u of
        AllTogether_Maybe ->
            "MAYBE"

        AllTogether_Yes ->
            "YES"

        AllTogether_No ->
            "NO"


toAllTogether_Answer : String -> AllTogether_Answer
toAllTogether_Answer str =
    case str of
        "MAYBE" ->
            AllTogether_Maybe

        "YES" ->
            AllTogether_Yes

        "NO" ->
            AllTogether_No

        _ ->
            AllTogether_Maybe


decodeAllTogether : PD.Decoder AllTogether
decodeAllTogether =
    let
        decodeAllTogether_Favourite =
            [ ( 3, PD.map AllTogether_MyStr PD.string )
            , ( 4, PD.map AllTogether_MyNum PD.int32 )
            , ( 5, PD.map AllTogether_Selection decodeScalar )
            ]

        decodeAllTogether_MyName =
            [ ( 6, PD.string )
            ]
    in
    PD.message emptyAllTogether
        [ PD.repeated 1 PD.string .myList (\v m -> { m | myList = v })
        , PD.mapped 2 ( "", False ) PD.string PD.bool .myMap (\v m -> { m | myMap = v })
        , PD.oneOf decodeAllTogether_Favourite (\v m -> { m | favourite = v })
        , PD.oneOf decodeAllTogether_MyName (\v m -> { m | myName = v })
        , PD.optional 7 decodeAllTogether_NestedAbc (\v m -> { m | abc = v })
        , PD.optional 8 decodeAllTogether_Answer (\v m -> { m | answer = v })
        ]


decodeAllTogether_NestedAbc : PD.Decoder AllTogether_NestedAbc
decodeAllTogether_NestedAbc =
    PD.message emptyAllTogether_NestedAbc
        [ PD.optional 1 PD.int32 (\v m -> { m | a = v })
        , PD.optional 2 PD.int32 (\v m -> { m | b = v })
        , PD.optional 3 PD.int32 (\v m -> { m | c = v })
        ]


decodeScalar : PD.Decoder Scalar
decodeScalar =
    PD.message emptyScalar
        [ PD.optional 1 PD.double (\v m -> { m | myDouble = v })
        , PD.optional 2 PD.float (\v m -> { m | myFloat = v })
        , PD.optional 3 PD.int32 (\v m -> { m | myInt32 = v })
        , PD.optional 5 PD.uint32 (\v m -> { m | myUint32 = v })
        , PD.optional 7 PD.sint32 (\v m -> { m | mySint32 = v })
        , PD.optional 9 PD.fixed32 (\v m -> { m | myFixed32 = v })
        , PD.optional 11 PD.sfixed32 (\v m -> { m | mySfixed32 = v })
        , PD.optional 13 PD.bool (\v m -> { m | myBool = v })
        , PD.optional 14 PD.string (\v m -> { m | myString = v })
        , PD.optional 15 PD.bytes (\v m -> { m | myBytes = v })
        ]


decodeAllTogether_Answer : PD.Decoder AllTogether_Answer
decodeAllTogether_Answer =
    let
        conv v =
            case v of
                0 ->
                    AllTogether_Maybe

                1 ->
                    AllTogether_Yes

                2 ->
                    AllTogether_No

                _ ->
                    AllTogether_Maybe
    in
    PD.map conv PD.int32


encodeAllTogether : AllTogether -> PE.Encoder
encodeAllTogether v =
    let
        encodeAllTogether_Favourite o =
            case o of
                Just (AllTogether_MyStr data) ->
                    [ ( 3, PE.string data ) ]

                Just (AllTogether_MyNum data) ->
                    [ ( 4, PE.int32 data ) ]

                Just (AllTogether_Selection data) ->
                    [ ( 5, encodeScalar data ) ]

                Nothing ->
                    []

        encodeAllTogether_MyName o =
            case o of
                Just data ->
                    [ ( 6, PE.string data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 1, PE.list PE.string v.myList )
        , ( 2, PE.dict PE.string PE.bool v.myMap )
        , ( 7, encodeAllTogether_NestedAbc v.abc )
        , ( 8, encodeAllTogether_Answer v.answer )
        ]
            ++ encodeAllTogether_Favourite v.favourite
            ++ encodeAllTogether_MyName v.myName


encodeAllTogether_NestedAbc : AllTogether_NestedAbc -> PE.Encoder
encodeAllTogether_NestedAbc v =
    PE.message <|
        [ ( 1, PE.int32 v.a )
        , ( 2, PE.int32 v.b )
        , ( 3, PE.int32 v.c )
        ]


encodeScalar : Scalar -> PE.Encoder
encodeScalar v =
    PE.message <|
        [ ( 1, PE.double v.myDouble )
        , ( 2, PE.float v.myFloat )
        , ( 3, PE.int32 v.myInt32 )
        , ( 5, PE.uint32 v.myUint32 )
        , ( 7, PE.sint32 v.mySint32 )
        , ( 9, PE.fixed32 v.myFixed32 )
        , ( 11, PE.sfixed32 v.mySfixed32 )
        , ( 13, PE.bool v.myBool )
        , ( 14, PE.string v.myString )
        , ( 15, PE.bytes v.myBytes )
        ]


encodeAllTogether_Answer : AllTogether_Answer -> PE.Encoder
encodeAllTogether_Answer v =
    let
        conv =
            case v of
                AllTogether_Maybe ->
                    0

                AllTogether_Yes ->
                    1

                AllTogether_No ->
                    2
    in
    PE.int32 conv
