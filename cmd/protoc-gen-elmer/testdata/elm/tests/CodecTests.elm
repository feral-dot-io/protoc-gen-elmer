module CodecTests exposing (..)

import Codec
import Expect
import Protobuf.Decode as PD
import Protobuf.Encode as PE
import Test exposing (Test, test)


multiTest : Test
multiTest =
    let
        run data =
            PE.encode (Codec.encodeMulti data)
                |> PD.decode Codec.decodeMulti
                |> Expect.equal (Just data)
    in
    test "empty" <|
        \_ ->
            run Codec.emptyMulti
