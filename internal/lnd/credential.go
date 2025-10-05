package lnd

import "context"

// MacaroonCredential implements the credentials.PerRPCCredentials interface
type MacaroonCredential struct {
	MacaroonHex string
}

func (m *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"macaroon": m.MacaroonHex,
	}, nil
}

func (m *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}
