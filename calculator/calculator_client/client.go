package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jisuhan3201/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hi i am client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot get client : %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
	// doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to calculator doUnary RPC...")
	req := &calculatorpb.SumRequest{
		Left:  10,
		Right: 12,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("cannot calculate : %v\n", err)
	}

	log.Printf("Response from Caculator : %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to calculator doServerStreaming RPC...")
	req := &calculatorpb.PrimeDecompositionRequest{
		Prime: 120,
	}
	resStream, err := c.Decompose(context.Background(), req)
	if err != nil {
		log.Fatalf("cannot decompose: %v\n", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while recieving steam response: %v", err)
		}
		log.Printf("Response from stream: %v\n", res)
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to calculator Compute Average Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("")
	}
	numbers := []int32{3, 12, 34, 28}
	for _, num := range numbers {
		fmt.Printf("sending number: %v\n", num)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v\n", err)
	}
	log.Printf("The Average is: %v\n", res.GetAvg())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to calculator FindMaximum Client Streaming RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while FindMaximum: %v\n", err)
	}
	waitc := make(chan struct{})

	// send go routine
	go func() {
		numbers := []int32{4, 7, 2, 9, 3, 15, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number : %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading server stream: %v\n", err)
			}
			maximum := res.GetMaximum()
			fmt.Printf("Receiving new response: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Err Unary RPC...")
	// correct call
	doErrorCall(c, 10)
	// incorrect call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC(user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent negative number")
			} else {
				log.Fatalf("Big error calling SquareRoot : %v\n", err)
			}
		}
	}
	fmt.Printf("Result of square root %v : %v\n", n, res.GetNumberRoot())
}
