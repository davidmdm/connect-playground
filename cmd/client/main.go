package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	signerv1 "github.com/davidmdm/connect-playground/internal/proto/signer/v1"
	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := signerv1connect.NewSignerServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	res, err := client.Sign(context.Background(), connect.NewRequest(&signerv1.SignRequest{
		OrgId:    "org_uuid",
		Subject:  "me-mario",
		Audience: []string{"ecr"},
	}))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res.Msg.Token)
}
