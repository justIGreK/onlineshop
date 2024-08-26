package grpcrequest

import (
	"context"
	"log"

	"github.com/justIGreK/emailcheck/go/emailcheck"
	"google.golang.org/grpc"
)

type GrpcRequest struct {
	conn *grpc.ClientConn
}

func NewGrpcRequst(conn *grpc.ClientConn)*GrpcRequest{
	return &GrpcRequest{conn:conn}
}

func (g *GrpcRequest) GetRequest(email string) bool{
	client := emailcheck.NewEmailServiceClient(g.conn)
	req := &emailcheck.EmailRequest{Email: email}
	resp, err := client.ValidateEmail(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling ValidateEmail: %v", err)
	}

	log.Printf("Validation Result: %v", resp)
	return resp.IsValid
}
