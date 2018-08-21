package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	grpcPool "github.com/zhiruchen/grpc-pool"
	pb "github.com/zhiruchen/grpc-pool/example/app"
)

var pool *grpcPool.Pool

func connectAppServer(address string) (*grpc.ClientConn, pb.AppClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewAppClient(conn)
	return conn, c, nil
}

func hello() error {
	conn, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Put(conn)

	c := pb.NewAppClient(conn)
	resp, err := c.Hello(context.Background(), &pb.AppReq{Hello: "Hello"})
	if err != nil {
		return err
	}

	log.Println(resp.Msg)
	return nil
}

func main() {
	// conn, client, err := connectAppServer("localhost:8989")
	// defer conn.Close()

	target := "localhost:8989"

	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	df := func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, target, opts...)
	}

	pool = grpcPool.NewGRPCPool(df, 2, target, opts)
	defer pool.Close()

	hello()
}
