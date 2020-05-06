package main

import (
	"context"
	"fmt"
	blogpb "github.com/hpierre74/go-grpc-firestore/proto"
	"github.com/hpierre74/go-grpc-firestore/server/api/blog"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	firebase "firebase.google.com/go"
)


func main() {
	ctx := context.Background()
	// Configure 'log' package to give file name and line number on eg. log.Fatal
	// just the filename & line number:
	// log.SetFlags(log.Lshortfile)
	// Or add timestamps and pipe file name and line number to it:
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting server on port :50051...")

	// 50051 is the default port for gRPC
	// Ideally we'd use 0.0.0.0 instead of localhost as well
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// Initialize firestore client
	conf := &firebase.Config{ProjectID: "<YOUR-PROJECT-ID>"}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestore init: %v", err)
	}
	defer client.Close()

	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}
	// var s *grpc.Server
	s := grpc.NewServer(opts...)
	// var srv *blogAPI.BlogServiceServer
	srv := &blog.ServiceServer{DB: &blog.DB{Collection: client.Collection("<YOUR-COLLECTION>")}}

	blogpb.RegisterBlogServiceServer(s, srv)

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	// Bad way to stop the server
	// if err := s.Serve(listener); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }
	// Right way to stop the server using a SHUTDOWN HOOK

	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed
	// And our main routine keeps running
	// If the main routine were to shutdown so would the child routine that is Serving the server
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing Firestore connection")
	client.Close()
	fmt.Println("Done.")
}
