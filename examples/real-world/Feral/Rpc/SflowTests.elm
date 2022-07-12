module Feral.Rpc.SflowTests exposing (..)

{-| Protobuf library for testing structures found in sflow.proto. This file was generated automatically by `protoc-gen-elmer`. See the base file for more information. Do not edit.
-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Dict exposing (Dict)
import Expect
import Feral.Rpc.Sflow
import Fuzz exposing (Fuzzer)
import Google.Protobuf
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.ElmerTests
import Protobuf.Encode as PE
import Test exposing (Test, fuzz, test)
import Time


fuzzAgent_Boot : Fuzzer Feral.Rpc.Sflow.Agent_Boot
fuzzAgent_Boot =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.Agent_Disk
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Pxe
        ]


fuzzAgent_Role : Fuzzer Feral.Rpc.Sflow.Agent_Role
fuzzAgent_Role =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.Agent_Opaque
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Router
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Server
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Oob
        ]


fuzzAgent_Slot : Fuzzer Feral.Rpc.Sflow.Agent_Slot
fuzzAgent_Slot =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.Agent_NotSlot
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Capacity
        , Fuzz.constant Feral.Rpc.Sflow.Agent_Capability
        ]


fuzzSampleTag_Type : Fuzzer Feral.Rpc.Sflow.SampleTag_Type
fuzzSampleTag_Type =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.SampleTag_Unkown
        , Fuzz.constant Feral.Rpc.Sflow.SampleTag_Filter
        , Fuzz.constant Feral.Rpc.Sflow.SampleTag_Top
        ]


fuzzSample_Group : Fuzzer Feral.Rpc.Sflow.Sample_Group
fuzzSample_Group =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.Sample_NoGroup
        , Fuzz.constant Feral.Rpc.Sflow.Sample_Role
        , Fuzz.constant Feral.Rpc.Sflow.Sample_Agent
        , Fuzz.constant Feral.Rpc.Sflow.Sample_Interface
        , Fuzz.constant Feral.Rpc.Sflow.Sample_Input
        , Fuzz.constant Feral.Rpc.Sflow.Sample_Output
        , Fuzz.constant Feral.Rpc.Sflow.Sample_InputDiscard
        , Fuzz.constant Feral.Rpc.Sflow.Sample_OutputDiscard
        , Fuzz.constant Feral.Rpc.Sflow.Sample_L3Protocol
        , Fuzz.constant Feral.Rpc.Sflow.Sample_L4Protocol
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcMac
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcAsn
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcNextAsn
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcPrefix
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcIp
        , Fuzz.constant Feral.Rpc.Sflow.Sample_SrcPort
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstMac
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstAsn
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstNextAsn
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstPrefix
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstIp
        , Fuzz.constant Feral.Rpc.Sflow.Sample_DstPort
        ]


fuzzState_Duplex : Fuzzer Feral.Rpc.Sflow.State_Duplex
fuzzState_Duplex =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.State_UnknownDuplex
        , Fuzz.constant Feral.Rpc.Sflow.State_FullDuplex
        , Fuzz.constant Feral.Rpc.Sflow.State_HalfDuplex
        , Fuzz.constant Feral.Rpc.Sflow.State_InDuplex
        , Fuzz.constant Feral.Rpc.Sflow.State_OutDuplex
        ]


fuzzState_Oper : Fuzzer Feral.Rpc.Sflow.State_Oper
fuzzState_Oper =
    Fuzz.oneOf
        [ Fuzz.constant Feral.Rpc.Sflow.State_NotUp
        , Fuzz.constant Feral.Rpc.Sflow.State_Up
        ]


fuzzAgent : Fuzzer Feral.Rpc.Sflow.Agent
fuzzAgent =
    Fuzz.map Feral.Rpc.Sflow.Agent
        Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap fuzzAgent_Role
        |> Fuzz.andMap fuzzAgent_Slot
        |> Fuzz.andMap fuzzAgent_Boot
        |> Fuzz.andMap Fuzz.string


fuzzInterface : Fuzzer Feral.Rpc.Sflow.Interface
fuzzInterface =
    Fuzz.map Feral.Rpc.Sflow.Interface
        Fuzz.string
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap fuzzState
        |> Fuzz.andMap fuzzRates


fuzzKnownTags : Fuzzer Feral.Rpc.Sflow.KnownTags
fuzzKnownTags =
    Fuzz.map Feral.Rpc.Sflow.KnownTags
        (Fuzz.list Fuzz.string)
        |> Fuzz.andMap (Fuzz.list Protobuf.ElmerTests.fuzzUInt32)
        |> Fuzz.andMap (Fuzz.list Protobuf.ElmerTests.fuzzUInt32)
        |> Fuzz.andMap (Fuzz.list Protobuf.ElmerTests.fuzzUInt32)
        |> Fuzz.andMap (Fuzz.list Protobuf.ElmerTests.fuzzUInt32)


