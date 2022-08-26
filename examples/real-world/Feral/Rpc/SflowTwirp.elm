module Feral.Rpc.SflowTwirp exposing (..)

{-| Protobuf library for executing RPC methods defined in package `feral.rpc.sflow`. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Dict exposing (Dict)
import Feral.Rpc.Sflow
import Google.Protobuf
import Http
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE
import Time



-- RPC


twirpSflow_ListAgents :
    (Result Http.Error Feral.Rpc.Sflow.ListAgentsResponse -> msg)
    -> String
    -> Feral.Rpc.Sflow.ListAgentsRequest
    -> Cmd msg
twirpSflow_ListAgents msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/feral.rpc.sflow.Sflow/ListAgents"
        , body =
            Feral.Rpc.Sflow.encodeListAgentsRequest data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Feral.Rpc.Sflow.decodeListAgentsResponse
        , timeout = Nothing
        , tracker = Nothing
        }


twirpSflow_ListKnownTags :
    (Result Http.Error Feral.Rpc.Sflow.ListKnownTagsResponse -> msg)
    -> String
    -> Feral.Rpc.Sflow.ListKnownTagsRequest
    -> Cmd msg
twirpSflow_ListKnownTags msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/feral.rpc.sflow.Sflow/ListKnownTags"
        , body =
            Feral.Rpc.Sflow.encodeListKnownTagsRequest data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Feral.Rpc.Sflow.decodeListKnownTagsResponse
        , timeout = Nothing
        , tracker = Nothing
        }


twirpSflow_ListRates :
    (Result Http.Error Feral.Rpc.Sflow.ListRatesResponse -> msg)
    -> String
    -> Feral.Rpc.Sflow.ListRatesRequest
    -> Cmd msg
twirpSflow_ListRates msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/feral.rpc.sflow.Sflow/ListRates"
        , body =
            Feral.Rpc.Sflow.encodeListRatesRequest data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Feral.Rpc.Sflow.decodeListRatesResponse
        , timeout = Nothing
        , tracker = Nothing
        }


twirpSflow_ListSamples :
    (Result Http.Error Feral.Rpc.Sflow.ListSamplesResponse -> msg)
    -> String
    -> Feral.Rpc.Sflow.ListSamplesRequest
    -> Cmd msg
twirpSflow_ListSamples msg api data =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = api ++ "/feral.rpc.sflow.Sflow/ListSamples"
        , body =
            Feral.Rpc.Sflow.encodeListSamplesRequest data
                |> PE.encode
                |> Http.bytesBody "application/protobuf"
        , expect = PD.expectBytes msg Feral.Rpc.Sflow.decodeListSamplesResponse
        , timeout = Nothing
        , tracker = Nothing
        }
