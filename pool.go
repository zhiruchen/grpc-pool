package grpc_pool

import (
	"fmt"

	"google.golang.org/grpc"
)

var (
	defaultSize = 5
)

type Pool struct {
	connChs chan *grpc.ClientConn
	df      DailFunc
	target  string
	opts    []grpc.DialOption
}

type DailFunc func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

func NewGRPCPool(df DailFunc, size int, target string, opts []grpc.DialOption) *Pool {
	if size == 0 {
		size = defaultSize
	}

	return &Pool{
		connChs: make(chan *grpc.ClientConn, size),
		df:      df,
		target:  target,
		opts:    opts,
	}
}

func (p *Pool) Get() (*grpc.ClientConn, error) {
	select {
	case conn := <-p.connChs:
		return conn, nil
	default:
		return p.df(p.target, p.opts...)
	}
}

func (p *Pool) Put(conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("put nil conn")
	}

	select {
	case p.connChs <- conn:
		return nil
	default:
		return conn.Close() // pool已满，关闭put的连接
	}
}

func (p *Pool) Close() {
	if p.connChs == nil {
		return
	}

	close(p.connChs)
	for conn := range p.connChs {
		conn.Close()
	}
}
