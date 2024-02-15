package jwtAuth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gtkmk/finder_api/infra/envMode"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuth_GenerateJWT(t *testing.T) {
	jwt := NewjwtAuth(os.Getenv(envMode.JwtSecretConst))
	id := "user123"
	layer := uint64(1)

	token, err := jwt.GenerateJWT(map[string]interface{}{
		"i": id,
		"l": layer,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJwtAuth_CheckJwt_ValidToken(t *testing.T) {
	jwt := NewjwtAuth(os.Getenv(envMode.JwtSecretConst))
	id := "user123"
	layer := uint64(1)
	expiration := time.Now().Add(time.Minute * 120).Unix()

	tokenString, _ := jwt.GenerateJWT(map[string]interface{}{
		"i": id,
		"l": layer,
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "token", Value: tokenString}
	req.AddCookie(cookie)

	claims, err := jwt.CheckJwt(req)
	mapClaims := *claims

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, true, mapClaims["authorized"])
	assert.Equal(t, id, mapClaims["i"])
	assert.Equal(t, float64(expiration), mapClaims["expiration"])
}

func TestJwtAuth_CheckJwt_ExpiredToken(t *testing.T) {
	jwt := NewjwtAuth(os.Getenv(envMode.JwtSecretConst))
	id := "user123"
	layer := uint64(1)

	tokenString, _ := jwt.GenerateJWT(map[string]interface{}{
		"i": id,
		"l": layer,
	})
	tokenString = tokenString + "modified"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "token", Value: tokenString}
	req.AddCookie(cookie)

	claims, err := jwt.CheckJwt(req)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.EqualError(t, err, "token signature is invalid: signature is invalid")
}

func TestJwtAuth_CheckJwt_InvalidToken(t *testing.T) {
	jwt := NewjwtAuth(os.Getenv(envMode.JwtSecretConst))
	id := "user123"
	layer := uint64(1)

	tokenString, _ := jwt.GenerateJWT(map[string]interface{}{
		"i": id,
		"l": layer,
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "token", Value: tokenString + "modified"}
	req.AddCookie(cookie)

	claims, err := jwt.CheckJwt(req)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.EqualError(t, err, "token signature is invalid: signature is invalid")
}
