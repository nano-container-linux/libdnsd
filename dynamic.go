package libdnsd

// BuildDynamicSubmitRequest mutualisé: construit la requête en signant le payload
func BuildDynamicSubmitRequest(payload string, privateKeyPath string, useAgent bool) (DynamicSubmitRequest, error) {
	pub, sig, err := SignString(payload, privateKeyPath, useAgent)
	if err != nil {
		return DynamicSubmitRequest{}, err
	}
	return DynamicSubmitRequest{PayloadHCL: payload, PublicKey: pub, Signature: sig}, nil
}

// SubmitDynamicPayload mutualisé: prépare la requête puis appelle la fonction de soumission fournie
func SubmitDynamicPayload(target string, payload string, privateKeyPath string, useAgent bool, submitFunc func(string, DynamicSubmitRequest) (*DynamicSubmitResponse, error)) (*DynamicSubmitResponse, error) {
	req, err := BuildDynamicSubmitRequest(payload, privateKeyPath, useAgent)
	if err != nil {
		return nil, err
	}
	return submitFunc(target, req)
}