fuzzListAgentsRequest : Fuzzer Feral.Rpc.Sflow.ListAgentsRequest
fuzzListAgentsRequest =
    Fuzz.constant Feral.Rpc.Sflow.ListAgentsRequest


fuzzListAgentsResponse : Fuzzer Feral.Rpc.Sflow.ListAgentsResponse
fuzzListAgentsResponse =
    Fuzz.map Feral.Rpc.Sflow.ListAgentsResponse
        (Fuzz.list fuzzAgent)
        |> Fuzz.andMap (Fuzz.list fuzzInterface)


fuzzListKnownTagsRequest : Fuzzer Feral.Rpc.Sflow.ListKnownTagsRequest
fuzzListKnownTagsRequest =
    Fuzz.constant Feral.Rpc.Sflow.ListKnownTagsRequest


fuzzListKnownTagsResponse : Fuzzer Feral.Rpc.Sflow.ListKnownTagsResponse
fuzzListKnownTagsResponse =
    Fuzz.map Feral.Rpc.Sflow.ListKnownTagsResponse
        fuzzKnownTags


fuzzListRatesRequest : Fuzzer Feral.Rpc.Sflow.ListRatesRequest
fuzzListRatesRequest =
    let
        fuzzListRatesRequest_Agent =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzListRatesRequest_IfIndex =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]
    in
    Fuzz.map Feral.Rpc.Sflow.ListRatesRequest
        fuzzWindow
        |> Fuzz.andMap fuzzAgent_Role
        |> Fuzz.andMap (Fuzz.maybe fuzzListRatesRequest_Agent)
        |> Fuzz.andMap (Fuzz.maybe fuzzListRatesRequest_IfIndex)


fuzzListRatesResponse : Fuzzer Feral.Rpc.Sflow.ListRatesResponse
fuzzListRatesResponse =
    Fuzz.map Feral.Rpc.Sflow.ListRatesResponse
        (Fuzz.list fuzzRates)


fuzzListSamplesRequest : Fuzzer Feral.Rpc.Sflow.ListSamplesRequest
fuzzListSamplesRequest =
    Fuzz.map Feral.Rpc.Sflow.ListSamplesRequest
        fuzzWindow
        |> Fuzz.andMap fuzzTagFilter
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap (Fuzz.list fuzzSample_Group)


fuzzListSamplesResponse : Fuzzer Feral.Rpc.Sflow.ListSamplesResponse
fuzzListSamplesResponse =
    Fuzz.map Feral.Rpc.Sflow.ListSamplesResponse
        fuzzSeries


fuzzRates : Fuzzer Feral.Rpc.Sflow.Rates
fuzzRates =
    Fuzz.map Feral.Rpc.Sflow.Rates
        Protobuf.ElmerTests.fuzzTimestamp
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzDoubleValue


fuzzSample : Fuzzer Feral.Rpc.Sflow.Sample
fuzzSample =
    Fuzz.map Feral.Rpc.Sflow.Sample
        Fuzz.float
        |> Fuzz.andMap Fuzz.float


fuzzSampleTag : Fuzzer Feral.Rpc.Sflow.SampleTag
fuzzSampleTag =
    Fuzz.map Feral.Rpc.Sflow.SampleTag
        fuzzSampleTag_Type
        |> Fuzz.andMap fuzzTagFilter


fuzzSamples : Fuzzer Feral.Rpc.Sflow.Samples
fuzzSamples =
    Fuzz.map Feral.Rpc.Sflow.Samples
        Protobuf.ElmerTests.fuzzTimestamp
        |> Fuzz.andMap (Fuzz.map Dict.fromList (Fuzz.list (Fuzz.tuple ( Protobuf.ElmerTests.fuzzInt32, fuzzSample ))))


fuzzSeries : Fuzzer Feral.Rpc.Sflow.Series
fuzzSeries =
    Fuzz.map Feral.Rpc.Sflow.Series
        (Fuzz.map Dict.fromList (Fuzz.list (Fuzz.tuple ( Protobuf.ElmerTests.fuzzInt32, fuzzSampleTag ))))
        |> Fuzz.andMap (Fuzz.list fuzzSamples)


fuzzState : Fuzzer Feral.Rpc.Sflow.State
fuzzState =
    Fuzz.map Feral.Rpc.Sflow.State
        Fuzz.string
        |> Fuzz.andMap Fuzz.string
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzTimestamp
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32
        |> Fuzz.andMap Fuzz.float
        |> Fuzz.andMap fuzzState_Duplex
        |> Fuzz.andMap Fuzz.bool
        |> Fuzz.andMap fuzzState_Oper
        |> Fuzz.andMap fuzzState_Oper


