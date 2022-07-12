-- This file is part of protoc-gen-elmer.
--
-- Protoc-gen-elmer is free software: you can redistribute it and/or modify it under the terms of the GNU Lesser General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
--
-- Protoc-gen-elmer is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.
--
-- You should have received a copy of the GNU Lesser General Public License along with Protoc-gen-elmer. If not, see <https:--www.gnu.org/licenses/>.


module Protobuf.Elmer exposing
    ( BoolValue, BytesValue, DoubleValue, FloatValue, Int32Value, Int64Value, StringValue, UInt32Value, UInt64Value
    , emptyAny, emptyApi, emptyBoolValue, emptyBytes, emptyBytesValue, emptyDoubleValue, emptyDuration, emptyEmpty, emptyEnum, emptyEnumValue, emptyField, emptyFieldMask, emptyField_Cardinality, emptyField_Kind, emptyFloatValue, emptyInt32Value, emptyInt64Value, emptyListValue, emptyMethod, emptyMixin, emptyNullValue, emptyOption, emptySourceContext, emptyStringValue, emptyStruct, emptySyntax, emptyTimestamp, emptyUInt32Value, emptyUInt64Value, emptyValue, emptyXType
    , decodeBoolValue, decodeBytesValue, decodeDoubleValue, decodeFloatValue, decodeInt32Value, decodeInt64Value, decodeStringValue, decodeTimestamp, decodeUInt32Value, decodeUInt64Value, decodeValue
    , encodeAny, encodeBoolValue, encodeBytesValue, encodeDoubleValue, encodeFloatValue, encodeInt32Value, encodeInt64Value, encodeStringValue, encodeTimestamp, encodeUInt32Value, encodeUInt64Value, encodeValue
    )

{-| Helper types and functions for `protoc-gen-elmer` codegen. This module should not be used directly.

See the project on how this may be used: <https://github.com/feral-dot-io/protoc-gen-elmer>


# Well-known types

@docs BoolValue, BytesValue, DoubleValue, FloatValue, Int32Value, Int64Value, StringValue, UInt32Value, UInt64Value


# Empty (zero) vlaues

@docs emptyAny, emptyApi, emptyBoolValue, emptyBytes, emptyBytesValue, emptyDoubleValue, emptyDuration, emptyEmpty, emptyEnum, emptyEnumValue, emptyField, emptyFieldMask, emptyField_Cardinality, emptyField_Kind, emptyFloatValue, emptyInt32Value, emptyInt64Value, emptyListValue, emptyMethod, emptyMixin, emptyNullValue, emptyOption, emptySourceContext, emptyStringValue, emptyStruct, emptySyntax, emptyTimestamp, emptyUInt32Value, emptyUInt64Value, emptyValue, emptyXType


# Decoders

@docs decodeBoolValue, decodeBytesValue, decodeDoubleValue, decodeFloatValue, decodeInt32Value, decodeInt64Value, decodeStringValue, decodeTimestamp, decodeUInt32Value, decodeUInt64Value, decodeValue


# Encoders

@docs encodeAny, encodeBoolValue, encodeBytesValue, encodeDoubleValue, encodeFloatValue, encodeInt32Value, encodeInt64Value, encodeStringValue, encodeTimestamp, encodeUInt32Value, encodeUInt64Value, encodeValue

-}

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Google.Protobuf as GP
import Protobuf.Decode as PD
import Protobuf.Encode as PE
import Time



-- Types


{-| -}
type alias BoolValue =
    Maybe Bool


{-| -}
type alias BytesValue =
    Maybe Bytes


{-| -}
type alias DoubleValue =
    Maybe Float


{-| -}
type alias FloatValue =
    Maybe Float


{-| -}
type alias Int32Value =
    Maybe Int


{-| -}
type alias Int64Value =
    Maybe Int


{-| -}
type alias StringValue =
    Maybe String


{-| -}
type alias UInt32Value =
    Maybe Int


{-| -}
type alias UInt64Value =
    Maybe Int



-- Zero values


{-| -}
emptyBoolValue : BoolValue
emptyBoolValue =
    Nothing


{-| -}
emptyBytes : Bytes
emptyBytes =
    BE.encode (BE.sequence [])


{-| -}
emptyBytesValue : BytesValue
emptyBytesValue =
    Nothing


{-| -}
emptyDoubleValue : FloatValue
emptyDoubleValue =
    Nothing


{-| -}
emptyFloatValue : FloatValue
emptyFloatValue =
    Nothing


{-| -}
emptyInt32Value : Int64Value
emptyInt32Value =
    Nothing


{-| -}
emptyInt64Value : Int64Value
emptyInt64Value =
    Nothing


{-| -}
emptyStringValue : StringValue
emptyStringValue =
    Nothing


{-| -}
emptyTimestamp : Time.Posix
emptyTimestamp =
    Time.millisToPosix 0


{-| -}
emptyUInt32Value : Maybe Int
emptyUInt32Value =
    Nothing


{-| -}
emptyUInt64Value : Maybe Int
emptyUInt64Value =
    Nothing



-- Decoders


{-| -}
decodeBoolValue : PD.Decoder BoolValue
decodeBoolValue =
    decodeValue PD.bool


{-| -}
decodeBytesValue : PD.Decoder BytesValue
decodeBytesValue =
    decodeValue PD.bytes


{-| -}
decodeDoubleValue : PD.Decoder FloatValue
decodeDoubleValue =
    decodeValue PD.double


