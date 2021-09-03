package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pb "github.com/zeihanaulia/reproduce-grpc-error/proto/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const (
	address     = "10.98.96.166:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 5 * time.Second,
		}))
	if err != nil {
		log.Println(fmt.Errorf("did not connect: %v", err))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler)
	r.Get("/readiness", readinessHandler)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		var a = []string{}
		i1 := a[0]

		fmt.Println(i1)
		runtime.Goexit()
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

	server := &http.Server{Addr: "0.0.0.0:3002", Handler: r}

	go func() {
		log.Printf("http server listening at %s", ":3002")
		if err := server.ListenAndServe(); err != nil {
			panic(fmt.Errorf("cannot start server, err: %v", err))
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
