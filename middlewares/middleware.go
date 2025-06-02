package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"user-service/common/responses"
	"user-service/config"
	"user-service/constants"
	errorConstants "user-service/constants/errors"
	services "user-service/services/users"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc { // HANDLE PANIC KETIKA ADA ERROR TETAPI APLIKASI MASIH BISA BERJALAN
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, responses.Response{
					Status:   constants.Error,
					Messsage: errorConstants.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Abort()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc { // MEMBATASI AKSES USER PER MENIT SEMISAL DALAM 1 MENIT DIIZINKAN CUMA 100 AKSES
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, responses.Response{
				Status:   constants.Error,
				Messsage: errorConstants.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

// FUNCTION PRIVATE BUAT API KEY

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}

	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, responses.Response{
		Status:   constants.Error,
		Messsage: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errorConstants.ErrUnauthorized
	}

	return nil
}

func validateBearerToken(c *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		return errorConstants.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return errorConstants.ErrUnauthorized
	}

	claims := &services.Claims{}
	tokenJWT, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errorConstants.ErrInvalidToken
		}

		jwtSecret := []byte(config.Config.JWTSecretKey)
		return jwtSecret, nil
	})

	if err != nil || !tokenJWT.Valid {
		return errorConstants.ErrUnauthorized
	}

	userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claims.User))
	c.Request = userLogin
	c.Set(constants.Token, token)

	return nil
}

// FUNC BUAT PENCEGAHAN KALAU ADA ERROR

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)
		if token != "" {
			responseUnauthorized(c, errorConstants.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(c, token)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
