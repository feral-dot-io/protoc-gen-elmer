module Example exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Protobuf.Decode as PD
import Protobuf.Encode as PE


{-| Enums!
-}
type AllTogether_Answer
    = AllTogether_Maybe Int
    | AllTogether_Yes
    | AllTogether_No


{-| A complex record with lots of features
-}
type alias AllTogether =
    { --  Lists
      myList : List String

    --  Maps
    , myMap : Dict String Bool

    --  A nilable sum type
    , favourite : Maybe AllTogether_Favourite
    , my_name : Maybe String
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

    --  Underling eriktim/elm-protocol-buffers library does not support 64-bit
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
    Scalar 0 0 0 0 0 0 0 False "" (BE.encode (BE.sequence []))


emptyAllTogether_Answer : AllTogether_Answer
emptyAllTogether_Answer =
    AllTogether_Maybe 0


allTogetherDecoder : PD.Decoder AllTogether
allTogetherDecoder =
    let
        allTogether_FavouriteDecoder =
            [ ( 3, PD.map AllTogether_MyStr PD.string )
            , ( 4, PD.map AllTogether_MyNum PD.int32 )
            , ( 5, PD.map AllTogether_Selection scalarDecoder )
            ]

        allTogether_MyNameDecoder =
            [ ( 6, PD.string )
            ]
    in
    PD.message emptyAllTogether
        [ PD.repeated 1 PD.string .myList (\v m -> { m | myList = v })
        , PD.mapped 2 ( "", False ) PD.string PD.bool .myMap (\v m -> { m | myMap = v })
        , PD.oneOf allTogether_FavouriteDecoder (\v m -> { m | favourite = v })
        , PD.oneOf allTogether_MyNameDecoder (\v m -> { m | my_name = v })
        , PD.optional 7 allTogether_NestedAbcDecoder (\v m -> { m | abc = v })
        , PD.optional 8 allTogether_AnswerDecoder (\v m -> { m | answer = v })
        ]


allTogether_NestedAbcDecoder : PD.Decoder AllTogether_NestedAbc
allTogether_NestedAbcDecoder =
    PD.message emptyAllTogether_NestedAbc
        [ PD.optional 1 PD.int32 (\v m -> { m | a = v })
        , PD.optional 2 PD.int32 (\v m -> { m | b = v })
        , PD.optional 3 PD.int32 (\v m -> { m | c = v })
        ]


scalarDecoder : PD.Decoder Scalar
scalarDecoder =
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


allTogether_AnswerDecoder : PD.Decoder AllTogether_Answer
allTogether_AnswerDecoder =
    let
        conv v =
            case v of
                1 ->
                    AllTogether_Yes

                2 ->
                    AllTogether_No

                wire ->
                    AllTogether_Maybe wire
    in
    PD.map conv PD.int32


allTogetherEncoder : AllTogether -> PE.Encoder
allTogetherEncoder v =
    let
        allTogether_FavouriteEncoder o =
            case o of
                Just (AllTogether_MyStr data) ->
                    [ ( 3, PE.string data ) ]

                Just (AllTogether_MyNum data) ->
                    [ ( 4, PE.int32 data ) ]

                Just (AllTogether_Selection data) ->
                    [ ( 5, scalarEncoder data ) ]

                Nothing ->
                    []

        allTogether_MyNameEncoder o =
            case o of
                Just data ->
                    [ ( 6, PE.string data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 1, PE.list PE.string v.myList )
        , ( 2, PE.dict PE.string PE.bool v.myMap )
        , ( 7, allTogether_NestedAbcEncoder v.abc )
        , ( 8, allTogether_AnswerEncoder v.answer )
        ]
            ++ allTogether_FavouriteEncoder v.favourite
            ++ allTogether_MyNameEncoder v.my_name


allTogether_NestedAbcEncoder : AllTogether_NestedAbc -> PE.Encoder
allTogether_NestedAbcEncoder v =
    PE.message <|
        [ ( 1, PE.int32 v.a )
        , ( 2, PE.int32 v.b )
        , ( 3, PE.int32 v.c )
        ]


scalarEncoder : Scalar -> PE.Encoder
scalarEncoder v =
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


allTogether_AnswerEncoder : AllTogether_Answer -> PE.Encoder
allTogether_AnswerEncoder v =
    let
        conv =
            case v of
                AllTogether_Maybe wire ->
                    wire

                AllTogether_Yes ->
                    1

                AllTogether_No ->
                    2
    in
    PE.int32 conv
