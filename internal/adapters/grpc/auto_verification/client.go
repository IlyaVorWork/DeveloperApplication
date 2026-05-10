package auto_verification_client

import (
	"context"

	gen "developerApplication/internal/adapters/grpc/auto_verification/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client gen.AutoVerificationServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: gen.NewAutoVerificationServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) StartVerification(ctx context.Context, apkFileName string) (string, error) {
	resp, err := c.client.StartVerificationProcess(ctx, &gen.StartVerificationProcessRequest{
		ApkFileName: apkFileName,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (c *Client) GetVerification(ctx context.Context, processID string) ([]*gen.VerificationProcess, error) {
	resp, err := c.client.GetVerificationProcess(ctx, &gen.GetVerificationRequest{
		ProcessID: processID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Process, nil
}
