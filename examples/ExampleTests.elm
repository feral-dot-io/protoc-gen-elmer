module ExampleTests exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Example
import Expect
import Fuzz exposing (Fuzzer)
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.ElmerTest
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


allTogether_AnswerFuzzer : Fuzzer Example.AllTogether_Answer
allTogether_AnswerFuzzer =
    Fuzz.oneOf
        [ Fuzz.map Example.AllTogether_Maybe Protobuf.ElmerTest.fuzzInt32
        , Fuzz.constant Example.AllTogether_Yes
        , Fuzz.constant Example.AllTogether_No
        ]


allTogetherFuzzer : Fuzzer Example.AllTogether
allTogetherFuzzer =
    let
        allTogether_FavouriteFuzzer =
            Fuzz.oneOf
                [ Fuzz.map Example.AllTogether_MyStr Fuzz.string
                , Fuzz.map Example.AllTogether_MyNum Protobuf.ElmerTest.fuzzInt32
                , Fuzz.map Example.AllTogether_Selection scalarFuzzer
                ]

        allTogether_MyNameFuzzer =
            Fuzz.oneOf
                [ Fuzz.string
                ]
    in
    Fuzz.map Example.AllTogether
        (Fuzz.list Fuzz.string)
        |> Fuzz.andMap (Fuzz.map Dict.fromList (Fuzz.list (Fuzz.tuple ( Fuzz.string, Fuzz.bool ))))
        |> Fuzz.andMap (Fuzz.maybe allTogether_FavouriteFuzzer)
        |> Fuzz.andMap (Fuzz.maybe allTogether_MyNameFuzzer)
        |> Fuzz.andMap allTogether_NestedAbcFuzzer
        |> Fuzz.andMap allTogether_AnswerFuzzer


allTogether_NestedAbcFuzzer : Fuzzer Example.AllTogether_NestedAbc
allTogether_NestedAbcFuzzer =
    Fuzz.map Example.AllTogether_NestedAbc
        Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzInt32


scalarFuzzer : Fuzzer Example.Scalar
scalarFuzzer =
    Fuzz.map Example.Scalar
        Fuzz.float
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzFloat32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Protobuf.ElmerTest.fuzzBytes


testAllTogether : Test
testAllTogether =
    let
        run =
            Protobuf.ElmerTest.runTest Example.allTogetherDecoder Example.allTogetherEncoder
    in
    Test.describe "encode then decode AllTogether"
        [ test "empty" (\_ -> run Example.emptyAllTogether)
        , fuzz allTogetherFuzzer "fuzzer" run
        ]


testAllTogether_NestedAbc : Test
testAllTogether_NestedAbc =
    let
        run =
            Protobuf.ElmerTest.runTest Example.allTogether_NestedAbcDecoder Example.allTogether_NestedAbcEncoder
    in
    Test.describe "encode then decode AllTogether_NestedAbc"
        [ test "empty" (\_ -> run Example.emptyAllTogether_NestedAbc)
        , fuzz allTogether_NestedAbcFuzzer "fuzzer" run
        ]


testScalar : Test
testScalar =
    let
        run =
            Protobuf.ElmerTest.runTest Example.scalarDecoder Example.scalarEncoder
    in
    Test.describe "encode then decode Scalar"
        [ test "empty" (\_ -> run Example.emptyScalar)
        , fuzz scalarFuzzer "fuzzer" run
        ]


testAllTogether_Answer : Test
testAllTogether_Answer =
    let
        run =
            Protobuf.ElmerTest.runTest Example.allTogether_AnswerDecoder Example.allTogether_AnswerEncoder
    in
    Test.describe "encode then decode AllTogether_Answer"
        [ test "empty" (\_ -> run Example.emptyAllTogether_Answer)
        , fuzz allTogether_AnswerFuzzer "fuzzer" run
        ]
