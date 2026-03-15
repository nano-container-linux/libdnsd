package libdnsd

import "testing"

func TestNormalizeName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"EXAMPLE.com", "example.com."},
		{"  foo.  ", "foo."},
		{"", "."},
		{"bar", "bar."},
		{"baz.", "baz."},
	}
	for _, tt := range cases {
		got := NormalizeName(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestAcmeChallengeFQDN(t *testing.T) {
	cases := map[string]string{
		"example.com":                  "_acme-challenge.example.com.",
		"example.com.":                 "_acme-challenge.example.com.",
		"*.example.com.":               "_acme-challenge.example.com.",
		"_acme-challenge.example.com.": "_acme-challenge.example.com.",
		"_ACME-CHALLENGE.EXAMPLE.COM":  "_acme-challenge.example.com.",
	}
	for in, want := range cases {
		got := AcmeChallengeFQDN(in)
		if got != want {
			t.Errorf("AcmeChallengeFQDN(%q) = %q, want %q", in, got, want)
		}
	}
}
