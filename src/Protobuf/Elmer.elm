module Protobuf.Elmer exposing
    ( BoolValue
    , BytesValue
    , DoubleValue
    , FloatValue
    , Int32Value
    , Int64Value
    , StringValue
    , Timestamp
    , UInt32Value
    , UInt64Value
    , anyEncoder
    , boolValueDecoder
    , boolValueEncoder
    , bytesValueDecoder
    , bytesValueEncoder
    , doubleValueDecoder
    , doubleValueEncoder
    , emptyAny
    , emptyApi
    , emptyBoolValue
    , emptyBytes
    , emptyBytesValue
    , emptyDoubleValue
    , emptyDuration
    , emptyEmpty
    , emptyEnum
    , emptyEnumValue
    , emptyField
    , emptyFieldMask
    , emptyField_Cardinality
    , emptyField_Kind
    , emptyFloatValue
    , emptyInt32Value
    , emptyInt64Value
    , emptyListValue
    , emptyMethod
    , emptyMixin
    , emptyNullValue
    , emptyOption
    , emptySourceContext
    , emptyStringValue
    , emptyStruct
    , emptySyntax
    , emptyTimestamp
    , emptyType
    , emptyUInt32Value
    , emptyUInt64Value
    , emptyValue
    , floatValueDecoder
    , floatValueEncoder
    , int32ValueDecoder
    , int32ValueEncoder
    , int64ValueDecoder
    , int64ValueEncoder
    , stringValueDecoder
    , stringValueEncoder
    , timestampDecoder
    , timestampEncoder
    , uInt32ValueDecoder
    , uInt32ValueEncoder
    , uInt64ValueDecoder
    , uInt64ValueEncoder
    )

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict
import Google.Protobuf as GP
import Protobuf.Decode as PD
import Protobuf.Encode as PE
import Time



-- Types


type alias BoolValue =
    Maybe Bool


type alias BytesValue =
    Maybe Bytes


type alias DoubleValue =
    Maybe Float


type alias FloatValue =
    Maybe Float


type alias Int32Value =
    Maybe Int


type alias Int64Value =
    Maybe Int


type alias StringValue =
    Maybe String


type alias Timestamp =
    Time.Posix


type alias UInt32Value =
    Maybe Int


type alias UInt64Value =
    Maybe Int



-- Zero values


emptyBoolValue : BoolValue
emptyBoolValue =
    Nothing


emptyBytes : Bytes
emptyBytes =
    BE.encode (BE.sequence [])


emptyBytesValue : BytesValue
emptyBytesValue =
    Nothing


emptyDoubleValue : FloatValue
emptyDoubleValue =
    Nothing


emptyFloatValue : FloatValue
emptyFloatValue =
    Nothing


emptyInt32Value : Int64Value
emptyInt32Value =
    Nothing


emptyInt64Value : Int64Value
emptyInt64Value =
    Nothing


emptyStringValue : StringValue
emptyStringValue =
    Nothing


emptyTimestamp : Time.Posix
emptyTimestamp =
    Time.millisToPosix 0


emptyUInt32Value : Maybe Int
emptyUInt32Value =
    Nothing


emptyUInt64Value : Maybe Int
emptyUInt64Value =
    Nothing



-- Decoders


boolValueDecoder : PD.Decoder BoolValue
boolValueDecoder =
    valueDecoder PD.bool


bytesValueDecoder : PD.Decoder BytesValue
bytesValueDecoder =
    valueDecoder PD.bytes


doubleValueDecoder : PD.Decoder FloatValue
doubleValueDecoder =
    valueDecoder PD.double


floatValueDecoder : PD.Decoder FloatValue
floatValueDecoder =
    valueDecoder PD.float


int32ValueDecoder : PD.Decoder Int32Value
int32ValueDecoder =
    valueDecoder PD.int32


int64ValueDecoder : PD.Decoder Int64Value
int64ValueDecoder =
    valueDecoder PD.int32


stringValueDecoder : PD.Decoder StringValue
stringValueDecoder =
    valueDecoder PD.string


