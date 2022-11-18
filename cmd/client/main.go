package main

import (
	"context"
	"fmt"
	"net/http"

	signerv1 "github.com/davidmdm/connect-playground/internal/proto/signer/v1"

	signer "github.com/davidmdm/connect-playground/pkgs/clients/go"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := signer.NewClient(http.DefaultClient, "http://localhost:8080")

	res, err := client.Sign(context.Background(), connect.NewRequest(&signerv1.SignRequest{
		OrgId:    "org_uuid",
		Subject:  "me-mario",
		Audience: []string{"ecr"},
	}))
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Msg.Token)
}
