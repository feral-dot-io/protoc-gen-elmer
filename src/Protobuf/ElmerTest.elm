module Protobuf.ElmerTest exposing
    ( boolValueFuzzer
    , bytesValueFuzzer
    , doubleValueFuzzer
    , floatValueFuzzer
    , fuzzBytes
    , fuzzFloat32
    , fuzzInt32
    , fuzzUInt32
    , int32ValueFuzzer
    , int64ValueFuzzer
    , stringValueFuzzer
    , timestampFuzzer
    , uInt32ValueFuzzer
    , uInt64ValueFuzzer
    )

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Fuzz exposing (Fuzzer)
import Protobuf.Elmer as Elmer
import Time



-- Helpers


fuzzInt32 : Fuzzer Int
fuzzInt32 =
    Fuzz.intRange -2147483648 2147483647


fuzzUInt32 : Fuzzer Int
fuzzUInt32 =
    Fuzz.intRange 0 4294967295


{-| Tests float32' exponent (8 bits).
Avoids trying to robusly map float64 (JS) -> float32
-}
fuzzFloat32 : Fuzzer Float
fuzzFloat32 =
    Fuzz.map (\i -> 2 ^ toFloat i) fuzzInt32


fuzzBytes : Fuzzer Bytes
fuzzBytes =
    Fuzz.intRange 0 255
        |> Fuzz.map BE.unsignedInt8
        |> Fuzz.list
        |> Fuzz.map (BE.sequence >> BE.encode)



-- Fuzzers for well-known types


boolValueFuzzer : Fuzzer Elmer.BoolValue
boolValueFuzzer =
    Fuzz.maybe Fuzz.bool


bytesValueFuzzer : Fuzzer Elmer.BytesValue
bytesValueFuzzer =
    Fuzz.maybe fuzzBytes


doubleValueFuzzer : Fuzzer Elmer.FloatValue
doubleValueFuzzer =
    Fuzz.maybe Fuzz.float


floatValueFuzzer : Fuzzer Elmer.FloatValue
floatValueFuzzer =
    Fuzz.maybe fuzzFloat32


int32ValueFuzzer : Fuzzer Elmer.Int32Value
int32ValueFuzzer =
    Fuzz.maybe fuzzInt32


int64ValueFuzzer : Fuzzer Elmer.Int64Value
int64ValueFuzzer =
    Fuzz.maybe fuzzInt32


stringValueFuzzer : Fuzzer Elmer.StringValue
stringValueFuzzer =
    Fuzz.maybe Fuzz.string


timestampFuzzer : Fuzzer Elmer.Timestamp
timestampFuzzer =
    fuzzUInt32 |> Fuzz.map Time.millisToPosix


uInt32ValueFuzzer : Fuzzer Elmer.UInt32Value
uInt32ValueFuzzer =
    Fuzz.maybe fuzzUInt32


uInt64ValueFuzzer : Fuzzer Elmer.UInt64Value
uInt64ValueFuzzer =
    Fuzz.maybe fuzzUInt32
