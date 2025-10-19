package thriftclient

import (
	"context"
	"monitoring-by-thrift/internal/domain"
	"time"
)

type Client struct {
	c *domain.ExternalService
}

func New(addr string, timeout time.Duration) (*Client, error) {
	/*transport, err := thrift.NewTSocketTimeout("http://")
	if err != nil {
		return nil, err
	}*/
	//	proto := thrift.NewTBinaryProtocolFactoryDefault()
	/*client := domain.ExternalService(transport) //extsvc.NewExternalServiceClientFactory(transport, proto)
	if err := transport.Open(); err != nil {
		return nil, err
	}*/
	return &Client{c: nil}, nil
}

func (cl *Client) ProcessStatus(ctx context.Context, transactionID string, proxyParams *map[string]string) (string, error) {
	// адаптируйте под ваш IDL
	resp := "t" ///cl.c.NormalizeOwner(ctx, ownerID)
	/*if err != nil {
		return "", err
	}*/
	return resp, nil
}
