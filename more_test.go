package libdnsd

import (
    "encoding/json"
    "testing"
)

func TestBuildDynamicSubmitRequest_NoKeyError(t *testing.T) {
    _, err := BuildDynamicSubmitRequest("payload", "", false)
    if err == nil {
        t.Fatalf("expected error when no private key and agent disabled")
    }
}

func TestSignPayloadWithPrivateKey_InvalidPath(t *testing.T) {
    _, _, err := SignPayloadWithPrivateKey("/non/existent/path.pem", "payload")
    if err == nil {
        t.Fatalf("expected error for invalid private key path")
    }
}

func TestJSONCodecName(t *testing.T) {
    var c jsonCodec
    if c.Name() != "json" {
        t.Fatalf("jsonCodec.Name() = %q, want \"json\"", c.Name())
    }
}

func TestAcmeTokenCreateRequest_JSONRoundtrip(t *testing.T) {
    in := AcmeTokenCreateRequest{
        FQDN:      "example.com",
        Signature: "sig",
        PublicKey: "pub",
    }
    b, err := json.Marshal(in)
    if err != nil {
        t.Fatalf("marshal error: %v", err)
    }
    var out AcmeTokenCreateRequest
    if err := json.Unmarshal(b, &out); err != nil {
        t.Fatalf("unmarshal error: %v", err)
    }
    if out.FQDN != in.FQDN || out.Signature != in.Signature || out.PublicKey != in.PublicKey {
        t.Fatalf("roundtrip mismatch: got %+v want %+v", out, in)
    }
}
