package middleware

import (
	"gopher-toolbox/token"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

func Test_authMiddleware(t *testing.T) {

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	paseto := token.New()
	user := token.UserPayload{
		Username: "test",
		// Permissions: map[string]bool{
		// 	"view_customer": true,
		// },
	}
	publicToken, error := paseto.Create(user)
	require.NoError(t, error)
	payload, error := paseto.Verify(publicToken)
	require.NoError(t, error)

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+publicToken)

	authorizationHeader := c.GetHeader("Authorization")
	require.Equal(t, "Bearer "+publicToken, authorizationHeader)

	authMW := AuthMiddleware(paseto)
	authMW(c)

	payloadFromMW, exists := c.Get("payload")
	require.True(t, exists)
	require.Equal(t, payload, payloadFromMW)
}
