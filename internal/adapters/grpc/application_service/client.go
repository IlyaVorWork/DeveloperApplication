package application_service_client

import (
	"context"

	gen "developerApplication/internal/adapters/grpc/application_service/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client gen.ApplicationServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: gen.NewApplicationServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateApplication(ctx context.Context, req *gen.CreateApplicationRequest) (string, error) {
	resp, err := c.client.CreateApplication(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.ApplicationID, nil
}

func (c *Client) CreateRepository(ctx context.Context, req *gen.CreateRepositoryRequest) (string, error) {
	resp, err := c.client.CreateRepository(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.RepositoryID, nil
}
