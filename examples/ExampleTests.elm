-- Generated by protoc-gen-elmgen. DO NOT EDIT!


module ExampleTests exposing (..)

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Example as Codec
import Expect
import Fuzz exposing (Fuzzer)
import Protobuf.Decode as PD
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


fuzzInt32 : Fuzzer Int
fuzzInt32 =
    Fuzz.intRange -2147483648 2147483647


fuzzUint32 : Fuzzer Int
fuzzUint32 =
    Fuzz.intRange 0 4294967295


fuzzFloat32 : Fuzzer Float
fuzzFloat32 =
    Fuzz.map (\i -> 2 ^ toFloat i) fuzzInt32


fuzzBytes : Fuzzer Bytes
fuzzBytes =
    Fuzz.intRange 0 255
        |> Fuzz.map BE.unsignedInt8
        |> Fuzz.list
        |> Fuzz.map (BE.sequence >> BE.encode)


allTogether_AnswerFuzzer : Fuzzer Codec.AllTogether_Answer
allTogether_AnswerFuzzer =
    Fuzz.oneOf
        [ Fuzz.map Codec.Maybe_AllTogether_Answer fuzzInt32
        , Fuzz.constant Codec.Yes_AllTogether_Answer
        , Fuzz.constant Codec.No_AllTogether_Answer
        ]


allTogetherFuzzer : Fuzzer Codec.AllTogether
allTogetherFuzzer =
    let
        allTogether_FavouriteFuzzer =
            Fuzz.oneOf
                [ Fuzz.map Codec.MyStr_AllTogether_Favourite Fuzz.string
                , Fuzz.map Codec.MyNum_AllTogether_Favourite fuzzInt32
                , Fuzz.map Codec.Selection_AllTogether_Favourite scalarFuzzer
                ]

        example_AllTogether_MyNameFuzzer =
            Fuzz.oneOf
                [ Fuzz.string
                ]
    in
    Fuzz.map Codec.AllTogether
        (Fuzz.list Fuzz.string)
        |> Fuzz.andMap
            (Fuzz.map Dict.fromList
                (Fuzz.list (Fuzz.tuple ( Fuzz.string, Fuzz.bool )))
            )
        |> Fuzz.andMap (Fuzz.maybe allTogether_FavouriteFuzzer)
        |> Fuzz.andMap (Fuzz.maybe example_AllTogether_MyNameFuzzer)
        |> Fuzz.andMap allTogether_NestedAbcFuzzer
        |> Fuzz.andMap allTogether_AnswerFuzzer


allTogether_NestedAbcFuzzer : Fuzzer Codec.AllTogether_NestedAbc
allTogether_NestedAbcFuzzer =
    Fuzz.map Codec.AllTogether_NestedAbc
        fuzzInt32
        |> Fuzz.andMap fuzzInt32
        |> Fuzz.andMap fuzzInt32


scalarFuzzer : Fuzzer Codec.Scalar
scalarFuzzer =
    Fuzz.map Codec.Scalar
        Fuzz.float
        |> Fuzz.andMap fuzzFloat32
        |> Fuzz.andMap fuzzInt32
        |> Fuzz.andMap fuzzUint32
        |> Fuzz.andMap fuzzInt32
        |> Fuzz.andMap fuzzUint32
        |> Fuzz.andMap fuzzInt32
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap fuzzBytes


testAllTogether : Test
testAllTogether =
    let
        run data =
            PE.encode (Codec.allTogetherEncoder data)
                |> PD.decode Codec.allTogetherDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether"
        [ test "empty" (\_ -> run Codec.emptyAllTogether)
        , fuzz allTogetherFuzzer "fuzzer" run
        ]


testAllTogether_NestedAbc : Test
testAllTogether_NestedAbc =
    let
        run data =
            PE.encode (Codec.allTogether_NestedAbcEncoder data)
                |> PD.decode Codec.allTogether_NestedAbcDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether_NestedAbc"
        [ test "empty" (\_ -> run Codec.emptyAllTogether_NestedAbc)
        , fuzz allTogether_NestedAbcFuzzer "fuzzer" run
        ]


testScalar : Test
testScalar =
    let
        run data =
            PE.encode (Codec.scalarEncoder data)
                |> PD.decode Codec.scalarDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Scalar"
        [ test "empty" (\_ -> run Codec.emptyScalar)
        , fuzz scalarFuzzer "fuzzer" run
        ]


testAllTogether_Answer : Test
testAllTogether_Answer =
    let
        run data =
            PE.encode (Codec.allTogether_AnswerEncoder data)
                |> PD.decode Codec.allTogether_AnswerDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether_Answer"
        [ test "empty" (\_ -> run Codec.emptyAllTogether_Answer)
        , fuzz allTogether_AnswerFuzzer "fuzzer" run
        ]
