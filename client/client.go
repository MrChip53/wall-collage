package client

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"wall-collage/pb"
)

func NewClient(socketPath string) (pb.WallCollageClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("unix:"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewWallCollageClient(conn)

	return client, conn
}
