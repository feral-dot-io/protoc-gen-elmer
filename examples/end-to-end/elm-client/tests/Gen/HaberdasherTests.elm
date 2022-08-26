module Gen.HaberdasherTests exposing (..)

{-| Protobuf library for testing structures found in package `gen.haberdasher`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Expect
import Fuzz exposing (Fuzzer)
import Gen.Haberdasher
import Protobuf.Decode as PD
import Protobuf.ElmerTests
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


fuzzHat : Fuzzer Gen.Haberdasher.Hat
fuzzHat =
    Fuzz.map Gen.Haberdasher.Hat
        Protobuf.ElmerTests.fuzzInt32
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string


fuzzSize : Fuzzer Gen.Haberdasher.Size
fuzzSize =
    Fuzz.map Gen.Haberdasher.Size
        Protobuf.ElmerTests.fuzzInt32


testHat : Test
testHat =
    let
        run =
            Protobuf.ElmerTests.runTest Gen.Haberdasher.decodeHat Gen.Haberdasher.encodeHat
    in
    Test.describe "encode then decode Hat"
        [ test "empty" (\_ -> run Gen.Haberdasher.emptyHat)
        , fuzz fuzzHat "fuzzer" run
        ]


testSize : Test
testSize =
    let
        run =
            Protobuf.ElmerTests.runTest Gen.Haberdasher.decodeSize Gen.Haberdasher.encodeSize
    in
    Test.describe "encode then decode Size"
        [ test "empty" (\_ -> run Gen.Haberdasher.emptySize)
        , fuzz fuzzSize "fuzzer" run
        ]
