package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Samples/BeyondREST/GoGrpcClient/greet"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := greet.NewGreeterClient(conn)
	response, err := client.SayHello(context.Background(), &greet.HelloRequest{Name: "FooBar"})
	if err != nil {
		panic(err)
	}

	log.Println(response.Message)

	fibClient := greet.NewMathGuruClient(conn)
	fibResult, err := fibClient.GetFibonacci(context.Background(), &greet.FromTo{From: 0, To: 50})
	if err != nil {
		panic(err)
	}

	for {
		msg, err := fibResult.Recv()
		if err != nil {
			break;
		}

		fmt.Printf("%d\n", msg.Result)
	}
}
