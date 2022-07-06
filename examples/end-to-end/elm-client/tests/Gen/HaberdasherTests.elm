module Gen.HaberdasherTests exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Expect
import Fuzz exposing (Fuzzer)
import Gen.Haberdasher
import Protobuf.Decode as PD
import Protobuf.ElmerTest
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)


hatFuzzer : Fuzzer Gen.Haberdasher.Hat
hatFuzzer =
    Fuzz.map Gen.Haberdasher.Hat
        Protobuf.ElmerTest.fuzzInt32
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string


sizeFuzzer : Fuzzer Gen.Haberdasher.Size
sizeFuzzer =
    Fuzz.map Gen.Haberdasher.Size
        Protobuf.ElmerTest.fuzzInt32


testHat : Test
testHat =
    let
        run data =
            PE.encode (Gen.Haberdasher.hatEncoder data)
                |> PD.decode Gen.Haberdasher.hatDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Hat"
        [ test "empty" (\_ -> run Gen.Haberdasher.emptyHat)
        , fuzz hatFuzzer "fuzzer" run
        ]


testSize : Test
testSize =
    let
        run data =
            PE.encode (Gen.Haberdasher.sizeEncoder data)
                |> PD.decode Gen.Haberdasher.sizeDecoder
                |> Expect.equal (Just data)
    in
    Test.describe "encode then decode Size"
        [ test "empty" (\_ -> run Gen.Haberdasher.emptySize)
        , fuzz sizeFuzzer "fuzzer" run
        ]
