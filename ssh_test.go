package libdnsd

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "encoding/base64"
    "os"
    "testing"

    "golang.org/x/crypto/ssh"
)

func TestSignPayloadWithPrivateKey_Verify(t *testing.T) {
    // generate RSA key
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

    payload := "test-payload"
    pub, sigB64, err := SignPayloadWithPrivateKey(f.Name(), payload)
    if err != nil {
        t.Fatalf("SignPayloadWithPrivateKey returned error: %v", err)
    }
    if pub == "" || sigB64 == "" {
        t.Fatalf("empty pub or sig")
    }
    sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
    if err != nil {
        t.Fatalf("invalid base64 signature: %v", err)
    }
    parsedPub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pub))
    if err != nil {
        // ssh.ParseAuthorizedKey returns (PublicKey, comment, options, rest, err) in older versions
        // try the simple parse
        pk, err2 := ssh.ParsePublicKey([]byte(pub))
        if err2 != nil {
            t.Fatalf("parse public key failed: %v / %v", err, err2)
        }
        parsedPub = pk
    }
    // construct signature; assume "ssh-rsa" algorithm for verification
    sig := &ssh.Signature{Format: ssh.KeyAlgoRSA, Blob: sigBytes}
    if err := parsedPub.Verify([]byte(payload), sig); err != nil {
        t.Fatalf("signature verification failed: %v", err)
    }
}
