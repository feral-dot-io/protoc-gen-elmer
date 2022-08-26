module Feral.Rpc.Sflow exposing (..)

{-| Protobuf library for decoding and encoding structures found in package `feral.rpc.sflow` along with helpers. This file was generated automatically by `protoc-gen-elmer`. Do not edit.

Records:

  - Agent
  - Interface
  - KnownTags
  - ListAgentsRequest
  - ListAgentsResponse
  - ListKnownTagsRequest
  - ListKnownTagsResponse
  - ListRatesRequest
  - ListRatesResponse
  - ListSamplesRequest
  - ListSamplesResponse
  - Rates
  - Sample
  - SampleTag
  - Samples
  - Series
  - State
  - TagFilter
  - Window

Unions:

  - Agent\_Boot
  - Agent\_Role
  - Agent\_Slot
  - SampleTag\_Type
  - Sample\_Group
  - State\_Duplex
  - State\_Oper

Each type defined has a: decoder, encoder and an empty (zero value) function. In addition to this enums have valuesOf, to and from (string) functions. All functions take the form `decodeDerivedIdent` where `decode` is the purpose and `DerivedIdent` comes from the Protobuf ident.

Elm identifiers are derived directly from the Protobuf ID (a full ident). The package maps to a module and the rest of the ID is the type. Since Protobuf names are hierachical (separated by a dot `.`), each namespace is mapped to an underscore `_` in an Elm ID. A Protobuf namespaced ident (parts between a dot `.`) are then cased to follow Elm naming conventions and do not include any undescores `_`. For example the enum `my.pkg.MyMessage.URLOptions` maps to the Elm module `My.Pkg` with ID `MyMessage_UrlOptions`.


# Types

@docs Agent, Interface, KnownTags, ListAgentsRequest, ListAgentsResponse, ListKnownTagsRequest, ListKnownTagsResponse, ListRatesRequest, ListRatesResponse, ListSamplesRequest, ListSamplesResponse, Rates, Sample, SampleTag, Samples, Series, State, TagFilter, Window, Agent_Boot, Agent_Role, Agent_Slot, SampleTag_Type, Sample_Group, State_Duplex, State_Oper


# Empty (zero values)

@docs emptyAgent, emptyInterface, emptyKnownTags, emptyListAgentsRequest, emptyListAgentsResponse, emptyListKnownTagsRequest, emptyListKnownTagsResponse, emptyListRatesRequest, emptyListRatesResponse, emptyListSamplesRequest, emptyListSamplesResponse, emptyRates, emptySample, emptySampleTag, emptySamples, emptySeries, emptyState, emptyTagFilter, emptyWindow, emptyAgent_Boot, emptyAgent_Role, emptyAgent_Slot, emptySampleTag_Type, emptySample_Group, emptyState_Duplex, emptyState_Oper


# Enum valuesOf

@docs valuesOfAgent_Boot, valuesOfAgent_Role, valuesOfAgent_Slot, valuesOfSampleTag_Type, valuesOfSample_Group, valuesOfState_Duplex, valuesOfState_Oper


# Enum and String converters

@docs fromAgent_Boot, toAgent_Boot, fromAgent_Role, toAgent_Role, fromAgent_Slot, toAgent_Slot, fromSampleTag_Type, toSampleTag_Type, fromSample_Group, toSample_Group, fromState_Duplex, toState_Duplex, fromState_Oper, toState_Oper


# Decoders

@docs decodeAgent, decodeInterface, decodeKnownTags, decodeListAgentsRequest, decodeListAgentsResponse, decodeListKnownTagsRequest, decodeListKnownTagsResponse, decodeListRatesRequest, decodeListRatesResponse, decodeListSamplesRequest, decodeListSamplesResponse, decodeRates, decodeSample, decodeSampleTag, decodeSamples, decodeSeries, decodeState, decodeTagFilter, decodeWindow, decodeAgent_Boot, decodeAgent_Role, decodeAgent_Slot, decodeSampleTag_Type, decodeSample_Group, decodeState_Duplex, decodeState_Oper


# Encoders

@docs encodeAgent, encodeInterface, encodeKnownTags, encodeListAgentsRequest, encodeListAgentsResponse, encodeListKnownTagsRequest, encodeListKnownTagsResponse, encodeListRatesRequest, encodeListRatesResponse, encodeListSamplesRequest, encodeListSamplesResponse, encodeRates, encodeSample, encodeSampleTag, encodeSamples, encodeSeries, encodeState, encodeTagFilter, encodeWindow, encodeAgent_Boot, encodeAgent_Role, encodeAgent_Slot, encodeSampleTag_Type, encodeSample_Group, encodeState_Duplex, encodeState_Oper

-}

-- // Code generated protoc-gen-elmer DO NOT EDIT \\

import Dict exposing (Dict)
import Google.Protobuf
import Protobuf.Decode as PD
import Protobuf.Elmer
import Protobuf.Encode as PE
import Time


type
    Agent_Boot
    -- TODO: pick a better first option to represent ???
    = Agent_Disk
    | Agent_Pxe


type Agent_Role
    = Agent_Opaque
    | Agent_Router
    | Agent_Server
    | Agent_Oob


type Agent_Slot
    = Agent_NotSlot
    | Agent_Capacity
    | Agent_Capability


type SampleTag_Type
    = SampleTag_Unkown
    | SampleTag_Filter
    | SampleTag_Top -- DROPPED = 3;


type Sample_Group
    = Sample_NoGroup
    | Sample_Role
    | Sample_Agent
    | Sample_Interface
    | Sample_Input
    | Sample_Output
    | Sample_InputDiscard
    | Sample_OutputDiscard
    | Sample_L3Protocol
    | Sample_L4Protocol
    | Sample_SrcMac
    | Sample_SrcAsn
    | Sample_SrcNextAsn
    | Sample_SrcPrefix
    | Sample_SrcIp
    | Sample_SrcPort
    | Sample_DstMac
    | Sample_DstAsn
    | Sample_DstNextAsn
    | Sample_DstPrefix
    | Sample_DstIp
    | Sample_DstPort


