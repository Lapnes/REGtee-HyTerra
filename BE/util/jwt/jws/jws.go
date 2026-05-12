package jws

import (
	variable "backend/constant"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
	ID          int
	PrivilageID int
	Email       string
}

// generateToken signs a short-lived access token (24 h).
func (claim *UserClaims) generateToken() (string, error) {
	claim.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	return t.SignedString([]byte(variable.PASSWORD_SALT))
}

// generateRefreshToken signs a refresh token (7 days).
func (claim *UserClaims) generateRefreshToken() (string, error) {
	refreshClaim := *claim // copy, so expiry doesn't alias
	refreshClaim.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, &refreshClaim)
	return t.SignedString([]byte(variable.PASSWORD_SALT))
}

// NewToken generates an access token and, optionally, a refresh token.
// BUG FIX: previously both token and refresh were generated with
// generateRefreshToken(); now the access token correctly uses generateToken().
func (claim *UserClaims) NewToken(withRefresh bool) (token string, refresh string, err error) {
	claim.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:   "services",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.New().String(),
	}

	token, err = claim.generateToken()
	if err != nil {
		return
	}
	if withRefresh {
		refresh, err = claim.generateRefreshToken()
	}
	return
}

func ParseToken(accessToken string, ignoreExpired bool) (*UserClaims, error) {
	token := strings.Replace(accessToken, "Bearer ", "", 1)
	var opts []jwt.ParserOption
	if ignoreExpired {
		opts = append(opts, jwt.WithoutClaimsValidation())
	}
	parsed, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(variable.PASSWORD_SALT), nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	return parsed.Claims.(*UserClaims), nil
}

func (claim *UserClaims) IsActive() bool {
	return claim.ExpiresAt.After(time.Now())
}

func ExtractClient(ah string) (*UserClaims, error) {
	ts := strings.Replace(ah, "Bearer ", "", 1)
	claimsMap, err := extractClaims(ts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token claims")
	}
	j, err := json.Marshal(&claimsMap)
	if err != nil {
		return nil, err
	}
	var claims UserClaims
	if err = json.Unmarshal(j, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(variable.PASSWORD_SALT), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid JWT token")
	}
	return claims, nil
}
