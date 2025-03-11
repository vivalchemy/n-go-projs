package main

import (
	"context"
	"io"
	"math/rand/v2"
	"net"
	"time"
	"vivalchemy/grpcdemo/pb"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

type Server struct {
	pb.UnimplementedCallTypesServer
}

func (srv *Server) UnaryCall(ctx context.Context, req *pb.EnumBody) (*pb.Body, error) {
	// if the user wants a predefined message, return it
	if req.ResponseType == pb.ResponseType_RT_PREDEFINED {
		return &pb.Body{Message: "hello from predefined message"}, nil
	}

	messages := []string{
		"Hello World",
		"The only way to do great work is to love what you do. – Steve Jobs",
		"Success is not final, failure is not fatal: it is the courage to continue that counts. – Winston Churchill",
		"Believe you can and you're halfway there. – Theodore Roosevelt",
		"Your time is limited, so don’t waste it living someone else’s life. – Steve Jobs",
		"Do what you can, with what you have, where you are. – Theodore Roosevelt",
		"It always seems impossible until it’s done. – Nelson Mandela",
		"Hardships often prepare ordinary people for an extraordinary destiny. – C.S. Lewis",
		"Don’t watch the clock; do what it does. Keep going. – Sam Levenson",
		"Doubt kills more dreams than failure ever will. – Suzy Kassem",
		"The best way to predict the future is to create it. – Peter Drucker",
	}

	return &pb.Body{Message: messages[rand.IntN(len(messages))]}, nil
}

func (srv *Server) ServerStreamingCall(req *pb.RepeatedBody, stream grpc.ServerStreamingServer[pb.Body]) error {
	log.Println("Strings Received: ", req.Message)
	for _, message := range req.Message {
		if err := stream.Send(&pb.Body{Message: "hello " + message}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}

func (srv *Server) ClientStreamingCall(stream grpc.ClientStreamingServer[pb.Body, pb.RepeatedBody]) error {
	messages := make([]string, 0)
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.RepeatedBody{Message: messages})
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Strings Received: ", message.Message)
		messages = append(messages, message.Message)
	}
}

func (srv *Server) BidiStreamingCall(stream grpc.BidiStreamingServer[pb.Body, pb.Body]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Strings Received: ", req.Message)
		if err := stream.Send(&pb.Body{Message: "hello" + req.Message}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) // utility to register all the services on grpcui

	pb.RegisterCallTypesServer(grpcServer, &Server{})

	log.Println("Starting server on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