timestampDecoder : PD.Decoder Timestamp
timestampDecoder =
    GP.timestampDecoder
        |> PD.map (\t -> t.seconds * 1000 + t.nanos // 1000000)
        |> PD.map Time.millisToPosix


uInt32ValueDecoder : PD.Decoder UInt32Value
uInt32ValueDecoder =
    valueDecoder PD.uint32


uInt64ValueDecoder : PD.Decoder UInt64Value
uInt64ValueDecoder =
    valueDecoder PD.uint32


valueDecoder : PD.Decoder w -> PD.Decoder (Maybe w)
valueDecoder dec =
    PD.message Nothing [ PD.optional 1 dec (\v _ -> Just v) ]



-- Encoders


anyEncoder : GP.Any -> PE.Encoder
anyEncoder =
    GP.toAnyEncoder


boolValueEncoder : BoolValue -> PE.Encoder
boolValueEncoder =
    valueEncoder PE.bool


bytesValueEncoder : BytesValue -> PE.Encoder
bytesValueEncoder =
    valueEncoder PE.bytes


doubleValueEncoder : FloatValue -> PE.Encoder
doubleValueEncoder =
    valueEncoder PE.double


floatValueEncoder : FloatValue -> PE.Encoder
floatValueEncoder =
    valueEncoder PE.float


int32ValueEncoder : Int32Value -> PE.Encoder
int32ValueEncoder =
    valueEncoder PE.int32


int64ValueEncoder : Int64Value -> PE.Encoder
int64ValueEncoder =
    valueEncoder PE.int32


stringValueEncoder : StringValue -> PE.Encoder
stringValueEncoder =
    valueEncoder PE.string


timestampEncoder : Time.Posix -> PE.Encoder
timestampEncoder p =
    let
        ms =
            Time.posixToMillis p
    in
    GP.toTimestampEncoder
        { seconds = ms // 1000
        , nanos = modBy 1000 ms * 1000000
        }


uInt32ValueEncoder : UInt32Value -> PE.Encoder
uInt32ValueEncoder =
    valueEncoder PE.uint32


uInt64ValueEncoder : UInt64Value -> PE.Encoder
uInt64ValueEncoder =
    valueEncoder PE.uint32


valueEncoder : (v -> PE.Encoder) -> Maybe v -> PE.Encoder
valueEncoder enc v =
    PE.message [ ( 1, v |> Maybe.map enc |> Maybe.withDefault PE.none ) ]



-- Zero values for Google.Protobuf pass through


emptyAny : GP.Any
emptyAny =
    GP.Any "" emptyBytes


emptyApi : GP.Api
emptyApi =
    GP.Api "" [] [] "" Nothing [] emptySyntax


emptyDuration : GP.Duration
emptyDuration =
    GP.Duration 0 0


emptyEmpty : GP.Empty
emptyEmpty =
    GP.Empty


emptyEnum : GP.Enum
emptyEnum =
    GP.Enum "" [] [] Nothing emptySyntax


emptyEnumValue : GP.EnumValue
emptyEnumValue =
    GP.EnumValue "" 0 []


emptyField : GP.Field
emptyField =
    GP.Field emptyField_Kind emptyField_Cardinality 0 "" "" 0 False [] "" ""


emptyField_Cardinality : GP.Cardinality
emptyField_Cardinality =
    GP.CardinalityUnknown


emptyField_Kind : GP.Kind
emptyField_Kind =
    GP.TypeUnknown


emptyFieldMask : GP.FieldMask
emptyFieldMask =
    GP.FieldMask []


emptyListValue : GP.ListValue
emptyListValue =
    GP.ListValue (GP.ListValueValues [])


emptyMethod : GP.Method
emptyMethod =
    GP.Method "" "" False "" False [] emptySyntax


emptyMixin : GP.Mixin
emptyMixin =
    GP.Mixin "" ""


emptyNullValue : GP.NullValue
emptyNullValue =
    GP.NullValue


emptyOption : GP.Option
emptyOption =
    GP.Option "" Nothing


emptySourceContext : GP.SourceContext
emptySourceContext =
    GP.SourceContext ""


emptyStruct : GP.Struct
emptyStruct =
    GP.Struct (GP.StructFields Dict.empty)


emptySyntax : GP.Syntax
emptySyntax =
    GP.SyntaxProto2


emptyType : GP.Type
emptyType =
    GP.Type "" [] [] [] Nothing emptySyntax


emptyValue : GP.Value
emptyValue =
    GP.Value (GP.ValueKind Nothing)
