package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/jisuhan3201/go-grpc/calculator/calculatorpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC %v", req)
	result := req.GetLeft() + req.GetRight()
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) Decompose(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_DecomposeServer) error {
	fmt.Printf("Received Deompose RPC %v", req)
	f := int32(2)
	p := req.GetPrime()
	for p > 1 {
		if p%f == 0 {
			stream.Send(&calculatorpb.PrimeDecompositionResponse{
				Factor: f,
			})
			p = p / f
			time.Sleep(1 * time.Second)
		} else {
			f++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Received Compute Average stream RPC\n")
	sum := 0
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageReponse{
				Avg: float64(sum) / float64(count),
			})
		}
		sum += int(req.GetNumber())
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Received FindMaximum stream RPC\n")
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending response to client: %v\n", err)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Received Square Root RPC %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Error(codes.InvalidArgument,
			fmt.Sprintf("Received an negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Caculator server..")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("cannot listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("cannot server : %v", err)
	}
}
