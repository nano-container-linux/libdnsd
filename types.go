package libdnsd

// Shared types for ACME and dynamic submission

type ACMETokenEntry struct {
	FQDN      string `json:"fqdn"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
	Revoked   bool   `json:"revoked"`
}

type AcmeTokenListRequest struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type AcmeTokenListResponse struct {
	Tokens []ACMETokenEntry `json:"tokens"`
}

type AcmeTokenCreateRequest struct {
	FQDN      string `json:"fqdn"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
}

type AcmeTokenCreateResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Path    string `json:"path"`
	ID      string `json:"id"`
	FQDN    string `json:"fqdn"`
}

type AcmeTokenRevokeRequest struct {
	FQDN      string `json:"fqdn"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
}

type AcmeTokenRevokeResponse struct {
	Message string `json:"message"`
	Revoked bool   `json:"revoked"`
}

// Types for dynamic submission
type DynamicSubmitRequest struct {
	PayloadHCL string `json:"payload_hcl"`
	PublicKey  string `json:"public_key"`
	Signature  string `json:"signature"`
}

type DynamicSubmitResponse struct {
	Message string `json:"message,omitempty"`
	ID      string `json:"id,omitempty"`
	Path    string `json:"path,omitempty"`
}
