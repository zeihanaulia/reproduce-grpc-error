package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pb "github.com/zeihanaulia/reproduce-grpc-error/proto/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 5 * time.Second,
		}))
	if err != nil {
		log.Println(fmt.Errorf("did not connect: %v", err))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	})

	r.Get("/connect", func(w http.ResponseWriter, r *http.Request) {

		// defer conn.Close()
		c := pb.NewGreeterClient(conn)

		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		ctx := context.Background()
		name := "test"
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			// log.Fatalf("could not greet: %v", err)
			log.Println(fmt.Errorf("could not greet: %v", err))
		}
		log.Printf("Greeting: %s", resp.GetMessage())

	})

	log.Printf("http server listening at %s", ":3000")
	if err := http.ListenAndServe(":3001", r); err != nil {
		panic(fmt.Errorf("cannot start server, err: %v", err))
	}
}