{-| Derived from MAU MIB (RFC 2668)
-}
type State_Duplex
    = State_UnknownDuplex
    | State_FullDuplex
    | State_HalfDuplex
    | State_InDuplex
    | State_OutDuplex


type State_Oper
    = State_NotUp
    | State_Up


aliasState_Down : State_Oper
aliasState_Down =
    State_NotUp


type alias Agent =
    { agent : String
    , name : String -- Blank = unknown
    , oob : String
    , role : Agent_Role
    , slot : Agent_Slot
    , boot : Agent_Boot
    , disk : String
    }


type alias Interface =
    { agent : String
    , ifIndex : Int
    , name : String
    , state : State
    , rates : Rates
    }



-- Samples


type alias KnownTags =
    { agents : List String
    , inputDiscards : List Int
    , outputDiscards : List Int
    , l3Protocols : List Int
    , l4Protocols : List Int
    }



-- Request / responses


type alias ListAgentsRequest =
    {}


type alias ListAgentsResponse =
    { agents : List Agent
    , interfaces : List Interface
    }


type alias ListKnownTagsRequest =
    {}


type alias ListKnownTagsResponse =
    { results : KnownTags
    }


type alias ListRatesRequest =
    { window : Window
    , role : Agent_Role
    , agent : Maybe String
    , ifIndex : Maybe Int
    }


type alias ListRatesResponse =
    { results : List Rates
    }


type alias ListSamplesRequest =
    { window : Window
    , filter : TagFilter
    , top : Int
    , groups : List Sample_Group
    }


type alias ListSamplesResponse =
    { results : Series
    }


type alias Rates =
    { interval : Time.Posix
    , agent : String
    , ifIndex : Int
    , inOctets : Protobuf.Elmer.DoubleValue
    , inUnicast : Protobuf.Elmer.DoubleValue
    , inMulticast : Protobuf.Elmer.DoubleValue
    , inBroadcast : Protobuf.Elmer.DoubleValue
    , inDiscards : Protobuf.Elmer.DoubleValue
    , inErrors : Protobuf.Elmer.DoubleValue
    , inUnknownProtos : Protobuf.Elmer.DoubleValue
    , outOctets : Protobuf.Elmer.DoubleValue
    , outUnicast : Protobuf.Elmer.DoubleValue
    , outMulticast : Protobuf.Elmer.DoubleValue
    , outBroadcast : Protobuf.Elmer.DoubleValue
    , outDiscards : Protobuf.Elmer.DoubleValue
    , outErrors : Protobuf.Elmer.DoubleValue
    }


type alias Sample =
    { packets : Float
    , bytes : Float
    }


type alias SampleTag =
    { xtype : SampleTag_Type
    , filter : TagFilter
    }


type alias Samples =
    { interval : Time.Posix
    , samples : Dict Int Sample -- Key is tag ID
    }


type alias Series =
    { tags : Dict Int SampleTag
    , series : List Samples
    }


type alias State =
    { id : String
    , agent : String
    , ifIndex : Int
    , received : Time.Posix
    , xtype : Int
    , speed : Float -- Mbit/s
    , direction : State_Duplex
    , promiscuous : Bool
    , admin : State_Oper
    , oper : State_Oper
    }


type alias TagFilter =
    { role : Agent_Role
    , agent : Maybe String
    , ifIndex : Maybe Int
    , input : Maybe Int
    , output : Maybe Int
    , inputDiscard : Maybe Int
    , outputDiscard : Maybe Int
    , l3Protocol : Maybe Int
    , l4Protocol : Maybe Int
    , srcMac : Maybe String
    , srcPrefix : Maybe String
    , srcIp : Maybe String
    , srcAsn : Maybe Int
    , srcNextAsn : Maybe Int
    , srcPort : Maybe Int
    , dstMac : Maybe String
    , dstPrefix : Maybe String
    , dstIp : Maybe String
    , dstAsn : Maybe Int
    , dstNextAsn : Maybe Int
    , dstPort : Maybe Int
    }



-- Generic counters


type alias Window =
    { before : Maybe Time.Posix
    , interval : Maybe Google.Protobuf.Duration
    , limit : Int
    }


emptyAgent : Agent
emptyAgent =
    Agent "" "" "" emptyAgent_Role emptyAgent_Slot emptyAgent_Boot ""


emptyInterface : Interface
emptyInterface =
    Interface "" 0 "" emptyState emptyRates


emptyKnownTags : KnownTags
emptyKnownTags =
    KnownTags [] [] [] [] []


emptyListAgentsRequest : ListAgentsRequest
emptyListAgentsRequest =
    ListAgentsRequest


emptyListAgentsResponse : ListAgentsResponse
emptyListAgentsResponse =
    ListAgentsResponse [] []


emptyListKnownTagsRequest : ListKnownTagsRequest
emptyListKnownTagsRequest =
    ListKnownTagsRequest


emptyListKnownTagsResponse : ListKnownTagsResponse
emptyListKnownTagsResponse =
    ListKnownTagsResponse emptyKnownTags


emptyListRatesRequest : ListRatesRequest
emptyListRatesRequest =
    ListRatesRequest emptyWindow emptyAgent_Role Nothing Nothing


emptyListRatesResponse : ListRatesResponse
emptyListRatesResponse =
    ListRatesResponse []


emptyListSamplesRequest : ListSamplesRequest
emptyListSamplesRequest =
    ListSamplesRequest emptyWindow emptyTagFilter 0 []


emptyListSamplesResponse : ListSamplesResponse
emptyListSamplesResponse =
    ListSamplesResponse emptySeries


emptyRates : Rates
emptyRates =
    Rates Protobuf.Elmer.emptyTimestamp "" 0 Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue Protobuf.Elmer.emptyDoubleValue


emptySample : Sample
emptySample =
    Sample 0 0


emptySampleTag : SampleTag
emptySampleTag =
    SampleTag emptySampleTag_Type emptyTagFilter


emptySamples : Samples
emptySamples =
    Samples Protobuf.Elmer.emptyTimestamp Dict.empty


emptySeries : Series
emptySeries =
    Series Dict.empty []


emptyState : State
emptyState =
    State "" "" 0 Protobuf.Elmer.emptyTimestamp 0 0 emptyState_Duplex False emptyState_Oper emptyState_Oper


