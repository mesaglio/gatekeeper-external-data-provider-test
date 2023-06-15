package cosign

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/sigstore/pkg/signature"
)

type CosignValidator struct {
	sigVerifier signature.Verifier
}

func New(cosignPKB64 string) (*CosignValidator, error) {
	pemData, err := base64.StdEncoding.DecodeString(cosignPKB64)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("could not decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %v", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public is not a RSA public key")
	}

	sv, err := signature.LoadRSAPKCS1v15Verifier(rsaPub, crypto.SHA256)
	if err != nil {
		return nil, fmt.Errorf("could not load signature verifier: %v", err)
	}

	return &CosignValidator{
		sigVerifier: sv,
	}, nil
}

func (nv CosignValidator) ValidKey(key string, results []externaldata.Item) []externaldata.Item {
	ref, err := name.ParseReference(key)

	if err != nil {
		fmt.Printf("ERROR (ParseReference(%q)): %v\n", key, err)
		return results
	}

	opts := &cosign.CheckOpts{
		SigVerifier: nv.sigVerifier,
	}
	_, _, err = cosign.VerifyImageSignatures(nil, ref, opts)

	if err != nil {
		fmt.Printf("could not verify image signature cosign: %v\n", err)
		results = validators.WriteInvalidKey(key, results)
		return results
	}

	results = validators.WriteValidKey(key, results)
	return results
}
