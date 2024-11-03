package token

import (
	"crypto/ed25519"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type UserPayload struct {
	Username string `json:"username"`
	// TODO: Add permissions
}

type Payload struct {
	UUID      uuid.UUID   `json:"token_uuid"`
	User      UserPayload `json:"user"`
	IssuedAt  time.Time   `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

type Paseto struct {
	paseto     *paseto.V2
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func New() *Paseto {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	return &Paseto{
		paseto:     paseto.NewV2(),
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func NewPayload(user UserPayload) *Payload {
	// TODO: add documentation and advert to developers: tokenID != user.UUID
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return NewPayload(user)
	}

	payload := &Payload{
		UUID:      tokenID,
		User:      user,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
	}

	return payload
}

func (m *Paseto) Create(user UserPayload) (string, error) {
	return m.paseto.Sign(m.privateKey, NewPayload(user), nil)
}

func (m *Paseto) Verify(token string) (*Payload, error) {
	var payload Payload
	err := m.paseto.Verify(token, m.publicKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
