# mTLS gRPC Server and Client Demo

A demonstration of mutual TLS (mTLS) authentication using gRPC in Go. This project showcases how to implement bidirectional certificate authentication between a gRPC server and client.

## Overview

This repository contains a simple gRPC service that uses mutual TLS authentication. Unlike standard TLS where only the server presents a certificate, mTLS requires both the server and client to present certificates, providing stronger security guarantees.

### Key Features

- **Server-side mTLS**: Server requires and verifies client certificates
- **Client-side authentication**: Client presents its own certificate to the server
- **CA-based trust**: Both server and client certificates are signed by a common Certificate Authority (CA)
- **gRPC integration**: Demonstrates mTLS with gRPC using Protocol Buffers

## How mTLS Works in This Project

1. **Certificate Authority (CA)**: A CA certificate and private key are used to sign both server and client certificates
2. **Server Certificate**: The server presents its certificate to clients and uses the CA to verify client certificates
3. **Client Certificate**: The client presents its certificate to the server and uses the CA to verify the server's certificate
4. **Mutual Verification**: Both parties verify each other's certificates before establishing a secure connection

## Prerequisites

- Go 1.25.4 or later
- Protocol Buffer compiler (`protoc`) and Go plugins (for regenerating proto files)
- OpenSSL (for generating certificates)

## Project Structure

```
mtls_service/
├── certs/              # Certificate files (CA, server, client)
│   ├── ca.crt          # CA certificate
│   ├── ca.key          # CA private key
│   ├── server.crt      # Server certificate
│   ├── server.key      # Server private key
│   ├── client.crt      # Client certificate
│   └── client.key      # Client private key
├── client/             # Client implementation
│   └── main.go
├── server/             # Server implementation
│   └── main.go
├── proto/              # Protocol Buffer definitions
│   ├── hello.proto     # Service definition
│   ├── hello.pb.go     # Generated Go code
│   └── hello_grpc.pb.go # Generated gRPC code
├── go.mod              # Go module dependencies
└── go.sum              # Go module checksums
```

## Certificate Generation

Before running the server or client, you need to generate the necessary certificates. Here's a script to generate them:

```bash
#!/bin/bash

# Create certs directory
mkdir -p certs
cd certs

# Generate CA private key
openssl genrsa -out ca.key 2048

# Generate CA certificate
openssl req -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=MyCA"

# Generate server private key
openssl genrsa -out server.key 2048

# Generate server certificate signing request
openssl req -new -key server.key -out server.csr -subj "/CN=localhost"

# Generate server certificate signed by CA
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extensions v3_req -extfile <(
cat <<EOF
[v3_req]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF
)

# Generate client private key
openssl genrsa -out client.key 2048

# Generate client certificate signing request
openssl req -new -key client.key -out client.csr -subj "/CN=client"

# Generate client certificate signed by CA
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt

# Clean up CSR files
rm *.csr *.srl

cd ..
```

Save this as `generate_certs.sh`, make it executable (`chmod +x generate_certs.sh`), and run it to generate all required certificates.

## Building and Running

### 1. Install Dependencies

```bash
go mod download
```

### 2. Generate Protocol Buffer Code (if needed)

If you modify `proto/hello.proto`, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/hello.proto
```

### 3. Run the Server

In one terminal:

```bash
go run server/main.go
```

The server will listen on port `50051` and log: `mTLS Server listening on port 50051...`

### 4. Run the Client

In another terminal:

```bash
go run client/main.go
```

The client will connect to the server, send a greeting request, and display the server's response.

## How It Works

### Server Side (`server/main.go`)

1. Loads the server's certificate and private key
2. Loads the CA certificate to verify client certificates
3. Creates TLS credentials with `ClientAuth: tls.RequireAndVerifyClientCert` to enforce mTLS
4. Creates a gRPC server with these credentials
5. Registers the `Greeter` service and starts listening

### Client Side (`client/main.go`)

1. Loads the client's certificate and private key
2. Loads the CA certificate to verify the server's certificate
3. Creates TLS credentials with the client certificate and CA trust pool
4. Connects to the server using these credentials
5. Makes an RPC call to `SayHello` and displays the response

## Testing

To verify mTLS is working correctly:

1. **Successful connection**: Run both server and client - you should see a successful greeting message
2. **Failed connection without client cert**: Modify the client to not include its certificate - the connection should fail
3. **Failed connection with wrong CA**: Use a different CA certificate - verification should fail

## Security Notes

- **Never commit private keys**: The `certs/` directory should be in `.gitignore` (or only commit `.crt` files for reference)
- **Use strong keys**: In production, use at least 2048-bit RSA keys or ECDSA keys
- **Certificate expiration**: Set appropriate expiration dates for production certificates
- **CA security**: Protect your CA private key - anyone with access can create trusted certificates

## License

This is a demonstration project for educational purposes.
