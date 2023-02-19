package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/circleci/backplane-go/environment"
	"github.com/circleci/samwise/observability"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
	"github.com/davidmdm/connect-playground/internal/signer"
	"github.com/davidmdm/muxter"
)

func main() {
	ctx, teardown := observability.MustReportingContext(context.Background(), environment.FromOS())
	defer teardown(ctx)

	privKey, pubKey := MustGenerateKeyPair()

	set := jwk.NewSet()
	set.AddKey(pubKey)

	handler := ServerHandler(ctx, signer.MakeService(privKey), set)

	http.ListenAndServe(":8080", h2c.NewHandler(handler, &http2.Server{}))
}

func ServerHandler(ctx context.Context, svc signer.Service, jwkSet jwk.Set) http.Handler {
	mux := muxter.New(muxter.MatchTrailingSlash(true))

	SignerServicePath, signerHandler := signerv1connect.NewSignerServiceHandler(svc)

	mux.UseGlobal(muxter.Logger(os.Stderr, func(o muxter.RespOverview) string {
		pattern := func() string {
			if o.Code == 404 {
				return o.Request.URL.Path
			}
			if pattern := o.Context.Pattern(); pattern != SignerServicePath {
				return pattern
			}
			return o.Request.URL.Path
		}()
		return fmt.Sprintf("%d %s %s %s", o.Code, o.Request.Method, pattern, o.TimeElapsed)
	}))

	mux.Handle(SignerServicePath, muxter.Adaptor(signerHandler, muxter.NoContext))

	mux.HandleFunc("/jwks.json", func(w http.ResponseWriter, r *http.Request, c muxter.Context) {
		json.NewEncoder(w).Encode(jwkSet)
	})

	mux.HandleFunc("/jwks/:id", func(w http.ResponseWriter, r *http.Request, c muxter.Context) {
		key, _ := jwkSet.LookupKeyID(c.Param("id"))
		json.NewEncoder(w).Encode(key)
	})

	return mux
}

func MustGenerateKeyPair() (priv, pub jwk.Key) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	priv, err = jwk.FromRaw(rsaKey)
	if err != nil {
		panic(err)
	}
	priv.Set("alg", "RS256")
	priv.Set("kid", "test")

	pub, err = priv.PublicKey()
	if err != nil {
		panic(err)
	}
	return
}