emptyTagFilter : TagFilter
emptyTagFilter =
    TagFilter emptyAgent_Role Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing Nothing


emptyWindow : Window
emptyWindow =
    Window Nothing Nothing 0


emptyAgent_Boot : Agent_Boot
emptyAgent_Boot =
    Agent_Disk


emptyAgent_Role : Agent_Role
emptyAgent_Role =
    Agent_Opaque


emptyAgent_Slot : Agent_Slot
emptyAgent_Slot =
    Agent_NotSlot


emptySampleTag_Type : SampleTag_Type
emptySampleTag_Type =
    SampleTag_Unkown


emptySample_Group : Sample_Group
emptySample_Group =
    Sample_NoGroup


emptyState_Duplex : State_Duplex
emptyState_Duplex =
    State_UnknownDuplex


emptyState_Oper : State_Oper
emptyState_Oper =
    State_NotUp


valuesOfAgent_Boot : List Agent_Boot
valuesOfAgent_Boot =
    [ Agent_Disk, Agent_Pxe ]


fromAgent_Boot : Agent_Boot -> String
fromAgent_Boot u =
    case u of
        Agent_Disk ->
            "DISK"

        Agent_Pxe ->
            "PXE"


toAgent_Boot : String -> Agent_Boot
toAgent_Boot str =
    case str of
        "DISK" ->
            Agent_Disk

        "PXE" ->
            Agent_Pxe

        _ ->
            Agent_Disk


valuesOfAgent_Role : List Agent_Role
valuesOfAgent_Role =
    [ Agent_Opaque, Agent_Router, Agent_Server, Agent_Oob ]


fromAgent_Role : Agent_Role -> String
fromAgent_Role u =
    case u of
        Agent_Opaque ->
            "OPAQUE"

        Agent_Router ->
            "ROUTER"

        Agent_Server ->
            "SERVER"

        Agent_Oob ->
            "OOB"


toAgent_Role : String -> Agent_Role
toAgent_Role str =
    case str of
        "OPAQUE" ->
            Agent_Opaque

        "ROUTER" ->
            Agent_Router

        "SERVER" ->
            Agent_Server

        "OOB" ->
            Agent_Oob

        _ ->
            Agent_Opaque


valuesOfAgent_Slot : List Agent_Slot
valuesOfAgent_Slot =
    [ Agent_NotSlot, Agent_Capacity, Agent_Capability ]


fromAgent_Slot : Agent_Slot -> String
fromAgent_Slot u =
    case u of
        Agent_NotSlot ->
            "NOT_SLOT"

        Agent_Capacity ->
            "CAPACITY"

        Agent_Capability ->
            "CAPABILITY"


toAgent_Slot : String -> Agent_Slot
toAgent_Slot str =
    case str of
        "NOT_SLOT" ->
            Agent_NotSlot

        "CAPACITY" ->
            Agent_Capacity

        "CAPABILITY" ->
            Agent_Capability

        _ ->
            Agent_NotSlot


valuesOfSampleTag_Type : List SampleTag_Type
valuesOfSampleTag_Type =
    [ SampleTag_Unkown, SampleTag_Filter, SampleTag_Top ]


fromSampleTag_Type : SampleTag_Type -> String
fromSampleTag_Type u =
    case u of
        SampleTag_Unkown ->
            "UNKOWN"

        SampleTag_Filter ->
            "FILTER"

        SampleTag_Top ->
            "TOP"


toSampleTag_Type : String -> SampleTag_Type
toSampleTag_Type str =
    case str of
        "UNKOWN" ->
            SampleTag_Unkown

        "FILTER" ->
            SampleTag_Filter

        "TOP" ->
            SampleTag_Top

        _ ->
            SampleTag_Unkown


valuesOfSample_Group : List Sample_Group
valuesOfSample_Group =
    [ Sample_NoGroup, Sample_Role, Sample_Agent, Sample_Interface, Sample_Input, Sample_Output, Sample_InputDiscard, Sample_OutputDiscard, Sample_L3Protocol, Sample_L4Protocol, Sample_SrcMac, Sample_SrcAsn, Sample_SrcNextAsn, Sample_SrcPrefix, Sample_SrcIp, Sample_SrcPort, Sample_DstMac, Sample_DstAsn, Sample_DstNextAsn, Sample_DstPrefix, Sample_DstIp, Sample_DstPort ]


fromSample_Group : Sample_Group -> String
fromSample_Group u =
    case u of
        Sample_NoGroup ->
            "NO_GROUP"

        Sample_Role ->
            "ROLE"

        Sample_Agent ->
            "AGENT"

        Sample_Interface ->
            "INTERFACE"

        Sample_Input ->
            "INPUT"

        Sample_Output ->
            "OUTPUT"

        Sample_InputDiscard ->
            "INPUT_DISCARD"

        Sample_OutputDiscard ->
            "OUTPUT_DISCARD"

        Sample_L3Protocol ->
            "L3_PROTOCOL"

        Sample_L4Protocol ->
            "L4_PROTOCOL"

        Sample_SrcMac ->
            "SRC_MAC"

        Sample_SrcAsn ->
            "SRC_ASN"

        Sample_SrcNextAsn ->
            "SRC_NEXT_ASN"

        Sample_SrcPrefix ->
            "SRC_PREFIX"

        Sample_SrcIp ->
            "SRC_IP"

        Sample_SrcPort ->
            "SRC_PORT"

        Sample_DstMac ->
            "DST_MAC"

        Sample_DstAsn ->
            "DST_ASN"

        Sample_DstNextAsn ->
            "DST_NEXT_ASN"

        Sample_DstPrefix ->
            "DST_PREFIX"

        Sample_DstIp ->
            "DST_IP"

        Sample_DstPort ->
            "DST_PORT"


