package signer

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"

	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
)

//go:generate moq -out ./client_mock.go -rm . Client
type Client = signerv1connect.SignerServiceClient

var dialer = net.Dialer{}

func NewClient(client *http.Client, url string) Client {
	if client == nil {
		client = http.DefaultClient
	}
	return signerv1connect.NewSignerServiceClient(client, url, connect.WithGRPC())
}

func H2C() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}
}
