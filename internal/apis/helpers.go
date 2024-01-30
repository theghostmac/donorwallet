package apis

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theghostmac/donorwallet/internal/jwtauth"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

// Extracts the user ID from the JWT token in the request handler.
func getUserIDFromToken(c *gin.Context) (uuid.UUID, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logger.Error("empty authorization header error ")
		return uuid.Nil, errors.New("authorization token not provided")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	userID, err := jwtauth.ValidateToken(tokenString)
	if err != nil {
		logger.Error("failed to authenticate the user: ", zap.String("user_id", userID.String()))
		return uuid.Nil, err
	}

	return *userID, nil
}