package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/yirenzhou/mtls_service/proto"
)

func main() {
	// 1. Load the Client's Certificate and Private Key
	clientCert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatalf("failed to load client key pair: %v", err)
	}

	// 2. Load the CA Certificate to verify the Server
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	// 3. Create the TLS Credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert}, // Our ID
		RootCAs:      certPool,                      // Trusted CAs for the Server
		ServerName:   "localhost",                   // Must match the SAN in server.crt
	})

	// 4. Create the gRPC Client Connection
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// 5. Call the RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "mTLS World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Server Response: %s", r.GetMessage())
}
