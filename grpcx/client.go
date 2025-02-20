package grpcx

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

// ClientCreator create a grpc client
type ClientCreator func(string, *grpc.ClientConn) interface{}

// GRPCClient is a grpc client
type GRPCClient struct {
	sync.RWMutex

	creator ClientCreator
	opts    *clientOptions
	clients map[string]interface{}
	conns   map[string]*grpc.ClientConn
}

// NewGRPCClient returns a GRPC Client
func NewGRPCClient(creator ClientCreator, opts ...ClientOption) *GRPCClient {
	copts := &clientOptions{}
	for _, opt := range opts {
		opt(copts)
	}

	return &GRPCClient{
		opts:    copts,
		creator: creator,
		clients: make(map[string]interface{}),
	}
}

// Close close
func (c *GRPCClient) Close() error {
	c.RLock()
	defer c.RUnlock()

	var err error
	for _, conn := range c.conns {
		err = conn.Close()
	}

	return err
}

// GetServiceClient returns a grpc client
func (c *GRPCClient) GetServiceClient(name string) (interface{}, error) {
	c.RLock()
	if cli, ok := c.clients[name]; ok {
		c.RUnlock()
		return cli, nil
	}
	c.RUnlock()

	client, err := c.createClient(name)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *GRPCClient) createClient(name string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if cli, ok := c.clients[name]; ok {
		return cli, nil
	}

	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcOptions = append(grpcOptions, grpc.WithIdleTimeout(c.opts.timeout)) //.WithTimeout(c.opts.timeout))
	//grpcOptions = append(grpcOptions, grpc.WithBlock())
	//grpcOptions = append(grpcOptions, grpc.WithResolvers(c.opts.resolver)) //.WithBalancer(grpc.RoundRobin(c.opts.resolver)))

	conn, err := grpc.NewClient(fmt.Sprintf("%s", c.opts.prefix), grpcOptions...)
	if err != nil {
		return nil, err
	}

	cli := c.creator(name, conn)
	c.clients[name] = cli

	return cli, nil
}
