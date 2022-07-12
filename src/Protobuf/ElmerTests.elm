-- This file is part of protoc-gen-elmer.
--
-- Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU Lesser General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
--
-- Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.
--
-- You should have received a copy of the GNU Lesser General Public License along with Protoc-gen-elmer. If not, see <https:--www.gnu.org/licenses/>.


module Protobuf.ElmerTests exposing
    ( runTest
    , fuzzAny, fuzzApi, fuzzBoolValue, fuzzBytes, fuzzBytesValue, fuzzDoubleValue, fuzzDuration, fuzzEmpty, fuzzEnum, fuzzEnumValue, fuzzField, fuzzFieldMask, fuzzField_Cardinality, fuzzField_Kind, fuzzFloat32, fuzzFloatValue, fuzzInt32, fuzzInt32Value, fuzzInt64Value, fuzzListValue, fuzzMethod, fuzzMinInt32, fuzzMixin, fuzzNullValue, fuzzOption, fuzzPosInt32, fuzzSourceContext, fuzzStringValue, fuzzStruct, fuzzSyntax, fuzzTimestamp, fuzzUInt32, fuzzUInt32Value, fuzzUInt64Value, fuzzValue, fuzzXType
    )

{-| Helper types and functions for `protoc-gen-elmer` codegen. This module should not be used directly.

See the project on how this may be used: <https://github.com/feral-dot-io/protoc-gen-elmer>


# Test runners

@docs runTest


# Fuzzers

@docs fuzzAny, fuzzApi, fuzzBoolValue, fuzzBytes, fuzzBytesValue, fuzzDoubleValue, fuzzDuration, fuzzEmpty, fuzzEnum, fuzzEnumValue, fuzzField, fuzzFieldMask, fuzzField_Cardinality, fuzzField_Kind, fuzzFloat32, fuzzFloatValue, fuzzInt32, fuzzInt32Value, fuzzInt64Value, fuzzListValue, fuzzMethod, fuzzMinInt32, fuzzMixin, fuzzNullValue, fuzzOption, fuzzPosInt32, fuzzSourceContext, fuzzStringValue, fuzzStruct, fuzzSyntax, fuzzTimestamp, fuzzUInt32, fuzzUInt32Value, fuzzUInt64Value, fuzzValue, fuzzXType

-}

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Expect
import Fuzz exposing (Fuzzer)
import Google.Protobuf as GP
import Protobuf.Decode as PD
import Protobuf.Elmer as Elmer
import Protobuf.Encode as PE
import Time


{-| Executes a test that runs data through an encoder then decodes it. Expect the result to be equal.
-}
runTest : PD.Decoder data -> (data -> PE.Encoder) -> data -> Expect.Expectation
runTest dec enc data =
    PE.encode (enc data)
        |> PD.decode dec
        |> Expect.equal (Just data)



-- Protobuf-specific fuzzers


{-| -}
fuzzInt32 : Fuzzer Int
fuzzInt32 =
    Fuzz.intRange -2147483648 2147483647


{-| -}
fuzzUInt32 : Fuzzer Int
fuzzUInt32 =
    Fuzz.intRange 0 4294967295


{-| Tests float32' exponent (8 bits).
Avoids trying to robusly map float64 (JS) -> float32
{-|-}
-}
fuzzFloat32 : Fuzzer Float
fuzzFloat32 =
    Fuzz.map (\i -> 2 ^ toFloat i) fuzzInt32


{-| -}
fuzzBytes : Fuzzer Bytes
fuzzBytes =
    Fuzz.intRange 0 255
        |> Fuzz.map BE.unsignedInt8
        |> Fuzz.list
        |> Fuzz.map (BE.sequence >> BE.encode)



-- Fuzzers for well-known types


{-| -}
fuzzBoolValue : Fuzzer Elmer.BoolValue
fuzzBoolValue =
    Fuzz.maybe Fuzz.bool


{-| -}
fuzzBytesValue : Fuzzer Elmer.BytesValue
fuzzBytesValue =
    Fuzz.maybe fuzzBytes


{-| -}
fuzzDoubleValue : Fuzzer Elmer.FloatValue
fuzzDoubleValue =
    Fuzz.maybe Fuzz.float


{-| -}
fuzzFloatValue : Fuzzer Elmer.FloatValue
fuzzFloatValue =
    Fuzz.maybe fuzzFloat32


{-| -}
fuzzInt32Value : Fuzzer Elmer.Int32Value
fuzzInt32Value =
    Fuzz.maybe fuzzInt32


{-| -}
fuzzInt64Value : Fuzzer Elmer.Int64Value
fuzzInt64Value =
    Fuzz.maybe fuzzInt32


{-| -}
fuzzStringValue : Fuzzer Elmer.StringValue
fuzzStringValue =
    Fuzz.maybe Fuzz.string


{-| -}
fuzzTimestamp : Fuzzer Time.Posix
fuzzTimestamp =
    fuzzUInt32 |> Fuzz.map Time.millisToPosix


{-| -}
fuzzUInt32Value : Fuzzer Elmer.UInt32Value
fuzzUInt32Value =
    Fuzz.maybe fuzzUInt32


{-| -}
fuzzUInt64Value : Fuzzer Elmer.UInt64Value
fuzzUInt64Value =
    Fuzz.maybe fuzzUInt32


{-| -}
fuzzPosInt32 : Fuzzer Int
fuzzPosInt32 =
    fuzzMinInt32 0