fuzzTagFilter : Fuzzer Feral.Rpc.Sflow.TagFilter
fuzzTagFilter =
    let
        fuzzTagFilter_Agent =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_IfIndex =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_Input =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_Output =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_InputDiscard =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_OutputDiscard =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_L3Protocol =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_L4Protocol =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_SrcMac =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_SrcPrefix =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_SrcIp =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_SrcAsn =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_SrcNextAsn =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_SrcPort =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_DstMac =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_DstPrefix =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_DstIp =
            Fuzz.oneOf
                [ Fuzz.string
                ]

        fuzzTagFilter_DstAsn =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_DstNextAsn =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]

        fuzzTagFilter_DstPort =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzUInt32
                ]
    in
    Fuzz.map Feral.Rpc.Sflow.TagFilter
        fuzzAgent_Role
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_Agent)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_IfIndex)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_Input)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_Output)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_InputDiscard)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_OutputDiscard)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_L3Protocol)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_L4Protocol)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcMac)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcPrefix)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcIp)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcAsn)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcNextAsn)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_SrcPort)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstMac)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstPrefix)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstIp)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstAsn)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstNextAsn)
        |> Fuzz.andMap (Fuzz.maybe fuzzTagFilter_DstPort)


fuzzWindow : Fuzzer Feral.Rpc.Sflow.Window
fuzzWindow =
    let
        fuzzWindow_Before =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzTimestamp
                ]

        fuzzWindow_Interval =
            Fuzz.oneOf
                [ Protobuf.ElmerTests.fuzzDuration
                ]
    in
    Fuzz.map Feral.Rpc.Sflow.Window
        (Fuzz.maybe fuzzWindow_Before)
        |> Fuzz.andMap (Fuzz.maybe fuzzWindow_Interval)
        |> Fuzz.andMap Protobuf.ElmerTests.fuzzUInt32


testAgent : Test
testAgent =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeAgent Feral.Rpc.Sflow.encodeAgent
    in
    Test.describe "encode then decode Agent"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyAgent)
        , fuzz fuzzAgent "fuzzer" run
        ]


testInterface : Test
testInterface =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeInterface Feral.Rpc.Sflow.encodeInterface
    in
    Test.describe "encode then decode Interface"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyInterface)
        , fuzz fuzzInterface "fuzzer" run
        ]


testKnownTags : Test
testKnownTags =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeKnownTags Feral.Rpc.Sflow.encodeKnownTags
    in
    Test.describe "encode then decode KnownTags"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyKnownTags)
        , fuzz fuzzKnownTags "fuzzer" run
        ]


testListAgentsRequest : Test
testListAgentsRequest =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListAgentsRequest Feral.Rpc.Sflow.encodeListAgentsRequest
    in
    Test.describe "encode then decode ListAgentsRequest"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListAgentsRequest)
        , fuzz fuzzListAgentsRequest "fuzzer" run
        ]


testListAgentsResponse : Test
testListAgentsResponse =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListAgentsResponse Feral.Rpc.Sflow.encodeListAgentsResponse
    in
    Test.describe "encode then decode ListAgentsResponse"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListAgentsResponse)
        , fuzz fuzzListAgentsResponse "fuzzer" run
        ]


testListKnownTagsRequest : Test
testListKnownTagsRequest =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListKnownTagsRequest Feral.Rpc.Sflow.encodeListKnownTagsRequest
    in
    Test.describe "encode then decode ListKnownTagsRequest"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListKnownTagsRequest)
        , fuzz fuzzListKnownTagsRequest "fuzzer" run
        ]


testListKnownTagsResponse : Test
testListKnownTagsResponse =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListKnownTagsResponse Feral.Rpc.Sflow.encodeListKnownTagsResponse
    in
    Test.describe "encode then decode ListKnownTagsResponse"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListKnownTagsResponse)
        , fuzz fuzzListKnownTagsResponse "fuzzer" run
        ]


testListRatesRequest : Test
testListRatesRequest =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListRatesRequest Feral.Rpc.Sflow.encodeListRatesRequest
    in
    Test.describe "encode then decode ListRatesRequest"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListRatesRequest)
        , fuzz fuzzListRatesRequest "fuzzer" run
        ]


testListRatesResponse : Test
testListRatesResponse =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListRatesResponse Feral.Rpc.Sflow.encodeListRatesResponse
    in
    Test.describe "encode then decode ListRatesResponse"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListRatesResponse)
        , fuzz fuzzListRatesResponse "fuzzer" run
        ]


