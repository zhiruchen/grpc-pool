package grpc_pool

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

var (
	defaultSize = 5
)

type Pool struct {
	mu      sync.RWMutex
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
	p.mu.RLock()
	conns := p.connChs
	df := p.df
	p.mu.RUnlock()

	select {
	case conn := <-conns:
		return conn, nil
	default:
		return df(p.target, p.opts...)
	}
}

func (p *Pool) Put(conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("put nil conn")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	select {
	case p.connChs <- conn:
		return nil
	default:
		return conn.Close() // pool已满，关闭put的连接
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	conns := p.connChs
	p.connChs = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}
