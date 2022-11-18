package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/structpb"

	signerv1 "github.com/davidmdm/connect-playground/internal/proto/signer/v1"
	signer "github.com/davidmdm/connect-playground/pkgs/clients/go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	client := signer.NewClient(signer.H2C(), "http://localhost:8080")

	customClaims, err := structpb.NewStruct(map[string]any{
		"project_id":  23,
		"context_ids": []any{"org-readonly", "github-readonly"},
	})
	if err != nil {
		return fmt.Errorf("failed to construct custom claims: %w", err)
	}

	request := connect.NewRequest(&signerv1.SignRequest{
		Issuer:       "peach-castle",
		Subject:      "me-mario",
		Audience:     []string{"smash-bros"},
		CustomClaims: customClaims,
	})

	res, err := client.Sign(context.Background(), request)
	if err != nil {
		if codeErr := new(connect.Error); errors.As(err, &codeErr) {
			return fmt.Errorf("failed to sign: %d: %s", codeErr.Code(), codeErr.Message())
		}
		return fmt.Errorf("failed to sign: %v", err)
	}

	fmt.Println(res.Msg.Token)
	return nil
}
