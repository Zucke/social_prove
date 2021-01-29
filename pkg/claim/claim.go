package claim

import (
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Errors
var (
	ErrInsufficientPrivileges    = errors.New("Insufficient privileges")
	ErrInvalidToken              = errors.New("invalid token")
	ErrInvalidClaim              = errors.New("invalid claim")
	ErrUserNotAuthorized         = errors.New("not authorized")
	ErrInvalidAutorizationFormat = errors.New("invalid autorization format")
)

// Claim what goes in token claims.
type Claim struct {
	jwt.StandardClaims
	ID   string `json:"id"`
	Role uint   `json:"role"`
}

// GenerateToken generete a new token.
func GenerateToken(signingString string, ID string, Role uint) (string, error) {
	claims := Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "User auth",
		},
		ID:   ID,
		Role: Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingString))
}

//TokenFromAuthorization get token from Authorization
func TokenFromAuthorization(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return "", ErrUserNotAuthorized
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", ErrInvalidAutorizationFormat
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", ErrInvalidAutorizationFormat
	}

	return l[1], nil
}

// GetFromToken get claims from a token string.
func GetFromToken(tokenString, signingString string) (*Claim, error) {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(signingString), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaim
	}

	role, ok := claims["role"].(uint)
	if !ok {
		return nil, ErrInsufficientPrivileges
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, ErrInsufficientPrivileges
	}

	return &Claim{ID: id, Role: role}, nil
}