toSample_Group : String -> Sample_Group
toSample_Group str =
    case str of
        "NO_GROUP" ->
            Sample_NoGroup

        "ROLE" ->
            Sample_Role

        "AGENT" ->
            Sample_Agent

        "INTERFACE" ->
            Sample_Interface

        "INPUT" ->
            Sample_Input

        "OUTPUT" ->
            Sample_Output

        "INPUT_DISCARD" ->
            Sample_InputDiscard

        "OUTPUT_DISCARD" ->
            Sample_OutputDiscard

        "L3_PROTOCOL" ->
            Sample_L3Protocol

        "L4_PROTOCOL" ->
            Sample_L4Protocol

        "SRC_MAC" ->
            Sample_SrcMac

        "SRC_ASN" ->
            Sample_SrcAsn

        "SRC_NEXT_ASN" ->
            Sample_SrcNextAsn

        "SRC_PREFIX" ->
            Sample_SrcPrefix

        "SRC_IP" ->
            Sample_SrcIp

        "SRC_PORT" ->
            Sample_SrcPort

        "DST_MAC" ->
            Sample_DstMac

        "DST_ASN" ->
            Sample_DstAsn

        "DST_NEXT_ASN" ->
            Sample_DstNextAsn

        "DST_PREFIX" ->
            Sample_DstPrefix

        "DST_IP" ->
            Sample_DstIp

        "DST_PORT" ->
            Sample_DstPort

        _ ->
            Sample_NoGroup


valuesOfState_Duplex : List State_Duplex
valuesOfState_Duplex =
    [ State_UnknownDuplex, State_FullDuplex, State_HalfDuplex, State_InDuplex, State_OutDuplex ]


fromState_Duplex : State_Duplex -> String
fromState_Duplex u =
    case u of
        State_UnknownDuplex ->
            "UNKNOWN_DUPLEX"

        State_FullDuplex ->
            "FULL_DUPLEX"

        State_HalfDuplex ->
            "HALF_DUPLEX"

        State_InDuplex ->
            "IN_DUPLEX"

        State_OutDuplex ->
            "OUT_DUPLEX"


toState_Duplex : String -> State_Duplex
toState_Duplex str =
    case str of
        "UNKNOWN_DUPLEX" ->
            State_UnknownDuplex

        "FULL_DUPLEX" ->
            State_FullDuplex

        "HALF_DUPLEX" ->
            State_HalfDuplex

        "IN_DUPLEX" ->
            State_InDuplex

        "OUT_DUPLEX" ->
            State_OutDuplex

        _ ->
            State_UnknownDuplex


valuesOfState_Oper : List State_Oper
valuesOfState_Oper =
    [ State_NotUp, State_Up ]


fromState_Oper : State_Oper -> String
fromState_Oper u =
    case u of
        State_NotUp ->
            "NOT_UP"

        State_Up ->
            "UP"


toState_Oper : String -> State_Oper
toState_Oper str =
    case str of
        "NOT_UP" ->
            State_NotUp

        "UP" ->
            State_Up

        _ ->
            State_NotUp


decodeAgent : PD.Decoder Agent
decodeAgent =
    PD.message emptyAgent
        [ PD.optional 1 PD.string (\v m -> { m | agent = v })
        , PD.optional 2 PD.string (\v m -> { m | name = v })
        , PD.optional 3 PD.string (\v m -> { m | oob = v })
        , PD.optional 4 decodeAgent_Role (\v m -> { m | role = v })
        , PD.optional 5 decodeAgent_Slot (\v m -> { m | slot = v })
        , PD.optional 6 decodeAgent_Boot (\v m -> { m | boot = v })
        , PD.optional 7 PD.string (\v m -> { m | disk = v })
        ]


decodeInterface : PD.Decoder Interface
decodeInterface =
    PD.message emptyInterface
        [ PD.optional 1 PD.string (\v m -> { m | agent = v })
        , PD.optional 2 PD.uint32 (\v m -> { m | ifIndex = v })
        , PD.optional 3 PD.string (\v m -> { m | name = v })
        , PD.optional 4 decodeState (\v m -> { m | state = v })
        , PD.optional 5 decodeRates (\v m -> { m | rates = v })
        ]


decodeKnownTags : PD.Decoder KnownTags
decodeKnownTags =
    PD.message emptyKnownTags
        [ PD.repeated 1 PD.string .agents (\v m -> { m | agents = v })
        , PD.repeated 2 PD.uint32 .inputDiscards (\v m -> { m | inputDiscards = v })
        , PD.repeated 3 PD.uint32 .outputDiscards (\v m -> { m | outputDiscards = v })
        , PD.repeated 4 PD.uint32 .l3Protocols (\v m -> { m | l3Protocols = v })
        , PD.repeated 5 PD.uint32 .l4Protocols (\v m -> { m | l4Protocols = v })
        ]


decodeListAgentsRequest : PD.Decoder ListAgentsRequest
decodeListAgentsRequest =
    PD.message emptyListAgentsRequest
        []


decodeListAgentsResponse : PD.Decoder ListAgentsResponse
decodeListAgentsResponse =
    PD.message emptyListAgentsResponse
        [ PD.repeated 1 decodeAgent .agents (\v m -> { m | agents = v })
        , PD.repeated 2 decodeInterface .interfaces (\v m -> { m | interfaces = v })
        ]


decodeListKnownTagsRequest : PD.Decoder ListKnownTagsRequest
decodeListKnownTagsRequest =
    PD.message emptyListKnownTagsRequest
        []


decodeListKnownTagsResponse : PD.Decoder ListKnownTagsResponse
decodeListKnownTagsResponse =
    PD.message emptyListKnownTagsResponse
        [ PD.optional 1 decodeKnownTags (\v m -> { m | results = v })
        ]


decodeListRatesRequest : PD.Decoder ListRatesRequest
decodeListRatesRequest =
    let
        decodeListRatesRequest_Agent =
            [ ( 2, PD.string )
            ]

        decodeListRatesRequest_IfIndex =
            [ ( 3, PD.uint32 )
            ]
    in
    PD.message emptyListRatesRequest
        [ PD.optional 1 decodeWindow (\v m -> { m | window = v })
        , PD.optional 4 decodeAgent_Role (\v m -> { m | role = v })
        , PD.oneOf decodeListRatesRequest_Agent (\v m -> { m | agent = v })
        , PD.oneOf decodeListRatesRequest_IfIndex (\v m -> { m | ifIndex = v })
        ]


decodeListRatesResponse : PD.Decoder ListRatesResponse
decodeListRatesResponse =
    PD.message emptyListRatesResponse
        [ PD.repeated 1 decodeRates .results (\v m -> { m | results = v })
        ]


