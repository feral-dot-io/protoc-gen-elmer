package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	pb "github.com/feral-dot-io/protoc-gen-elmer/examples/end-to-end/go-server/gen"
	"github.com/rs/cors"
	"github.com/twitchtv/twirp"
)

func main() {
	impl := &randomHaberdasher{}
	// We're not doing things like auth
	server := pb.NewHaberdasherServer(impl,
		twirp.WithServerHooks(NewLoggingHooks()))

	// Allow CORS (net/http wrapper)
	handler := cors.New(cors.Options{
		AllowOriginFunc:  func(string) bool { return true },
		AllowCredentials: true,
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type"}}).
		Handler(server)

	// Listen for requests
	log.Printf("Listening for RPC requests on http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		log.Fatalf("error listening to RPC server: %s\n", err)
	}
}

// NewLoggingServerHooks logs request and errors to stdout in the service
func NewLoggingHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			pkg, _ := twirp.PackageName(ctx)
			service, _ := twirp.ServiceName(ctx)
			method, _ := twirp.MethodName(ctx)
			log.Printf("Request on `%s.%s.%s`\n", pkg, service, method)
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			log.Printf("Error %s\n", twerr)
			return ctx
		},
	}
}

// Implements our RPC service
type randomHaberdasher struct{}

func (h *randomHaberdasher) MakeHat(ctx context.Context, size *pb.Size) (*pb.Hat, error) {
	if size.Inches <= 0 {
		return nil, twirp.InvalidArgumentError("Inches", "I can't make a hat that small!")
	}
	colors := []string{"white", "black", "brown", "red", "blue"}
	names := []string{"bowler", "baseball cap", "top hat", "derby"}
	return &pb.Hat{
		Size:  size.Inches,
		Color: colors[rand.Intn(len(colors))],
		Name:  names[rand.Intn(len(names))],
	}, nil
}
