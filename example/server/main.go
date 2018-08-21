package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/zhiruchen/grpc-pool/example/app"
)

type app struct{}

func (a *app) Hello(ctx context.Context, req *pb.AppReq) (*pb.AppResp, error) {
	return &pb.AppResp{
		Msg: req.Hello + " My Friend ",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8989")
	if err != nil {
		log.Println(err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAppServer(s, &app{})

	if err = s.Serve(lis); err != nil {
		fmt.Println("serve error: %v", err)
	}
}
