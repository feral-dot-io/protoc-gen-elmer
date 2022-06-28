-- Generated by protoc-gen-elmgen. DO NOT EDIT!


module Gen.HaberdasherTwirp exposing (..)

import Gen.Haberdasher as Data
import Http
import Protobuf.Decode as PD
import Protobuf.Encode as PE


makeHat : (Result Http.Error Data.Hat -> msg) -> String -> Data.Size -> Cmd msg
makeHat msg api data =
    Http.post
        { url = api ++ "/gen.haberdasher.Haberdasher/MakeHat"
        , body =
            Data.sizeEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Data.hatDecoder
        }
