# End to end example

An example showing the full realisation of types without borders using Elm and Go.

All commands are assumed to be run from this directory.

## Server

The server implements Twirp's [hat making example](https://github.com/twitchtv/twirp/tree/main/example) 🎩 It can be started with:
```
go run go-server/main.go
```

Test the server makes hats with a manual request in a new terminal:
```
echo 'inches:12' \
    | protoc --encode gen.haberdasher.Size api.proto \
    | curl -s --request POST \
      --header "Content-Type: application/protobuf" \
      --data-binary @- \
      http://localhost:8080/twirp/gen.haberdasher.Haberdasher/MakeHat \
    | protoc --decode gen.haberdasher.Hat api.proto
```

## Client

It's recommended that you check in your generated code. So it's already prepared under `elm-client/src/Gen`. You can build the client with the standard tools such as `elm reactor`. I recommend using [elm-live](https://www.elm-live.com/) using:
```
cd elm-client
elm-live src/Main.elm -- --debug
```

Then visit [http://localhost:8000] and have fun making hats 🤠

Finally, in the same directory, you can run the generated fuzz tests with `elm-test`:
```
Compiling > Starting tests

elm-test 0.19.1-revision7
-------------------------

Running 4 tests. To reproduce these results, run: elm-test --fuzz 100 --seed 396783544611329


TEST RUN PASSED

Duration: 136 ms
Passed:   4
Failed:   0
```
