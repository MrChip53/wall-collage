package client

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"wall-collage/pb"
)

func NewClientWithError(socketPath string) (pb.WallCollageClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient("unix:"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}))

	client := pb.NewWallCollageClient(conn)

	return client, conn, err
}

func NewClient(socketPath string) (pb.WallCollageClient, *grpc.ClientConn) {
	client, conn, err := NewClientWithError(socketPath)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}
	return client, conn
}
