package libdnsd

import (
    "reflect"
    "testing"
)

func TestJSONCodecMarshalUnmarshal(t *testing.T) {
    type S struct {
        A string `json:"a"`
        N int    `json:"n"`
    }
    in := S{A: "hello", N: 42}
    codec := jsonCodec{}
    b, err := codec.Marshal(in)
    if err != nil {
        t.Fatalf("Marshal error: %v", err)
    }
    var out S
    if err := codec.Unmarshal(b, &out); err != nil {
        t.Fatalf("Unmarshal error: %v", err)
    }
    if !reflect.DeepEqual(in, out) {
        t.Fatalf("codec roundtrip mismatch: got %+v want %+v", out, in)
    }
}
