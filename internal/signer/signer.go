package signer

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	signerv1 "github.com/davidmdm/connect-playground/internal/proto/signer/v1"
	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
)

type Service struct {
	key jwk.Key
	signerv1connect.UnimplementedSignerServiceHandler
}

func MakeService(key jwk.Key) Service {
	return Service{key: key}
}

func (svc Service) Sign(ctx context.Context, req *connect.Request[signerv1.SignRequest]) (*connect.Response[signerv1.SignResponse], error) {
	now := time.Now()

	builder := jwt.
		NewBuilder().
		Subject(req.Msg.Subject).
		Audience(req.Msg.Audience).
		Issuer(req.Msg.Issuer).
		IssuedAt(now).
		Expiration(now.Add(time.Hour))

	for k, v := range req.Msg.CustomClaims.AsMap() {
		builder.Claim("x.claim/"+k, v)
	}

	tkn, err := builder.Build()
	if err != nil {
		return nil, err
	}

	tkn.Options().Enable(jwt.FlattenAudience)

	signed, err := jwt.Sign(tkn, jwt.WithKey(jwa.RS256, svc.key))
	if err != nil {
		return nil, err
	}

	return &connect.Response[signerv1.SignResponse]{
		Msg: &signerv1.SignResponse{Token: string(signed)},
	}, nil
}
