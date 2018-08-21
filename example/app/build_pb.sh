#!/usr/local/bin/fish

echo $GOPATH
protoc  \
	--proto_path=$GOPATH/src  \
	--proto_path=.  \
	--go_out=plugins=grpc:$GOPATH/src/github.com/zhiruchen/grpc-pool/example/app \
	test.proto
