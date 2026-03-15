package libdnsd

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "encoding/base64"
    "os"
    "strings"
    "testing"
)

func TestBuildDynamicSubmitRequest(t *testing.T) {
    // generate an RSA private key and write PEM to a temp file
    key, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        t.Fatalf("generate key: %v", err)
    }
    privBytes := x509.MarshalPKCS1PrivateKey(key)
    pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}

    f, err := os.CreateTemp("", "libdnsd-priv-*.pem")
    if err != nil {
        t.Fatalf("create temp file: %v", err)
    }
    defer os.Remove(f.Name())
    if err := pem.Encode(f, pemBlock); err != nil {
        t.Fatalf("write pem: %v", err)
    }
    f.Close()

    payload := "example-payload"
    req, err := BuildDynamicSubmitRequest(payload, f.Name(), false)
    if err != nil {
        t.Fatalf("BuildDynamicSubmitRequest returned error: %v", err)
    }
    if req.PayloadHCL != payload {
        t.Fatalf("unexpected payload: got %q want %q", req.PayloadHCL, payload)
    }
    if req.PublicKey == "" {
        t.Fatalf("public key empty")
    }
    if req.Signature == "" {
        t.Fatalf("signature empty")
    }
    if _, err := base64.StdEncoding.DecodeString(req.Signature); err != nil {
        t.Fatalf("signature is not valid base64: %v", err)
    }
    if !strings.HasPrefix(req.PublicKey, "ssh-rsa") {
        // ssh.MarshalAuthorizedKey returns a key line beginning with algorithm
        t.Fatalf("public key does not look like an ssh key: %q", req.PublicKey)
    }
}
