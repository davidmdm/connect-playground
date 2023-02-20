package main

import (
	"context"
	"fmt"
	"net/http"

	signerv1 "github.com/davidmdm/connect-playground/internal/proto/signer/v1"
	"google.golang.org/protobuf/types/known/structpb"

	signer "github.com/davidmdm/connect-playground/pkgs/clients/go"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := signer.NewClient(http.DefaultClient, "http://localhost:8080")

	customClaims, err := structpb.NewStruct(map[string]any{
		"project_id":  23,
		"context_ids": []any{"org-readonly", "github-readonly"},
	})
	if err != nil {
		panic(err)
	}

	request := connect.NewRequest(&signerv1.SignRequest{
		OrgId:        "org_uuid",
		Subject:      "me-mario",
		Audience:     []string{"ECR"},
		CustomClaims: customClaims,
	})

	res, err := client.Sign(context.Background(), request)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Msg.Token)
}
