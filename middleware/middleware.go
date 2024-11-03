package middleware

import (
	"context"
	"gopher-toolbox/token"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(token *token.Paseto) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if len(authorizationHeader) == 0 {
			slog.Error("authorization header is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization_header_required"})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			slog.Error("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_authorization_header"})
			return
		}

		accessToken := fields[1]
		payload, err := token.Verify(accessToken)
		if err != nil {
			slog.Error("error verifying token", "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_signature"})
			return
		}

		c.Set("payload", payload)
		c.Next()
	}
}

func PermissionMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// payloadInterface, exists := c.Get("payload")
		// if !exists {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		// 	return
		// }

		// payload, ok := payloadInterface.(*token.Payload)
		// if !ok {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid payload type"})
		// 	return
		// }

		// for _, requiredPermission := range requiredPermissions {
		// 	hasPermission, exists := payload.User.Permissions[requiredPermission]
		// 	if !exists || !hasPermission {
		// 		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission '%s' required", requiredPermission)})
		// 		return
		// 	}
		// }

		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
