module Gen.Haberdasher exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Protobuf.Decode as PD
import Protobuf.ElmerTest
import Protobuf.Encode as PE


{-| A Hat is a piece of headwear made by a Haberdasher.
-}
type alias Hat =
    --  The size of a hat should always be in inches.
    { size : Int

    --  The color of a hat will never be 'invisible', but other than
    --  that, anything is fair game.
    , color : String

    --  The name of a hat is it's type. Like, 'bowler', or something.
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


hatDecoder : PD.Decoder Hat
hatDecoder =
    PD.message emptyHat
        [ PD.optional 1 PD.int32 (\v m -> { m | size = v })
        , PD.optional 2 PD.string (\v m -> { m | color = v })
        , PD.optional 3 PD.string (\v m -> { m | name = v })
        ]


sizeDecoder : PD.Decoder Size
sizeDecoder =
    PD.message emptySize
        [ PD.optional 1 PD.int32 (\v m -> { m | inches = v })
        ]


hatEncoder : Hat -> PE.Encoder
hatEncoder v =
    PE.message <|
        [ ( 1, PE.int32 v.size )
        , ( 2, PE.string v.color )
        , ( 3, PE.string v.name )
        ]


sizeEncoder : Size -> PE.Encoder
sizeEncoder v =
    PE.message <|
        [ ( 1, PE.int32 v.inches )
        ]