decodeListSamplesRequest : PD.Decoder ListSamplesRequest
decodeListSamplesRequest =
    PD.message emptyListSamplesRequest
        [ PD.optional 1 decodeWindow (\v m -> { m | window = v })
        , PD.optional 2 decodeTagFilter (\v m -> { m | filter = v })
        , PD.optional 3 PD.uint32 (\v m -> { m | top = v })
        , PD.repeated 4 decodeSample_Group .groups (\v m -> { m | groups = v })
        ]


decodeListSamplesResponse : PD.Decoder ListSamplesResponse
decodeListSamplesResponse =
    PD.message emptyListSamplesResponse
        [ PD.optional 1 decodeSeries (\v m -> { m | results = v })
        ]


decodeRates : PD.Decoder Rates
decodeRates =
    PD.message emptyRates
        [ PD.optional 1 Protobuf.Elmer.decodeTimestamp (\v m -> { m | interval = v })
        , PD.optional 2 PD.string (\v m -> { m | agent = v })
        , PD.optional 3 PD.uint32 (\v m -> { m | ifIndex = v })
        , PD.optional 4 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inOctets = v })
        , PD.optional 5 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inUnicast = v })
        , PD.optional 6 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inMulticast = v })
        , PD.optional 7 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inBroadcast = v })
        , PD.optional 8 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inDiscards = v })
        , PD.optional 9 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inErrors = v })
        , PD.optional 10 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | inUnknownProtos = v })
        , PD.optional 11 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outOctets = v })
        , PD.optional 12 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outUnicast = v })
        , PD.optional 13 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outMulticast = v })
        , PD.optional 14 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outBroadcast = v })
        , PD.optional 15 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outDiscards = v })
        , PD.optional 16 Protobuf.Elmer.decodeDoubleValue (\v m -> { m | outErrors = v })
        ]


decodeSample : PD.Decoder Sample
decodeSample =
    PD.message emptySample
        [ PD.optional 1 PD.double (\v m -> { m | packets = v })
        , PD.optional 2 PD.double (\v m -> { m | bytes = v })
        ]


decodeSampleTag : PD.Decoder SampleTag
decodeSampleTag =
    PD.message emptySampleTag
        [ PD.optional 1 decodeSampleTag_Type (\v m -> { m | xtype = v })
        , PD.optional 2 decodeTagFilter (\v m -> { m | filter = v })
        ]


decodeSamples : PD.Decoder Samples
decodeSamples =
    PD.message emptySamples
        [ PD.optional 1 Protobuf.Elmer.decodeTimestamp (\v m -> { m | interval = v })
        , PD.mapped 2 ( 0, emptySample ) PD.int32 decodeSample .samples (\v m -> { m | samples = v })
        ]


decodeSeries : PD.Decoder Series
decodeSeries =
    PD.message emptySeries
        [ PD.mapped 1 ( 0, emptySampleTag ) PD.int32 decodeSampleTag .tags (\v m -> { m | tags = v })
        , PD.repeated 2 decodeSamples .series (\v m -> { m | series = v })
        ]


decodeState : PD.Decoder State
decodeState =
    PD.message emptyState
        [ PD.optional 1 PD.string (\v m -> { m | id = v })
        , PD.optional 2 PD.string (\v m -> { m | agent = v })
        , PD.optional 3 PD.uint32 (\v m -> { m | ifIndex = v })
        , PD.optional 4 Protobuf.Elmer.decodeTimestamp (\v m -> { m | received = v })
        , PD.optional 5 PD.uint32 (\v m -> { m | xtype = v })
        , PD.optional 6 PD.double (\v m -> { m | speed = v })
        , PD.optional 7 decodeState_Duplex (\v m -> { m | direction = v })
        , PD.optional 8 PD.bool (\v m -> { m | promiscuous = v })
        , PD.optional 9 decodeState_Oper (\v m -> { m | admin = v })
        , PD.optional 10 decodeState_Oper (\v m -> { m | oper = v })
        ]


decodeTagFilter : PD.Decoder TagFilter
decodeTagFilter =
    let
        decodeTagFilter_Agent =
            [ ( 2, PD.string )
            ]

        decodeTagFilter_IfIndex =
            [ ( 3, PD.uint32 )
            ]

        decodeTagFilter_Input =
            [ ( 4, PD.uint32 )
            ]

        decodeTagFilter_Output =
            [ ( 5, PD.uint32 )
            ]

        decodeTagFilter_InputDiscard =
            [ ( 6, PD.uint32 )
            ]

        decodeTagFilter_OutputDiscard =
            [ ( 7, PD.uint32 )
            ]

        decodeTagFilter_L3Protocol =
            [ ( 8, PD.uint32 )
            ]

        decodeTagFilter_L4Protocol =
            [ ( 9, PD.uint32 )
            ]

        decodeTagFilter_SrcMac =
            [ ( 10, PD.string )
            ]

        decodeTagFilter_SrcPrefix =
            [ ( 11, PD.string )
            ]

        decodeTagFilter_SrcIp =
            [ ( 12, PD.string )
            ]

        decodeTagFilter_SrcAsn =
            [ ( 13, PD.uint32 )
            ]

        decodeTagFilter_SrcNextAsn =
            [ ( 14, PD.uint32 )
            ]

        decodeTagFilter_SrcPort =
            [ ( 15, PD.uint32 )
            ]

        decodeTagFilter_DstMac =
            [ ( 16, PD.string )
            ]

        decodeTagFilter_DstPrefix =
            [ ( 17, PD.string )
            ]

        decodeTagFilter_DstIp =
            [ ( 18, PD.string )
            ]

        decodeTagFilter_DstAsn =
            [ ( 19, PD.uint32 )
            ]

        decodeTagFilter_DstNextAsn =
            [ ( 20, PD.uint32 )
            ]

        decodeTagFilter_DstPort =
            [ ( 21, PD.uint32 )
            ]
    in
    PD.message emptyTagFilter
        [ PD.optional 1 decodeAgent_Role (\v m -> { m | role = v })
        , PD.oneOf decodeTagFilter_Agent (\v m -> { m | agent = v })
        , PD.oneOf decodeTagFilter_IfIndex (\v m -> { m | ifIndex = v })
        , PD.oneOf decodeTagFilter_Input (\v m -> { m | input = v })
        , PD.oneOf decodeTagFilter_Output (\v m -> { m | output = v })
        , PD.oneOf decodeTagFilter_InputDiscard (\v m -> { m | inputDiscard = v })
        , PD.oneOf decodeTagFilter_OutputDiscard (\v m -> { m | outputDiscard = v })
        , PD.oneOf decodeTagFilter_L3Protocol (\v m -> { m | l3Protocol = v })
        , PD.oneOf decodeTagFilter_L4Protocol (\v m -> { m | l4Protocol = v })
        , PD.oneOf decodeTagFilter_SrcMac (\v m -> { m | srcMac = v })
        , PD.oneOf decodeTagFilter_SrcPrefix (\v m -> { m | srcPrefix = v })
        , PD.oneOf decodeTagFilter_SrcIp (\v m -> { m | srcIp = v })
        , PD.oneOf decodeTagFilter_SrcAsn (\v m -> { m | srcAsn = v })
        , PD.oneOf decodeTagFilter_SrcNextAsn (\v m -> { m | srcNextAsn = v })
        , PD.oneOf decodeTagFilter_SrcPort (\v m -> { m | srcPort = v })
        , PD.oneOf decodeTagFilter_DstMac (\v m -> { m | dstMac = v })
        , PD.oneOf decodeTagFilter_DstPrefix (\v m -> { m | dstPrefix = v })
        , PD.oneOf decodeTagFilter_DstIp (\v m -> { m | dstIp = v })
        , PD.oneOf decodeTagFilter_DstAsn (\v m -> { m | dstAsn = v })
        , PD.oneOf decodeTagFilter_DstNextAsn (\v m -> { m | dstNextAsn = v })
        , PD.oneOf decodeTagFilter_DstPort (\v m -> { m | dstPort = v })
        ]


