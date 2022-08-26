module Gen.HaberdasherTwirp exposing (..)

{-| Protobuf library for executing RPC methods defined in package `gen.haberdasher`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Gen.Haberdasher
import Http
import Protobuf.Decode as PD
import Protobuf.Encode as PE



-- A Haberdasher makes hats for clients.


{-| MakeHat produces a hat of mysterious, randomly-selected color!
-}
twirpHaberdasher_MakeHat :
    (Result Http.Error Gen.Haberdasher.Hat -> msg)
    -> String
    -> Gen.Haberdasher.Size
    -> Cmd msg
twirpHaberdasher_MakeHat msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/gen.haberdasher.Haberdasher/MakeHat"
        , body =
            Gen.Haberdasher.encodeSize data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Gen.Haberdasher.decodeHat
        , timeout = Nothing
        , tracker = Nothing
        }
