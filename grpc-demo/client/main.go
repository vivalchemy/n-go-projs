package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"
	"vivalchemy/grpcdemo/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	// --------------------------------------------------------------------------------
	// Code to connect to the server
	// --------------------------------------------------------------------------------
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewCallTypesClient(conn)

	// --------------------------------------------------------------------------------
	// Code to call services on the server
	// --------------------------------------------------------------------------------
	callUnary(client)
	callServerStreaming(client)
	callClientStreaming(client)
	callBidiStreaming(client)
}

func callUnary(client pb.CallTypesClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.UnaryCall(ctx, &pb.EnumBody{}) // default is random message
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Unary Response: ", res.Message)
}

func callServerStreaming(client pb.CallTypesClient) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// since this is streaming call, we don't need to cancel it
	ctx := context.Background()

	stream, err := client.ServerStreamingCall(ctx, &pb.RepeatedBody{
		Message: []string{"Keep", "pushing", "forward", "and", "never", "stop", "chasing", "your", "dreams", "!"},
	})
	if err != nil {
		log.Fatal(err)
	}

	counter := 0
	for {
		message, err := stream.Recv()
		// stream ended
		if err == io.EOF {
			break
		}
		// some other error
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Server Streaming Response: ", counter, "=> ", message.Message)
		counter++
	}
	log.Println("Server Streaming Response total messages received: ", counter)
}

func callClientStreaming(client pb.CallTypesClient) {
	log.Println("1")
	ctx := context.Background() // no timeout for streaming call
	log.Println("2")
	stream, err := client.ClientStreamingCall(ctx)
	log.Println("3")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("4")

	messages := [10]string{"Believe", "in", "yourself", "and", "anything", "is", "possible", "with", "hard", "work"}
	log.Println("5")
	for _, message := range messages {
		log.Println("7")
		if err := stream.Send(&pb.Body{Message: message}); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client Streaming Response: ", res.Message)
}

func callBidiStreaming(client pb.CallTypesClient) {
	ctx := context.Background()
	stream, err := client.BidiStreamingCall(ctx)
	if err != nil {
		log.Fatal(err)
	}

	messages := [10]string{"Believe", "in", "yourself", "and", "anything", "is", "possible", "with", "hard", "work"}

	wg := sync.WaitGroup{}
	wg.Add(1)
	// receive messages from server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Bidi Streaming Response: ", res.Message)
		}
		wg.Done()
	}()

	for _, message := range messages {
		if err := stream.Send(&pb.Body{Message: message}); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
	stream.CloseSend()

	wg.Wait() // wait till you receive all messages from server

	log.Println("Bidi Streaming Response total messages received: ", len(messages))
}
