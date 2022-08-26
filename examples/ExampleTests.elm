module ExampleTests exposing (..)

{-| Protobuf library for testing structures found in package `example`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Example
import Expect
import Fuzz exposing (Fuzzer)
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.ElmerTests
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


fuzzAllTogether_Answer : Fuzzer Example.AllTogether_Answer
fuzzAllTogether_Answer =
    Fuzz.oneOf
        [ Fuzz.constant Example.AllTogether_Maybe
        , Fuzz.constant Example.AllTogether_Yes
        , Fuzz.constant Example.AllTogether_No
        ]


fuzzAllTogether : Fuzzer Example.AllTogether
fuzzAllTogether =
    let
        fuzzAllTogether_Favourite =
            Fuzz.oneOf
                [ Fuzz.map Example.AllTogether_MyStr Fuzz.string
                , Fuzz.map Example.AllTogether_MyNum Protobuf.ElmerTests.fuzzInt32
                , Fuzz.map Example.AllTogether_Selection fuzzScalar
                ]

        fuzzAllTogether_MyName =
            Fuzz.oneOf
                [ Fuzz.string
                ]
    in
    Fuzz.map Example.AllTogether
        (Fuzz.list Fuzz.string)
        |> Fuzz.andMap (Fuzz.map Dict.fromList (Fuzz.list (Fuzz.tuple ( Fuzz.string, Fuzz.bool ))))
        |> Fuzz.andMap (Fuzz.maybe fuzzAllTogether_Favourite)
        |> Fuzz.andMap (Fuzz.maybe fuzzAllTogether_MyName)
        |> Fuzz.andMap fuzzAllTogether_NestedAbc
        |> Fuzz.andMap fuzzAllTogether_Answer


fuzzAllTogether_NestedAbc : Fuzzer Example.AllTogether_NestedAbc
fuzzAllTogether_NestedAbc =
    Fuzz.map Example.AllTogether_NestedAbc
        Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzInt32


fuzzScalar : Fuzzer Example.Scalar
fuzzScalar =
    Fuzz.map Example.Scalar
        Fuzz.float
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzFloat32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzBytes


testAllTogether : Test
testAllTogether =
    let
        run =
            Protobuf.ElmerTests.runTest Example.decodeAllTogether Example.encodeAllTogether
    in
    Test.describe "encode then decode AllTogether"
        [ test "empty" (\_ -> run Example.emptyAllTogether)
        , fuzz fuzzAllTogether "fuzzer" run
        ]


testAllTogether_NestedAbc : Test
testAllTogether_NestedAbc =
    let
        run =
            Protobuf.ElmerTests.runTest Example.decodeAllTogether_NestedAbc Example.encodeAllTogether_NestedAbc
    in
    Test.describe "encode then decode AllTogether_NestedAbc"
        [ test "empty" (\_ -> run Example.emptyAllTogether_NestedAbc)
        , fuzz fuzzAllTogether_NestedAbc "fuzzer" run
        ]


testScalar : Test
testScalar =
    let
        run =
            Protobuf.ElmerTests.runTest Example.decodeScalar Example.encodeScalar
    in
    Test.describe "encode then decode Scalar"
        [ test "empty" (\_ -> run Example.emptyScalar)
        , fuzz fuzzScalar "fuzzer" run
        ]


testAllTogether_Answer : Test
testAllTogether_Answer =
    let
        run =
            Protobuf.ElmerTests.runTest Example.decodeAllTogether_Answer Example.encodeAllTogether_Answer
    in
    Test.describe "encode then decode AllTogether_Answer"
        [ test "empty" (\_ -> run Example.emptyAllTogether_Answer)
        , fuzz fuzzAllTogether_Answer "fuzzer" run
        ]
