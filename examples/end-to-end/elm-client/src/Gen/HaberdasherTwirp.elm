module Gen.HaberdasherTwirp exposing (..)

{-
   // Code generated protoc-gen-elmer DO NOT EDIT \\
-}

import Gen.Haberdasher
import Http
import Protobuf.Decode as PD
import Protobuf.ElmerTest
import Protobuf.Encode as PE



--  A Haberdasher makes hats for clients.


{-| MakeHat produces a hat of mysterious, randomly-selected color!
-}
haberdasher_MakeHat : (Result Http.Error Gen.Haberdasher.Hat -> msg) -> String -> Gen.Haberdasher.Size -> Cmd msg
haberdasher_MakeHat msg api data =
    Http.post
        { url = api ++ "/gen.haberdasher.Haberdasher/MakeHat"
        , body =
            Gen.Haberdasher.sizeEncoder data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Gen.Haberdasher.hatDecoder
        }
