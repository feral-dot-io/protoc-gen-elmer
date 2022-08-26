module Ex04Twirp exposing (..)

{-| Protobuf library for executing RPC methods defined in package `Ex04`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Ex04
import Google.Protobuf
import Http
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE


{-| Example API for our RPC client. No corresponding implementation.
Check out `end-to-end` example of everything working together.
-}
twirpSpeaker_HelloWorld :
    (Result Http.Error Ex04.Response -> msg)
    -> String
    -> Google.Protobuf.Empty
    -> Cmd msg
twirpSpeaker_HelloWorld msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/Ex04.Speaker/HelloWorld"
        , body =
            Google.Protobuf.toEmptyEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Ex04.decodeResponse
        , timeout = Nothing
        , tracker = Nothing
        }
