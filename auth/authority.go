package auth

import "time"

const (
	AuthorityName = "Authority"
)

// Authority is the interface that issue and verify token.
type Authority interface {
	Issuer
	Verifier
}

type Issuer interface {
	Issue(identity *Identity, duration time.Duration) (string, int64, error)
}

type Verifier interface {
	Verify(token string) (*Identity, int64, error)
}

type Session struct {
	Identity      *Identity
	Authorization string
}
