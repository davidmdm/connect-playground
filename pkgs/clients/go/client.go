package signer

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
)

type Client = signerv1connect.SignerServiceClient

func NewClient(client *http.Client, url string) Client {
	return signerv1connect.NewSignerServiceClient(client, url, connect.WithGRPC())
}