decodeWindow : PD.Decoder Window
decodeWindow =
    let
        decodeWindow_Before =
            [ ( 1, Protobuf.Elmer.decodeTimestamp )
            ]

        decodeWindow_Interval =
            [ ( 2, Google.Protobuf.durationDecoder )
            ]
    in
    PD.message emptyWindow
        [ PD.oneOf decodeWindow_Before (\v m -> { m | before = v })
        , PD.oneOf decodeWindow_Interval (\v m -> { m | interval = v })
        , PD.optional 3 PD.uint32 (\v m -> { m | limit = v })
        ]


decodeAgent_Boot : PD.Decoder Agent_Boot
decodeAgent_Boot =
    let
        conv v =
            case v of
                0 ->
                    Agent_Disk

                1 ->
                    Agent_Pxe

                _ ->
                    Agent_Disk
    in
    PD.map conv PD.int32


decodeAgent_Role : PD.Decoder Agent_Role
decodeAgent_Role =
    let
        conv v =
            case v of
                0 ->
                    Agent_Opaque

                1 ->
                    Agent_Router

                2 ->
                    Agent_Server

                3 ->
                    Agent_Oob

                _ ->
                    Agent_Opaque
    in
    PD.map conv PD.int32


decodeAgent_Slot : PD.Decoder Agent_Slot
decodeAgent_Slot =
    let
        conv v =
            case v of
                0 ->
                    Agent_NotSlot

                1 ->
                    Agent_Capacity

                2 ->
                    Agent_Capability

                _ ->
                    Agent_NotSlot
    in
    PD.map conv PD.int32


decodeSampleTag_Type : PD.Decoder SampleTag_Type
decodeSampleTag_Type =
    let
        conv v =
            case v of
                0 ->
                    SampleTag_Unkown

                1 ->
                    SampleTag_Filter

                2 ->
                    SampleTag_Top

                _ ->
                    SampleTag_Unkown
    in
    PD.map conv PD.int32


decodeSample_Group : PD.Decoder Sample_Group
decodeSample_Group =
    let
        conv v =
            case v of
                0 ->
                    Sample_NoGroup

                1 ->
                    Sample_Role

                2 ->
                    Sample_Agent

                3 ->
                    Sample_Interface

                4 ->
                    Sample_Input

                5 ->
                    Sample_Output

                6 ->
                    Sample_InputDiscard

                7 ->
                    Sample_OutputDiscard

                8 ->
                    Sample_L3Protocol

                9 ->
                    Sample_L4Protocol

                10 ->
                    Sample_SrcMac

                11 ->
                    Sample_SrcAsn

                12 ->
                    Sample_SrcNextAsn

                13 ->
                    Sample_SrcPrefix

                14 ->
                    Sample_SrcIp

                15 ->
                    Sample_SrcPort

                16 ->
                    Sample_DstMac

                17 ->
                    Sample_DstAsn

                18 ->
                    Sample_DstNextAsn

                19 ->
                    Sample_DstPrefix

                20 ->
                    Sample_DstIp

                21 ->
                    Sample_DstPort

                _ ->
                    Sample_NoGroup
    in
    PD.map conv PD.int32


decodeState_Duplex : PD.Decoder State_Duplex
decodeState_Duplex =
    let
        conv v =
            case v of
                0 ->
                    State_UnknownDuplex

                1 ->
                    State_FullDuplex

                2 ->
                    State_HalfDuplex

                3 ->
                    State_InDuplex

                4 ->
                    State_OutDuplex

                _ ->
                    State_UnknownDuplex
    in
    PD.map conv PD.int32


decodeState_Oper : PD.Decoder State_Oper
decodeState_Oper =
    let
        conv v =
            case v of
                0 ->
                    State_NotUp

                1 ->
                    State_Up

                _ ->
                    State_NotUp
    in
    PD.map conv PD.int32


encodeAgent : Agent -> PE.Encoder
encodeAgent v =
    PE.message <|
        [ ( 1, PE.string v.agent )
        , ( 2, PE.string v.name )
        , ( 3, PE.string v.oob )
        , ( 4, encodeAgent_Role v.role )
        , ( 5, encodeAgent_Slot v.slot )
        , ( 6, encodeAgent_Boot v.boot )
        , ( 7, PE.string v.disk )
        ]


