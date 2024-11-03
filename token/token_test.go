package token

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zepyrshut/esfaker"
)

func Test_New(t *testing.T) {
	paseto := New()

	require.NotNil(t, paseto.paseto)
	require.NotNil(t, paseto.privateKey)
	require.NotNil(t, paseto.publicKey)
}

func Test_NewPayload(t *testing.T) {
	user := createRandomUser()
	payload := NewPayload(user)

	require.True(t, isValidUUID(payload.UUID))
	require.Equal(t, user, payload.User)
}

func Test_CreateToken(t *testing.T) {
	token := New()
	user := createRandomUser()
	signature, err := token.Create(user)

	require.Nil(t, err)
	require.NotEmpty(t, signature)
}

func Test_VerifyToken(t *testing.T) {
	token := New()
	user := createRandomUser()
	signature, _ := token.Create(user)
	payload, err := token.Verify(signature)

	require.Nil(t, err)
	require.Equal(t, user, payload.User)
}

func Test_VerifyToken_InvalidToken(t *testing.T) {
	token := New()
	_, err := token.Verify("invalid-token")

	require.NotNil(t, err)
}

func isValidUUID(u uuid.UUID) bool {
	_, err := uuid.Parse(u.String())
	return err == nil
}

func createRandomUser() UserPayload {
	return UserPayload{
		Username: esfaker.Chars(5, 10),
	}
}