{-| -}
decodeFloatValue : PD.Decoder FloatValue
decodeFloatValue =
    decodeValue PD.float


{-| -}
decodeInt32Value : PD.Decoder Int32Value
decodeInt32Value =
    decodeValue PD.int32


{-| -}
decodeInt64Value : PD.Decoder Int64Value
decodeInt64Value =
    decodeValue PD.int32


{-| -}
decodeStringValue : PD.Decoder StringValue
decodeStringValue =
    decodeValue PD.string


{-| -}
decodeTimestamp : PD.Decoder Time.Posix
decodeTimestamp =
    GP.timestampDecoder
        |> PD.map (\t -> t.seconds * 1000 + t.nanos // 1000000)
        |> PD.map Time.millisToPosix


{-| -}
decodeUInt32Value : PD.Decoder UInt32Value
decodeUInt32Value =
    decodeValue PD.uint32


{-| -}
decodeUInt64Value : PD.Decoder UInt64Value
decodeUInt64Value =
    decodeValue PD.uint32


{-| -}
decodeValue : PD.Decoder w -> PD.Decoder (Maybe w)
decodeValue dec =
    PD.message Nothing [ PD.optional 1 dec (\v _ -> Just v) ]



-- Encoders


{-| -}
encodeAny : GP.Any -> PE.Encoder
encodeAny =
    GP.toAnyEncoder


{-| -}
encodeBoolValue : BoolValue -> PE.Encoder
encodeBoolValue =
    encodeValue PE.bool


{-| -}
encodeBytesValue : BytesValue -> PE.Encoder
encodeBytesValue =
    encodeValue PE.bytes


{-| -}
encodeDoubleValue : FloatValue -> PE.Encoder
encodeDoubleValue =
    encodeValue PE.double


{-| -}
encodeFloatValue : FloatValue -> PE.Encoder
encodeFloatValue =
    encodeValue PE.float


{-| -}
encodeInt32Value : Int32Value -> PE.Encoder
encodeInt32Value =
    encodeValue PE.int32


{-| -}
encodeInt64Value : Int64Value -> PE.Encoder
encodeInt64Value =
    encodeValue PE.int32


{-| -}
encodeStringValue : StringValue -> PE.Encoder
encodeStringValue =
    encodeValue PE.string


{-| -}
encodeTimestamp : Time.Posix -> PE.Encoder
encodeTimestamp p =
    let
        ms =
            Time.posixToMillis p
    in
    GP.toTimestampEncoder
        { seconds = ms // 1000
        , nanos = modBy 1000 ms * 1000000
        }


{-| -}
encodeUInt32Value : UInt32Value -> PE.Encoder
encodeUInt32Value =
    encodeValue PE.uint32


{-| -}
encodeUInt64Value : UInt64Value -> PE.Encoder
encodeUInt64Value =
    encodeValue PE.uint32


{-| -}
encodeValue : (v -> PE.Encoder) -> Maybe v -> PE.Encoder
encodeValue enc v =
    PE.message [ ( 1, v |> Maybe.map enc |> Maybe.withDefault PE.none ) ]



-- Zero values for Google.Protobuf pass through


{-| -}
emptyAny : GP.Any
emptyAny =
    GP.Any "" emptyBytes


{-| -}
emptyApi : GP.Api
emptyApi =
    GP.Api "" [] [] "" Nothing [] emptySyntax


{-| -}
emptyDuration : GP.Duration
emptyDuration =
    GP.Duration 0 0


{-| -}
emptyEmpty : GP.Empty
emptyEmpty =
    GP.Empty


{-| -}
emptyEnum : GP.Enum
emptyEnum =
    GP.Enum "" [] [] Nothing emptySyntax


{-| -}
emptyEnumValue : GP.EnumValue
emptyEnumValue =
    GP.EnumValue "" 0 []


{-| -}
emptyField : GP.Field
emptyField =
    GP.Field emptyField_Kind emptyField_Cardinality 0 "" "" 0 False [] "" ""


{-| -}
emptyField_Cardinality : GP.Cardinality
emptyField_Cardinality =
    GP.CardinalityUnknown


{-| -}
emptyField_Kind : GP.Kind
emptyField_Kind =
    GP.TypeUnknown


{-| -}
emptyFieldMask : GP.FieldMask
emptyFieldMask =
    GP.FieldMask []


{-| -}
emptyListValue : GP.ListValue
emptyListValue =
    GP.ListValue (GP.ListValueValues [])


{-| -}
emptyMethod : GP.Method
emptyMethod =
    GP.Method "" "" False "" False [] emptySyntax


{-| -}
emptyMixin : GP.Mixin
emptyMixin =
    GP.Mixin "" ""


{-| -}
emptyNullValue : GP.NullValue
emptyNullValue =
    GP.NullValue


{-| -}
emptyOption : GP.Option
emptyOption =
    GP.Option "" Nothing


{-| -}
emptySourceContext : GP.SourceContext
emptySourceContext =
    GP.SourceContext ""


{-| -}
emptyStruct : GP.Struct
emptyStruct =
    GP.Struct (GP.StructFields Dict.empty)


{-| -}
emptySyntax : GP.Syntax
emptySyntax =
    GP.SyntaxProto2


{-| -}
emptyXType : GP.Type
emptyXType =
    GP.Type "" [] [] [] Nothing emptySyntax


{-| -}
emptyValue : GP.Value
emptyValue =
    GP.Value (GP.ValueKind Nothing)
