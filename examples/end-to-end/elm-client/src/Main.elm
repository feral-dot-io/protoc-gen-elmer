module Main exposing (main)

import Browser
import Gen.Haberdasher as Haberdasher
import Gen.HaberdasherTwirp as Rpc
import Html as H exposing (Html)
import Html.Attributes as HA
import Html.Events as HE
import Http


api : String
api =
    "http://localhost:8080/twirp"


{-| Standard program that can make HTTP requests
-}
main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = \_ -> Sub.none
        , view = view
        }


type alias Model =
    -- Stores our list of hats or an error if we ever receive one.
    { hats : Result Http.Error (List Haberdasher.Hat)

    -- Form value
    , selectedInches : Int

    -- Are we fetching data? We'd normally use something like RemoteData.Webdata here
    , fetching : Bool
    }


{-| Builds an empty Model
-}
init : flags -> ( Model, Cmd Msg )
init _ =
    ( Model (Ok []) 8 False
    , Cmd.none
    )


type Msg
    = SetSelectedInches String
    | MakeHatRequest
    | HatResult (Result Http.Error Haberdasher.Hat)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SetSelectedInches str ->
            ( { model
                | selectedInches =
                    String.toInt str
                        |> Maybe.withDefault model.selectedInches
                        |> clamp 1 100
              }
            , Cmd.none
            )

        MakeHatRequest ->
            ( model
            , Haberdasher.Size model.selectedInches
                |> Rpc.haberdasher_MakeHat HatResult api
            )

        HatResult result ->
            ( { model | hats = Result.map2 (::) result model.hats }
            , Cmd.none
            )


view : Model -> Html Msg
view model =
    H.div [] <|
        [ H.h1 [] [ H.text "✨ Hats! ✨" ]
        , case model.hats of
            Ok hats ->
                viewArmoire model hats

            Err err ->
                viewHttpError err
        ]


viewArmoire : Model -> List Haberdasher.Hat -> Html Msg
viewArmoire { selectedInches, fetching } hats =
    H.div []
        [ H.p [] [ H.text "Your armoire contains:" ]

        -- List our hats
        , H.ul [] <|
            if List.isEmpty hats then
                [ H.li []
                    [ H.em [] [ H.text "No hats!" ]
                    , H.text " 🤯"
                    ]
                ]

            else
                let
                    liHatter hat =
                        H.li []
                            [ H.text <|
                                hat.name
                                    ++ ", "
                                    ++ hat.color
                                    ++ ", "
                                    ++ String.fromInt hat.size
                                    ++ "″"
                            ]
                in
                List.map liHatter hats

        -- Form
        , H.p [] [ H.text "Would you like a new hat?" ]
        , H.p []
            [ H.label []
                [ H.text "Size (in inches)? "
                , H.input
                    [ HA.type_ "number"
                    , HA.min "1"
                    , HA.max "100"
                    , HA.step "1"
                    , HE.onInput SetSelectedInches
                    , HA.value (String.fromInt selectedInches)
                    ]
                    []
                ]
            ]

        -- Submit
        , H.p []
            [ H.button
                [ HA.disabled fetching
                , HE.onClick MakeHatRequest
                ]
                [ H.text "Why yes, I would love a new hat" ]
            ]
        ]


viewHttpError : Http.Error -> Html Msg
viewHttpError err =
    let
        str =
            case err of
                Http.BadUrl url ->
                    "bad URL: " ++ url

                Http.Timeout ->
                    "timeout"

                Http.NetworkError ->
                    "network problem"

                Http.BadStatus status ->
                    "bad status: " ++ String.fromInt status

                Http.BadBody bodyErr ->
                    "bad body: " ++ bodyErr
    in
    H.p [] [ H.text ("There was an error: " ++ str) ]