{-| -}
fuzzMinInt32 : Int -> Fuzzer Int
fuzzMinInt32 min =
    Fuzz.intRange min 2147483647



-- Fuzzers for Google.Protobuf pass through. Avoids deepd nesting


{-| -}
fuzzAny : Fuzzer GP.Any
fuzzAny =
    Fuzz.map2 GP.Any Fuzz.string fuzzBytes


{-| -}
fuzzApi : Fuzzer GP.Api
fuzzApi =
    Fuzz.map GP.Api Fuzz.string
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap (Fuzz.maybe fuzzSourceContext)
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap fuzzSyntax


{-| -}
fuzzDuration : Fuzzer GP.Duration
fuzzDuration =
    Fuzz.map2 GP.Duration fuzzPosInt32 (Fuzz.intRange 0 1000)


{-| -}
fuzzEmpty : Fuzzer GP.Empty
fuzzEmpty =
    Fuzz.constant GP.Empty


{-| -}
fuzzEnum : Fuzzer GP.Enum
fuzzEnum =
    Fuzz.map5 GP.Enum
        Fuzz.string
        (Fuzz.constant [])
        (Fuzz.constant [])
        (Fuzz.maybe fuzzSourceContext)
        fuzzSyntax


{-| -}
fuzzEnumValue : Fuzzer GP.EnumValue
fuzzEnumValue =
    Fuzz.map3 GP.EnumValue Fuzz.string fuzzPosInt32 (Fuzz.list fuzzOption)


{-| -}
fuzzField : Fuzzer GP.Field
fuzzField =
    Fuzz.map GP.Field fuzzField_Kind
        |> Fuzz.andMap fuzzField_Cardinality
        |> Fuzz.andMap fuzzPosInt32
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap fuzzPosInt32
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string


{-| -}
fuzzField_Cardinality : Fuzzer GP.Cardinality
fuzzField_Cardinality =
    Fuzz.oneOf
        [ Fuzz.constant GP.CardinalityUnknown
        , Fuzz.constant GP.CardinalityOptional
        , Fuzz.constant GP.CardinalityRequired
        , Fuzz.constant GP.CardinalityRepeated
        , Fuzz.map GP.CardinalityUnrecognized_ (fuzzMinInt32 4)
        ]


{-| -}
fuzzField_Kind : Fuzzer GP.Kind
fuzzField_Kind =
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
        , Fuzz.map GP.KindUnrecognized_ (fuzzMinInt32 19)
        ]


{-| -}
fuzzFieldMask : Fuzzer GP.FieldMask
fuzzFieldMask =
    Fuzz.map GP.FieldMask (Fuzz.constant [])


{-| -}
fuzzListValue : Fuzzer GP.ListValue
fuzzListValue =
    Fuzz.map GP.ListValue (Fuzz.map GP.ListValueValues (Fuzz.constant []))


{-| -}
fuzzMethod : Fuzzer GP.Method
fuzzMethod =
    Fuzz.map GP.Method Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap fuzzSyntax


{-| -}
fuzzMixin : Fuzzer GP.Mixin
fuzzMixin =
    Fuzz.map2 GP.Mixin Fuzz.string Fuzz.string


{-| -}
fuzzNullValue : Fuzzer GP.NullValue
fuzzNullValue =
    Fuzz.oneOf
        [ Fuzz.constant GP.NullValue
        , Fuzz.map GP.NullValueUnrecognized_ (fuzzMinInt32 1)
        ]


{-| -}
fuzzOption : Fuzzer GP.Option
fuzzOption =
    Fuzz.map2 GP.Option Fuzz.string (Fuzz.maybe fuzzAny)


{-| -}
fuzzSourceContext : Fuzzer GP.SourceContext
fuzzSourceContext =
    Fuzz.map GP.SourceContext Fuzz.string


{-| -}
fuzzStruct : Fuzzer GP.Struct
fuzzStruct =
    Fuzz.map GP.Struct (Fuzz.map GP.StructFields (Fuzz.constant Dict.empty))


{-| -}
fuzzSyntax : Fuzzer GP.Syntax
fuzzSyntax =
    Fuzz.oneOf
        [ Fuzz.constant GP.SyntaxProto2
        , Fuzz.constant GP.SyntaxProto3
        , Fuzz.map GP.SyntaxUnrecognized_ (fuzzMinInt32 2)
        ]


{-| -}
fuzzXType : Fuzzer GP.Type
fuzzXType =
    Fuzz.map GP.Type Fuzz.string
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.constant [])
        |> Fuzz.andMap (Fuzz.maybe fuzzSourceContext)
        |> Fuzz.andMap fuzzSyntax


{-| -}
fuzzValue : Fuzzer GP.Value
fuzzValue =
    let
        kindTypeFuzzer =
            Fuzz.oneOf
                [ Fuzz.map GP.KindNullValue fuzzNullValue
                , Fuzz.map GP.KindNumberValue Fuzz.float
                , Fuzz.map GP.KindStringValue Fuzz.string
                , Fuzz.map GP.KindBoolValue Fuzz.bool
                , Fuzz.map GP.KindStructValue fuzzStruct
                , Fuzz.map GP.KindListValue fuzzListValue
                ]
    in
    Fuzz.map GP.Value (Fuzz.map GP.ValueKind (Fuzz.maybe kindTypeFuzzer))
