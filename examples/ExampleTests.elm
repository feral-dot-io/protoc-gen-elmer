module ExampleTests exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Example
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


allTogether_AnswerFuzzer : Fuzzer Example.AllTogether_Answer
allTogether_AnswerFuzzer =
    Fuzz.oneOf
        [ Fuzz.map Example.AllTogether_Maybe fuzzInt32
        , Fuzz.constant Example.AllTogether_Yes
        , Fuzz.constant Example.AllTogether_No
        ]


allTogetherFuzzer : Fuzzer Example.AllTogether
allTogetherFuzzer =
    let
        allTogether_FavouriteFuzzer =
            Fuzz.oneOf
                [ Fuzz.map Example.AllTogether_MyStr Fuzz.string
                , Fuzz.map Example.AllTogether_MyNum fuzzInt32
                , Fuzz.map Example.AllTogether_Selection scalarFuzzer
                ]

        allTogether_MyNameFuzzer =
            Fuzz.oneOf
                [ Fuzz.string
                ]
    in
    Fuzz.map Example.AllTogether
        (Fuzz.list Fuzz.string)
        |> Fuzz.andMap
            (Fuzz.map Dict.fromList
                (Fuzz.list (Fuzz.tuple ( Fuzz.string, Fuzz.bool )))
            )
        |> Fuzz.andMap (Fuzz.maybe allTogether_FavouriteFuzzer)
        |> Fuzz.andMap (Fuzz.maybe allTogether_MyNameFuzzer)
        |> Fuzz.andMap allTogether_NestedAbcFuzzer
        |> Fuzz.andMap allTogether_AnswerFuzzer


allTogether_NestedAbcFuzzer : Fuzzer Example.AllTogether_NestedAbc
allTogether_NestedAbcFuzzer =
    Fuzz.map Example.AllTogether_NestedAbc
        fuzzInt32
        |> Fuzz.andMap fuzzInt32
        |> Fuzz.andMap fuzzInt32


scalarFuzzer : Fuzzer Example.Scalar
scalarFuzzer =
    Fuzz.map Example.Scalar
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
            PE.encode (Example.allTogetherEncoder data)
                |> PD.decode Example.allTogetherDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether"
        [ test "empty" (\_ -> run Example.emptyAllTogether)
        , fuzz allTogetherFuzzer "fuzzer" run
        ]


testAllTogether_NestedAbc : Test
testAllTogether_NestedAbc =
    let
        run data =
            PE.encode (Example.allTogether_NestedAbcEncoder data)
                |> PD.decode Example.allTogether_NestedAbcDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether_NestedAbc"
        [ test "empty" (\_ -> run Example.emptyAllTogether_NestedAbc)
        , fuzz allTogether_NestedAbcFuzzer "fuzzer" run
        ]


testScalar : Test
testScalar =
    let
        run data =
            PE.encode (Example.scalarEncoder data)
                |> PD.decode Example.scalarDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Scalar"
        [ test "empty" (\_ -> run Example.emptyScalar)
        , fuzz scalarFuzzer "fuzzer" run
        ]


testAllTogether_Answer : Test
testAllTogether_Answer =
    let
        run data =
            PE.encode (Example.allTogether_AnswerEncoder data)
                |> PD.decode Example.allTogether_AnswerDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode AllTogether_Answer"
        [ test "empty" (\_ -> run Example.emptyAllTogether_Answer)
        , fuzz allTogether_AnswerFuzzer "fuzzer" run
        ]
