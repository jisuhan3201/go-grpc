package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jisuhan3201/go-grpc/calculator/calculatorpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculate method invoked %v", req.GetCalculating())
	result := req.GetCalculating().GetLeft() + req.Calculating.GetRight()
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Caculator server..")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("cannot listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("cannot server : %v", err)
	}
}
