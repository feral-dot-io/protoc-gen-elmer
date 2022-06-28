-- Generated by protoc-gen-elmgen. DO NOT EDIT!


module Gen.HaberdasherTests exposing (..)

import Expect
import Fuzz exposing (Fuzzer)
import Gen.Haberdasher as Codec
import Protobuf.Decode as PD
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


fuzzInt32 : Fuzzer Int
fuzzInt32 =
    Fuzz.intRange -2147483648 2147483647


hatFuzzer : Fuzzer Codec.Hat
hatFuzzer =
    Fuzz.map Codec.Hat
        fuzzInt32
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string


sizeFuzzer : Fuzzer Codec.Size
sizeFuzzer =
    Fuzz.map Codec.Size
        fuzzInt32


testHat : Test
testHat =
    let
        run data =
            PE.encode (Codec.hatEncoder data)
                |> PD.decode Codec.hatDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Hat"
        [ test "empty" (\_ -> run Codec.emptyHat)
        , fuzz hatFuzzer "fuzzer" run
        ]


testSize : Test
testSize =
    let
        run data =
            PE.encode (Codec.sizeEncoder data)
                |> PD.decode Codec.sizeDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Size"
        [ test "empty" (\_ -> run Codec.emptySize)
        , fuzz sizeFuzzer "fuzzer" run
        ]