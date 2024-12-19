package jwt

import (
	"fmt"
	"net/http"
	"project/helper"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// constant
const (
	REQUEST = "invalid request"
	TOKEN   = "invalid token"
)

type JWT struct {
	PrivateKey string
	PublicKey  string
	Log        *zap.Logger
	UserID     string
}

type customClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	IP    string `json:"ip"`
	jwt.StandardClaims
}

func NewJWT(privateKey, publicKey string, log *zap.Logger) JWT {
	return JWT{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Log:        log,
	}
}

func (j *JWT) CreateToken(email, ip string, ID string) (string, error) {
	//prepare private key parsing
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(j.PrivateKey))
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &customClaims{
		ID:             ID,
		Email:          email,
		IP:             ip,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

// JWT for API
func (j *JWT) AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.PublicKey))
		if err != nil {
			helper.BadResponse(c, "unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims := &customClaims{}
		tokenValue := c.GetHeader("Authorization")
		if len(tokenValue) == 0 {
			c.Abort()
			return
		}

		tkn, err := jwt.ParseWithClaims(strings.Split(tokenValue, "Bearer ")[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				helper.BadResponse(c, fmt.Sprintf("unexpected method: %s", token.Header["alg"]), http.StatusUnauthorized)
				c.Abort()
				return nil, err
			}
			return key, nil
		})

		if err != nil {
			helper.BadResponse(c, "fail to validate signature or session expired", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !tkn.Valid {
			helper.BadResponse(c, "invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("user-id", claims.ID)
		//j.UserID = claims.ID

		c.Next()
	}
}
