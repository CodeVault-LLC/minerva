package fingerprint

import (
	"context"
	"time"

	"github.com/codevault-llc/minerva/pkg/logger"
	pb "github.com/codevault-llc/minerva/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.FingerprintServiceClient
}

var FingerprintClient *Client

// NewClient initializes a new gRPC client for the Fingerprint Service.
func NewClient(address string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	logger.Log.Info("Connected to Fingerprint Service")
	return &Client{
		conn:   conn,
		client: pb.NewFingerprintServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() {
	c.conn.Close()
}

// AddFingerprint calls the AddFingerprint RPC.
func (c *Client) AddFingerprint(ctx context.Context, req *pb.AddFingerprintRequest) (*pb.AddFingerprintResponse, error) {
	return c.client.AddFingerprint(ctx, req)
}

// GetFingerprint calls the GetFingerprint RPC.
func (c *Client) GetFingerprint(ctx context.Context, req *pb.GetFingerprintRequest) (*pb.GetFingerprintResponse, error) {
	return c.client.GetFingerprint(ctx, req)
}

// MatchFingerprint calls the MatchFingerprint RPC.
func (c *Client) MatchFingerprint(ctx context.Context, req *pb.MatchFingerprintRequest) (*pb.MatchFingerprintResponse, error) {
	return c.client.MatchFingerprint(ctx, req)
}
