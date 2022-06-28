-- Generated by protoc-gen-elmgen. DO NOT EDIT!


module ExampleTwirp exposing (..)

import Example as Data
import Http
import Protobuf.Decode as PD
import Protobuf.Encode as PE


anotherMethod : (Result Http.Error Data.Scalar -> msg) -> String -> Data.AllTogether -> Cmd msg
anotherMethod msg api data =
    Http.post
        { url = api ++ "/example.OurService/AnotherMethod"
        , body =
            Data.allTogetherEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Data.scalarDecoder
        }


ourRpcMethod : (Result Http.Error Data.AllTogether -> msg) -> String -> Data.Scalar -> Cmd msg
ourRpcMethod msg api data =
    Http.post
        { url = api ++ "/example.OurService/OurRPCMethod"
        , body =
            Data.scalarEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Data.allTogetherDecoder
        }