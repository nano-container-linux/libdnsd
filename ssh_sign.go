package libdnsd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SignString choisit la méthode de signature selon les paramètres.
func SignString(signed string, privateKeyPath string, useAgent bool) (string, string, error) {
	if strings.TrimSpace(privateKeyPath) != "" {
		return SignPayloadWithPrivateKey(privateKeyPath, signed)
	}
	if !useAgent {
		return "", "", fmt.Errorf("no private key provided and ssh-agent usage disabled")
	}
	return SignPayloadWithAgent(signed)
}

// SignPayloadWithPrivateKey signe un payload avec une clé privée SSH.
func SignPayloadWithPrivateKey(privateKeyPath string, payload string) (string, string, error) {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read private key %s: %w", privateKeyPath, err)
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse private key %s: %w", privateKeyPath, err)
	}
	sig, err := signer.Sign(rand.Reader, []byte(payload))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign payload with private key %s: %w", privateKeyPath, err)
	}
	pub := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey())))
	sigRaw := sig.Blob
	return pub, base64.StdEncoding.EncodeToString(sigRaw), nil
}

// SignPayloadWithAgent signe un payload avec la première clé disponible dans ssh-agent.
func SignPayloadWithAgent(payload string) (string, string, error) {
	sock := strings.TrimSpace(os.Getenv("SSH_AUTH_SOCK"))
	if sock == "" {
		return "", "", fmt.Errorf("SSH_AUTH_SOCK is not set")
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}
	defer conn.Close()
	agentClient := agent.NewClient(conn)
	signers, err := agentClient.Signers()
	if err != nil {
		return "", "", fmt.Errorf("failed to list ssh-agent signers: %w", err)
	}
	if len(signers) == 0 {
		return "", "", fmt.Errorf("ssh-agent has no loaded keys")
	}
	signer := signers[0]
	sig, err := signer.Sign(rand.Reader, []byte(payload))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign payload with ssh-agent key: %w", err)
	}
	pub := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey())))
	sigRaw := sig.Blob
	return pub, base64.StdEncoding.EncodeToString(sigRaw), nil
}
