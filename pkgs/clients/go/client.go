package signer

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
	"golang.org/x/net/http2"
)

type Client = signerv1connect.SignerServiceClient

var dialer = net.Dialer{}

func NewClient(client *http.Client, url string) Client {
	if client == nil {
		client = http.DefaultClient
	}
	c := *client
	c.Transport = &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}
	return signerv1connect.NewSignerServiceClient(&c, url, connect.WithGRPC())
}
