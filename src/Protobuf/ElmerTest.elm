module Protobuf.ElmerTest exposing
    ( anyFuzzer
    , apiFuzzer
    , boolValueFuzzer
    , bytesValueFuzzer
    , doubleValueFuzzer
    , durationFuzzer
    , emptyFuzzer
    , enumFuzzer
    , enumValueFuzzer
    , fieldFuzzer
    , fieldMaskFuzzer
    , field_CardinalityFuzzer
    , field_KindFuzzer
    , floatValueFuzzer
    , fuzzBytes
    , fuzzFloat32
    , fuzzInt32
    , fuzzUInt32
    , int32ValueFuzzer
    , int64ValueFuzzer
    , listValueFuzzer
    , methodFuzzer
    , mixinFuzzer
    , nullValueFuzzer
    , optionFuzzer
    , sourceContextFuzzer
    , stringValueFuzzer
    , structFuzzer
    , syntaxFuzzer
    , timestampFuzzer
    , typeFuzzer
    , uInt32ValueFuzzer
    , uInt64ValueFuzzer
    , valueFuzzer
    )

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Fuzz exposing (Fuzzer)
import Google.Protobuf as GP
import Protobuf.Elmer as Elmer
import Time



-- Protobuf-specific fuzzers


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


posInt32Fuzzer : Fuzzer Int
posInt32Fuzzer =
    minInt32Fuzzer 0


minInt32Fuzzer : Int -> Fuzzer Int
minInt32Fuzzer min =
    Fuzz.intRange min 2147483647



-- Fuzzers for Google.Protobuf pass through. Avoids deepd nesting


anyFuzzer : Fuzzer GP.Any
anyFuzzer =
    Fuzz.map2 GP.Any Fuzz.string fuzzBytes


apiFuzzer : Fuzzer GP.Api
apiFuzzer =
    Fuzz.map GP.Api Fuzz.string
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap (Fuzz.maybe sourceContextFuzzer)
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap syntaxFuzzer


durationFuzzer : Fuzzer GP.Duration
durationFuzzer =
    Fuzz.map2 GP.Duration posInt32Fuzzer (Fuzz.intRange 0 1000)


emptyFuzzer : Fuzzer GP.Empty
emptyFuzzer =
    Fuzz.constant GP.Empty


enumFuzzer : Fuzzer GP.Enum
enumFuzzer =
    Fuzz.map5 GP.Enum
        Fuzz.string
        (Fuzz.constant [])
        (Fuzz.constant [])
        (Fuzz.maybe sourceContextFuzzer)
        syntaxFuzzer


enumValueFuzzer : Fuzzer GP.EnumValue
enumValueFuzzer =
    Fuzz.map3 GP.EnumValue Fuzz.string posInt32Fuzzer (Fuzz.list optionFuzzer)


fieldFuzzer : Fuzzer GP.Field
fieldFuzzer =
    Fuzz.map GP.Field field_KindFuzzer
        |> Fuzz.andMap field_CardinalityFuzzer
        |> Fuzz.andMap posInt32Fuzzer
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap posInt32Fuzzer
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string


field_CardinalityFuzzer : Fuzzer GP.Cardinality
field_CardinalityFuzzer =
    Fuzz.oneOf
        [ Fuzz.constant GP.CardinalityUnknown
        , Fuzz.constant GP.CardinalityOptional
        , Fuzz.constant GP.CardinalityRequired
        , Fuzz.constant GP.CardinalityRepeated
        , Fuzz.map GP.CardinalityUnrecognized_ (minInt32Fuzzer 4)
        ]


field_KindFuzzer : Fuzzer GP.Kind
field_KindFuzzer =
    Fuzz.oneOf
        [ Fuzz.constant GP.TypeUnknown
        , Fuzz.constant GP.TypeDouble
        , Fuzz.constant GP.TypeFloat
        , Fuzz.constant GP.TypeInt64
        , Fuzz.constant GP.TypeUint64
        , Fuzz.constant GP.TypeInt32
        , Fuzz.constant GP.TypeFixed64
        , Fuzz.constant GP.TypeFixed32
        , Fuzz.constant GP.TypeBool
        , Fuzz.constant GP.TypeString
        , Fuzz.constant GP.TypeGroup
        , Fuzz.constant GP.TypeMessage
        , Fuzz.constant GP.TypeBytes
        , Fuzz.constant GP.TypeUint32
        , Fuzz.constant GP.TypeEnum
        , Fuzz.constant GP.TypeSfixed32
        , Fuzz.constant GP.TypeSfixed64
        , Fuzz.constant GP.TypeSint32
        , Fuzz.constant GP.TypeSint64
        , Fuzz.map GP.KindUnrecognized_ (minInt32Fuzzer 19)
        ]


fieldMaskFuzzer : Fuzzer GP.FieldMask
fieldMaskFuzzer =
    Fuzz.map GP.FieldMask (Fuzz.constant [])


listValueFuzzer : Fuzzer GP.ListValue
listValueFuzzer =
    Fuzz.map GP.ListValue (Fuzz.map GP.ListValueValues (Fuzz.constant []))


methodFuzzer : Fuzzer GP.Method
methodFuzzer =
    Fuzz.map GP.Method Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap syntaxFuzzer


mixinFuzzer : Fuzzer GP.Mixin
mixinFuzzer =
    Fuzz.map2 GP.Mixin Fuzz.string Fuzz.string


nullValueFuzzer : Fuzzer GP.NullValue
nullValueFuzzer =
    Fuzz.oneOf
        [ Fuzz.constant GP.NullValue
        , Fuzz.map GP.NullValueUnrecognized_ (minInt32Fuzzer 1)
        ]


optionFuzzer : Fuzzer GP.Option
optionFuzzer =
    Fuzz.map2 GP.Option Fuzz.string (Fuzz.maybe anyFuzzer)


sourceContextFuzzer : Fuzzer GP.SourceContext
sourceContextFuzzer =
    Fuzz.map GP.SourceContext Fuzz.string


structFuzzer : Fuzzer GP.Struct
structFuzzer =
    Fuzz.map GP.Struct (Fuzz.map GP.StructFields (Fuzz.constant Dict.empty))


syntaxFuzzer : Fuzzer GP.Syntax
syntaxFuzzer =
    Fuzz.oneOf
        [ Fuzz.constant GP.SyntaxProto2
        , Fuzz.constant GP.SyntaxProto3
        , Fuzz.map GP.SyntaxUnrecognized_ (minInt32Fuzzer 2)
        ]


typeFuzzer : Fuzzer GP.Type
typeFuzzer =
    Fuzz.map GP.Type Fuzz.string
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.maybe sourceContextFuzzer)
        |> Fuzz.andMap syntaxFuzzer


valueFuzzer : Fuzzer GP.Value
valueFuzzer =
    let
        kindTypeFuzzer =
            Fuzz.oneOf
                [ Fuzz.map GP.KindNullValue nullValueFuzzer
                , Fuzz.map GP.KindNumberValue Fuzz.float
                , Fuzz.map GP.KindStringValue Fuzz.string
                , Fuzz.map GP.KindBoolValue Fuzz.bool
                , Fuzz.map GP.KindStructValue structFuzzer
                , Fuzz.map GP.KindListValue listValueFuzzer
                ]
    in
    Fuzz.map GP.Value (Fuzz.map GP.ValueKind (Fuzz.maybe kindTypeFuzzer))
