module Ex03Complex exposing (..)

{-| Protobuf library for decoding and encoding structures found in package `Ex03Complex` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - Ex
  - Ex\_Questions
  - Ex\_WellKnownHolder

Unions: (none)

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs Ex, Ex_Questions, Ex_WellKnownHolder


# Empty (zero values)

@docs emptyEx, emptyEx_Questions, emptyEx_WellKnownHolder


# Decoders

@docs decodeEx, decodeEx_Questions, decodeEx_WellKnownHolder


# Encoders

@docs encodeEx, encodeEx_Questions, encodeEx_WellKnownHolder

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Ex02Enums
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE
import Time


type alias Ex =
    -- We can have put things in lists
    { marbles : List Float -- Diameter (in cm) -- post comments
    , lastMarbleLost : Maybe Float

    -- A special kind of enum that holds `Maybe pick_one`
    , pickOne : Maybe Ex_PickOne

    -- Let's not forget -- maps map to a Dict
    , pressingQuestions : Dict String Ex02Enums.Answer
    }


{-| We can nest messages
-}
type alias Ex_Questions =
    -- Reference another .proto's message
    { brunch : Ex02Enums.Answer

    -- Note how message fields aren't wrapped in a Maybe
    , pudding : Ex02Enums.Answer
    }


{-| Note the runs of caps gets converted to `WellKnownHolder`
-}
type alias Ex_WellKnownHolder =
    -- Timestamp becomes `Time.Posix` directly
    { createdOn : Time.Posix

    -- Wrappers allow us to put scalars behind a `Maybe Int`
    , uncertainInteger : Protobuf.Elmer.Int32Value
    }


type
    Ex_PickOne
    -- We use `elm/Bytes` directly
    = Ex_PngOfMarblesLostOverTime Bytes
      -- Maps to `String`
    | Ex_EssayOnNotLosingMyMarbles String
      -- We can also include other messages
    | Ex_Complex Ex_WellKnownHolder


emptyEx : Ex
emptyEx =
    Ex [] Nothing Nothing Dict.empty


emptyEx_Questions : Ex_Questions
emptyEx_Questions =
    Ex_Questions Ex02Enums.emptyAnswer Ex02Enums.emptyAnswer


emptyEx_WellKnownHolder : Ex_WellKnownHolder
emptyEx_WellKnownHolder =
    Ex_WellKnownHolder Protobuf.Elmer.emptyTimestamp Protobuf.Elmer.emptyInt32Value


decodeEx : PD.Decoder Ex
decodeEx =
    let
        decodeEx_LastMarbleLost =
            [ ( 2, PD.double )
            ]

        decodeEx_PickOne =
            [ ( 3, PD.map Ex_PngOfMarblesLostOverTime PD.bytes )
            , ( 4, PD.map Ex_EssayOnNotLosingMyMarbles PD.string )
            , ( 5, PD.map Ex_Complex decodeEx_WellKnownHolder )
            ]
    in
    PD.message emptyEx
        [ PD.repeated 1 PD.double .marbles (\v m -> { m | marbles = v })
        , PD.oneOf decodeEx_LastMarbleLost (\v m -> { m | lastMarbleLost = v })
        , PD.oneOf decodeEx_PickOne (\v m -> { m | pickOne = v })
        , PD.mapped 6 ( "", Ex02Enums.emptyAnswer ) PD.string Ex02Enums.decodeAnswer .pressingQuestions (\v m -> { m | pressingQuestions = v })
        ]


decodeEx_Questions : PD.Decoder Ex_Questions
decodeEx_Questions =
    PD.message emptyEx_Questions
        [ PD.optional 1 Ex02Enums.decodeAnswer (\v m -> { m | brunch = v })
        , PD.optional 2 Ex02Enums.decodeAnswer (\v m -> { m | pudding = v })
        ]


decodeEx_WellKnownHolder : PD.Decoder Ex_WellKnownHolder
decodeEx_WellKnownHolder =
    PD.message emptyEx_WellKnownHolder
        [ PD.optional 1 Protobuf.Elmer.decodeTimestamp (\v m -> { m | createdOn = v })
        , PD.optional 16 Protobuf.Elmer.decodeInt32Value (\v m -> { m | uncertainInteger = v })
        ]


encodeEx : Ex -> PE.Encoder
encodeEx v =
    let
        encodeEx_LastMarbleLost o =
            case o of
                Just data ->
                    [ ( 2, PE.double data ) ]

                Nothing ->
                    []

        encodeEx_PickOne o =
            case o of
                Just (Ex_PngOfMarblesLostOverTime data) ->
                    [ ( 3, PE.bytes data ) ]

                Just (Ex_EssayOnNotLosingMyMarbles data) ->
                    [ ( 4, PE.string data ) ]

                Just (Ex_Complex data) ->
                    [ ( 5, encodeEx_WellKnownHolder data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 1, PE.list PE.double v.marbles )
        , ( 6, PE.dict PE.string Ex02Enums.encodeAnswer v.pressingQuestions )
        ]
            ++ encodeEx_LastMarbleLost v.lastMarbleLost
            ++ encodeEx_PickOne v.pickOne


encodeEx_Questions : Ex_Questions -> PE.Encoder
encodeEx_Questions v =
    PE.message <|
        [ ( 1, Ex02Enums.encodeAnswer v.brunch )
        , ( 2, Ex02Enums.encodeAnswer v.pudding )
        ]


encodeEx_WellKnownHolder : Ex_WellKnownHolder -> PE.Encoder
encodeEx_WellKnownHolder v =
    PE.message <|
        [ ( 1, Protobuf.Elmer.encodeTimestamp v.createdOn )
        , ( 16, Protobuf.Elmer.encodeInt32Value v.uncertainInteger )
        ]
