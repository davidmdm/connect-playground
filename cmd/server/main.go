package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/davidmdm/connect-playground/internal/proto/signer/v1/signerv1connect"
	"github.com/davidmdm/connect-playground/internal/signer"
	"github.com/davidmdm/x/xcontext"
)

func main() {
	ctx, cancel := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer cancel()

	privKey, pubKey := MustGenerateKeyPair()

	set := jwk.NewSet()
	set.AddKey(pubKey)

	handler := ServerHandler(ctx, signer.MakeService(privKey), set)

	http.ListenAndServe(":8080", handler)
}

func ServerHandler(ctx context.Context, svc signer.Service, jwkSet jwk.Set) http.Handler {
	mux := http.NewServeMux()

	SignerServicePath, signerHandler := signerv1connect.NewSignerServiceHandler(svc)

	mux.Handle(SignerServicePath, signerHandler)

	mux.HandleFunc("/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(jwkSet)
	})

	mux.HandleFunc("/jwks/{id}", func(w http.ResponseWriter, r *http.Request) {
		key, _ := jwkSet.LookupKeyID(r.PathValue("id"))
		json.NewEncoder(w).Encode(key)
	})

	return h2c.NewHandler(TrafficLogger(mux), &http2.Server{})
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

func TrafficLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		fmt.Printf("%s %s - (%s)\n", r.Method, r.Pattern, time.Since(start))
	})
}
