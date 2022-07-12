module ExampleTwirp exposing (..)

{-| Protobuf library for executing RPC methods defined in example.proto. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Bytes exposing (Bytes)
import Bytes.Encode as BE
import Dict exposing (Dict)
import Example
import Http
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE



-- We can define RPC methods and generate a Twirp client


twirpOurService_AnotherMethod :
    (Result Http.Error Example.Scalar -> msg)
    -> String
    -> Example.AllTogether
    -> Cmd msg
twirpOurService_AnotherMethod msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/example.OurService/AnotherMethod"
        , body =
            Example.encodeAllTogether data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Example.decodeScalar
        , timeout = Nothing
        , tracker = Nothing
        }


{-| Each method is an HTTP request
-}
twirpOurService_OurRpcMethod :
    (Result Http.Error Example.AllTogether -> msg)
    -> String
    -> Example.Scalar
    -> Cmd msg
twirpOurService_OurRpcMethod msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/example.OurService/OurRPCMethod"
        , body =
            Example.encodeScalar data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Example.decodeAllTogether
        , timeout = Nothing
        , tracker = Nothing
        }