testListSamplesRequest : Test
testListSamplesRequest =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListSamplesRequest Feral.Rpc.Sflow.encodeListSamplesRequest
    in
    Test.describe "encode then decode ListSamplesRequest"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListSamplesRequest)
        , fuzz fuzzListSamplesRequest "fuzzer" run
        ]


testListSamplesResponse : Test
testListSamplesResponse =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeListSamplesResponse Feral.Rpc.Sflow.encodeListSamplesResponse
    in
    Test.describe "encode then decode ListSamplesResponse"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyListSamplesResponse)
        , fuzz fuzzListSamplesResponse "fuzzer" run
        ]


testRates : Test
testRates =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeRates Feral.Rpc.Sflow.encodeRates
    in
    Test.describe "encode then decode Rates"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyRates)
        , fuzz fuzzRates "fuzzer" run
        ]


testSample : Test
testSample =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSample Feral.Rpc.Sflow.encodeSample
    in
    Test.describe "encode then decode Sample"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySample)
        , fuzz fuzzSample "fuzzer" run
        ]


testSampleTag : Test
testSampleTag =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSampleTag Feral.Rpc.Sflow.encodeSampleTag
    in
    Test.describe "encode then decode SampleTag"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySampleTag)
        , fuzz fuzzSampleTag "fuzzer" run
        ]


testSamples : Test
testSamples =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSamples Feral.Rpc.Sflow.encodeSamples
    in
    Test.describe "encode then decode Samples"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySamples)
        , fuzz fuzzSamples "fuzzer" run
        ]


testSeries : Test
testSeries =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSeries Feral.Rpc.Sflow.encodeSeries
    in
    Test.describe "encode then decode Series"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySeries)
        , fuzz fuzzSeries "fuzzer" run
        ]


testState : Test
testState =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeState Feral.Rpc.Sflow.encodeState
    in
    Test.describe "encode then decode State"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyState)
        , fuzz fuzzState "fuzzer" run
        ]


testTagFilter : Test
testTagFilter =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeTagFilter Feral.Rpc.Sflow.encodeTagFilter
    in
    Test.describe "encode then decode TagFilter"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyTagFilter)
        , fuzz fuzzTagFilter "fuzzer" run
        ]


testWindow : Test
testWindow =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeWindow Feral.Rpc.Sflow.encodeWindow
    in
    Test.describe "encode then decode Window"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyWindow)
        , fuzz fuzzWindow "fuzzer" run
        ]


testAgent_Boot : Test
testAgent_Boot =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeAgent_Boot Feral.Rpc.Sflow.encodeAgent_Boot
    in
    Test.describe "encode then decode Agent_Boot"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyAgent_Boot)
        , fuzz fuzzAgent_Boot "fuzzer" run
        ]


testAgent_Role : Test
testAgent_Role =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeAgent_Role Feral.Rpc.Sflow.encodeAgent_Role
    in
    Test.describe "encode then decode Agent_Role"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyAgent_Role)
        , fuzz fuzzAgent_Role "fuzzer" run
        ]


testAgent_Slot : Test
testAgent_Slot =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeAgent_Slot Feral.Rpc.Sflow.encodeAgent_Slot
    in
    Test.describe "encode then decode Agent_Slot"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyAgent_Slot)
        , fuzz fuzzAgent_Slot "fuzzer" run
        ]


testSampleTag_Type : Test
testSampleTag_Type =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSampleTag_Type Feral.Rpc.Sflow.encodeSampleTag_Type
    in
    Test.describe "encode then decode SampleTag_Type"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySampleTag_Type)
        , fuzz fuzzSampleTag_Type "fuzzer" run
        ]


testSample_Group : Test
testSample_Group =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeSample_Group Feral.Rpc.Sflow.encodeSample_Group
    in
    Test.describe "encode then decode Sample_Group"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptySample_Group)
        , fuzz fuzzSample_Group "fuzzer" run
        ]


testState_Duplex : Test
testState_Duplex =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeState_Duplex Feral.Rpc.Sflow.encodeState_Duplex
    in
    Test.describe "encode then decode State_Duplex"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyState_Duplex)
        , fuzz fuzzState_Duplex "fuzzer" run
        ]


testState_Oper : Test
testState_Oper =
    let
        run =
            Protobuf.ElmerTests.runTest Feral.Rpc.Sflow.decodeState_Oper Feral.Rpc.Sflow.encodeState_Oper
    in
    Test.describe "encode then decode State_Oper"
        [ test "empty" (\_ -> run Feral.Rpc.Sflow.emptyState_Oper)
        , fuzz fuzzState_Oper "fuzzer" run
        ]