encodeInterface : Interface -> PE.Encoder
encodeInterface v =
    PE.message <|
        [ ( 1, PE.string v.agent )
        , ( 2, PE.uint32 v.ifIndex )
        , ( 3, PE.string v.name )
        , ( 4, encodeState v.state )
        , ( 5, encodeRates v.rates )
        ]


encodeKnownTags : KnownTags -> PE.Encoder
encodeKnownTags v =
    PE.message <|
        [ ( 1, PE.list PE.string v.agents )
        , ( 2, PE.list PE.uint32 v.inputDiscards )
        , ( 3, PE.list PE.uint32 v.outputDiscards )
        , ( 4, PE.list PE.uint32 v.l3Protocols )
        , ( 5, PE.list PE.uint32 v.l4Protocols )
        ]


encodeListAgentsRequest : ListAgentsRequest -> PE.Encoder
encodeListAgentsRequest _ =
    PE.message <|
        []


encodeListAgentsResponse : ListAgentsResponse -> PE.Encoder
encodeListAgentsResponse v =
    PE.message <|
        [ ( 1, PE.list encodeAgent v.agents )
        , ( 2, PE.list encodeInterface v.interfaces )
        ]


encodeListKnownTagsRequest : ListKnownTagsRequest -> PE.Encoder
encodeListKnownTagsRequest _ =
    PE.message <|
        []


encodeListKnownTagsResponse : ListKnownTagsResponse -> PE.Encoder
encodeListKnownTagsResponse v =
    PE.message <|
        [ ( 1, encodeKnownTags v.results )
        ]


encodeListRatesRequest : ListRatesRequest -> PE.Encoder
encodeListRatesRequest v =
    let
        encodeListRatesRequest_Agent o =
            case o of
                Just data ->
                    [ ( 2, PE.string data ) ]

                Nothing ->
                    []

        encodeListRatesRequest_IfIndex o =
            case o of
                Just data ->
                    [ ( 3, PE.uint32 data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 1, encodeWindow v.window )
        , ( 4, encodeAgent_Role v.role )
        ]
            ++ encodeListRatesRequest_Agent v.agent
            ++ encodeListRatesRequest_IfIndex v.ifIndex


encodeListRatesResponse : ListRatesResponse -> PE.Encoder
encodeListRatesResponse v =
    PE.message <|
        [ ( 1, PE.list encodeRates v.results )
        ]


encodeListSamplesRequest : ListSamplesRequest -> PE.Encoder
encodeListSamplesRequest v =
    PE.message <|
        [ ( 1, encodeWindow v.window )
        , ( 2, encodeTagFilter v.filter )
        , ( 3, PE.uint32 v.top )
        , ( 4, PE.list encodeSample_Group v.groups )
        ]


encodeListSamplesResponse : ListSamplesResponse -> PE.Encoder
encodeListSamplesResponse v =
    PE.message <|
        [ ( 1, encodeSeries v.results )
        ]


encodeRates : Rates -> PE.Encoder
encodeRates v =
    PE.message <|
        [ ( 1, Protobuf.Elmer.encodeTimestamp v.interval )
        , ( 2, PE.string v.agent )
        , ( 3, PE.uint32 v.ifIndex )
        , ( 4, Protobuf.Elmer.encodeDoubleValue v.inOctets )
        , ( 5, Protobuf.Elmer.encodeDoubleValue v.inUnicast )
        , ( 6, Protobuf.Elmer.encodeDoubleValue v.inMulticast )
        , ( 7, Protobuf.Elmer.encodeDoubleValue v.inBroadcast )
        , ( 8, Protobuf.Elmer.encodeDoubleValue v.inDiscards )
        , ( 9, Protobuf.Elmer.encodeDoubleValue v.inErrors )
        , ( 10, Protobuf.Elmer.encodeDoubleValue v.inUnknownProtos )
        , ( 11, Protobuf.Elmer.encodeDoubleValue v.outOctets )
        , ( 12, Protobuf.Elmer.encodeDoubleValue v.outUnicast )
        , ( 13, Protobuf.Elmer.encodeDoubleValue v.outMulticast )
        , ( 14, Protobuf.Elmer.encodeDoubleValue v.outBroadcast )
        , ( 15, Protobuf.Elmer.encodeDoubleValue v.outDiscards )
        , ( 16, Protobuf.Elmer.encodeDoubleValue v.outErrors )
        ]


encodeSample : Sample -> PE.Encoder
encodeSample v =
    PE.message <|
        [ ( 1, PE.double v.packets )
        , ( 2, PE.double v.bytes )
        ]


encodeSampleTag : SampleTag -> PE.Encoder
encodeSampleTag v =
    PE.message <|
        [ ( 1, encodeSampleTag_Type v.xtype )
        , ( 2, encodeTagFilter v.filter )
        ]


encodeSamples : Samples -> PE.Encoder
encodeSamples v =
    PE.message <|
        [ ( 1, Protobuf.Elmer.encodeTimestamp v.interval )
        , ( 2, PE.dict PE.int32 encodeSample v.samples )
        ]


encodeSeries : Series -> PE.Encoder
encodeSeries v =
    PE.message <|
        [ ( 1, PE.dict PE.int32 encodeSampleTag v.tags )
        , ( 2, PE.list encodeSamples v.series )
        ]


encodeState : State -> PE.Encoder
encodeState v =
    PE.message <|
        [ ( 1, PE.string v.id )
        , ( 2, PE.string v.agent )
        , ( 3, PE.uint32 v.ifIndex )
        , ( 4, Protobuf.Elmer.encodeTimestamp v.received )
        , ( 5, PE.uint32 v.xtype )
        , ( 6, PE.double v.speed )
        , ( 7, encodeState_Duplex v.direction )
        , ( 8, PE.bool v.promiscuous )
        , ( 9, encodeState_Oper v.admin )
        , ( 10, encodeState_Oper v.oper )
        ]


