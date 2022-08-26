module Gen.Haberdasher exposing (..)

{-| Protobuf library for decoding and encoding structures found in package `gen.haberdasher` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - Hat
  - Size

Unions: (none)

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs Hat, Size


# Empty (zero values)

@docs emptyHat, emptySize


# Decoders

@docs decodeHat, decodeSize


# Encoders

@docs encodeHat, encodeSize

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Protobuf.Decode as PD
import Protobuf.Encode as PE


{-| A Hat is a piece of headwear made by a Haberdasher.
-}
type alias Hat =
    -- The size of a hat should always be in inches.
    { size : Int

    -- The color of a hat will never be 'invisible', but other than
    -- that, anything is fair game.
    , color : String

    -- The name of a hat is it's type. Like, 'bowler', or something.
    , name : String
    }


{-| Size is passed when requesting a new hat to be made. It's always
measured in inches.
-}
type alias Size =
    { inches : Int
    }


emptyHat : Hat
emptyHat =
    Hat 0 "" ""


emptySize : Size
emptySize =
    Size 0


decodeHat : PD.Decoder Hat
decodeHat =
    PD.message emptyHat
        [ PD.optional 1 PD.int32 (\v m -> { m | size = v })
        , PD.optional 2 PD.string (\v m -> { m | color = v })
        , PD.optional 3 PD.string (\v m -> { m | name = v })
        ]


decodeSize : PD.Decoder Size
decodeSize =
    PD.message emptySize
        [ PD.optional 1 PD.int32 (\v m -> { m | inches = v })
        ]


encodeHat : Hat -> PE.Encoder
encodeHat v =
    PE.message <|
        [ ( 1, PE.int32 v.size )
        , ( 2, PE.string v.color )
        , ( 3, PE.string v.name )
        ]


encodeSize : Size -> PE.Encoder
encodeSize v =
    PE.message <|
        [ ( 1, PE.int32 v.inches )
        ]
