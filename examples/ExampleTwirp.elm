module ExampleTwirp exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Example
import Http
import Protobuf.Decode as PD
import Protobuf.Encode as PE



--  We can define RPC methods and generate a Twirp client


ourService_AnotherMethod : (Result Http.Error Example.Scalar -> msg) -> String -> Example.AllTogether -> Cmd msg
ourService_AnotherMethod msg api data =
    Http.post
        { url = api ++ "/example.OurService/AnotherMethod"
        , body =
            Example.allTogetherEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Example.scalarDecoder
        }


{-| Each method is an HTTP request
-}
ourService_OurRpcMethod : (Result Http.Error Example.AllTogether -> msg) -> String -> Example.Scalar -> Cmd msg
ourService_OurRpcMethod msg api data =
    Http.post
        { url = api ++ "/example.OurService/OurRPCMethod"
        , body =
            Example.scalarEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Example.allTogetherDecoder
        }