encodeTagFilter : TagFilter -> PE.Encoder
encodeTagFilter v =
    let
        encodeTagFilter_Agent o =
            case o of
                Just data ->
                    [ ( 2, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_IfIndex o =
            case o of
                Just data ->
                    [ ( 3, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_Input o =
            case o of
                Just data ->
                    [ ( 4, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_Output o =
            case o of
                Just data ->
                    [ ( 5, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_InputDiscard o =
            case o of
                Just data ->
                    [ ( 6, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_OutputDiscard o =
            case o of
                Just data ->
                    [ ( 7, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_L3Protocol o =
            case o of
                Just data ->
                    [ ( 8, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_L4Protocol o =
            case o of
                Just data ->
                    [ ( 9, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcMac o =
            case o of
                Just data ->
                    [ ( 10, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcPrefix o =
            case o of
                Just data ->
                    [ ( 11, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcIp o =
            case o of
                Just data ->
                    [ ( 12, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcAsn o =
            case o of
                Just data ->
                    [ ( 13, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcNextAsn o =
            case o of
                Just data ->
                    [ ( 14, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_SrcPort o =
            case o of
                Just data ->
                    [ ( 15, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstMac o =
            case o of
                Just data ->
                    [ ( 16, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstPrefix o =
            case o of
                Just data ->
                    [ ( 17, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstIp o =
            case o of
                Just data ->
                    [ ( 18, PE.string data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstAsn o =
            case o of
                Just data ->
                    [ ( 19, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstNextAsn o =
            case o of
                Just data ->
                    [ ( 20, PE.uint32 data ) ]

                Nothing ->
                    []

        encodeTagFilter_DstPort o =
            case o of
                Just data ->
                    [ ( 21, PE.uint32 data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 1, encodeAgent_Role v.role )
        ]
            ++ encodeTagFilter_Agent v.agent
            ++ encodeTagFilter_IfIndex v.ifIndex
            ++ encodeTagFilter_Input v.input
            ++ encodeTagFilter_Output v.output
            ++ encodeTagFilter_InputDiscard v.inputDiscard
            ++ encodeTagFilter_OutputDiscard v.outputDiscard
            ++ encodeTagFilter_L3Protocol v.l3Protocol
            ++ encodeTagFilter_L4Protocol v.l4Protocol
            ++ encodeTagFilter_SrcMac v.srcMac
            ++ encodeTagFilter_SrcPrefix v.srcPrefix
            ++ encodeTagFilter_SrcIp v.srcIp
            ++ encodeTagFilter_SrcAsn v.srcAsn
            ++ encodeTagFilter_SrcNextAsn v.srcNextAsn
            ++ encodeTagFilter_SrcPort v.srcPort
            ++ encodeTagFilter_DstMac v.dstMac
            ++ encodeTagFilter_DstPrefix v.dstPrefix
            ++ encodeTagFilter_DstIp v.dstIp
            ++ encodeTagFilter_DstAsn v.dstAsn
            ++ encodeTagFilter_DstNextAsn v.dstNextAsn
            ++ encodeTagFilter_DstPort v.dstPort


encodeWindow : Window -> PE.Encoder
encodeWindow v =
    let
        encodeWindow_Before o =
            case o of
                Just data ->
                    [ ( 1, Protobuf.Elmer.encodeTimestamp data ) ]

                Nothing ->
                    []

        encodeWindow_Interval o =
            case o of
                Just data ->
                    [ ( 2, Google.Protobuf.toDurationEncoder data ) ]

                Nothing ->
                    []
    in
    PE.message <|
        [ ( 3, PE.uint32 v.limit )
        ]
            ++ encodeWindow_Before v.before
            ++ encodeWindow_Interval v.interval


encodeAgent_Boot : Agent_Boot -> PE.Encoder
encodeAgent_Boot v =
    let
        conv =
            case v of
                Agent_Disk ->
                    0

                Agent_Pxe ->
                    1
    in
    PE.int32 conv


encodeAgent_Role : Agent_Role -> PE.Encoder
encodeAgent_Role v =
    let
        conv =
            case v of
                Agent_Opaque ->
                    0

                Agent_Router ->
                    1

                Agent_Server ->
                    2

                Agent_Oob ->
                    3
    in
    PE.int32 conv


encodeAgent_Slot : Agent_Slot -> PE.Encoder
encodeAgent_Slot v =
    let
        conv =
            case v of
                Agent_NotSlot ->
                    0

                Agent_Capacity ->
                    1

                Agent_Capability ->
                    2
    in
    PE.int32 conv


encodeSampleTag_Type : SampleTag_Type -> PE.Encoder
encodeSampleTag_Type v =
    let
        conv =
            case v of
                SampleTag_Unkown ->
                    0

                SampleTag_Filter ->
                    1

                SampleTag_Top ->
                    2
    in
    PE.int32 conv


encodeSample_Group : Sample_Group -> PE.Encoder
encodeSample_Group v =
    let
        conv =
            case v of
                Sample_NoGroup ->
                    0

                Sample_Role ->
                    1

                Sample_Agent ->
                    2

                Sample_Interface ->
                    3

                Sample_Input ->
                    4

                Sample_Output ->
                    5

                Sample_InputDiscard ->
                    6

                Sample_OutputDiscard ->
                    7

                Sample_L3Protocol ->
                    8

                Sample_L4Protocol ->
                    9

                Sample_SrcMac ->
                    10

                Sample_SrcAsn ->
                    11

                Sample_SrcNextAsn ->
                    12

                Sample_SrcPrefix ->
                    13

                Sample_SrcIp ->
                    14

                Sample_SrcPort ->
                    15

                Sample_DstMac ->
                    16

                Sample_DstAsn ->
                    17

                Sample_DstNextAsn ->
                    18

                Sample_DstPrefix ->
                    19

                Sample_DstIp ->
                    20

                Sample_DstPort ->
                    21
    in
    PE.int32 conv


encodeState_Duplex : State_Duplex -> PE.Encoder
encodeState_Duplex v =
    let
        conv =
            case v of
                State_UnknownDuplex ->
                    0

                State_FullDuplex ->
                    1

                State_HalfDuplex ->
                    2

                State_InDuplex ->
                    3

                State_OutDuplex ->
                    4
    in
    PE.int32 conv


encodeState_Oper : State_Oper -> PE.Encoder
encodeState_Oper v =
    let
        conv =
            case v of
                State_NotUp ->
                    0

                State_Up ->
                    1
    in
    PE.int32 conv
