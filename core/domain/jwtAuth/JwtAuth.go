package jwtAuth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gtkmk/finder_api/core/domain/helper"
)

type JwtAuth struct {
	expiration int64
	secretKey  []byte
}

func NewjwtAuth(
	secret string,
) *JwtAuth {
	return &JwtAuth{
		expiration: time.Now().Add(time.Minute * 120).Unix(),
		secretKey:  []byte(secret),
	}
}

func (jwtAuth *JwtAuth) GenerateJWT(tokenContent map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	permissionHashed, err := jwtAuth.generatePermissionJWT(tokenContent)

	if err != nil {
		return "", helper.ErrorBuilder(helper.ErrorGeneratingEncryptionConst, err.Error())
	}

	unityHashed, err := jwtAuth.generateUnityJWT(tokenContent)

	if err != nil {
		return "", helper.ErrorBuilder(helper.ErrorGeneratingEncryptionConst, err.Error())
	}

	claims["authorized"] = true
	claims["l"] = tokenContent["l"]
	claims["i"] = tokenContent["i"]
	claims["p"] = permissionHashed
	claims["u"] = unityHashed
	claims["expiration"] = jwtAuth.expiration

	tokenString, err := token.SignedString(jwtAuth.secretKey)

	if err != nil {
		return "", helper.ErrorBuilder(helper.SomethingWentWrongConst, err.Error())
	}

	return tokenString, nil
}

func (jwtAuth *JwtAuth) generatePermissionJWT(tokenContent map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["p"] = tokenContent["p"]
	claims["expiration"] = jwtAuth.expiration

	tokenString, err := token.SignedString(jwtAuth.secretKey)

	if err != nil {
		return "", helper.ErrorBuilder(helper.SomethingWentWrongConst, err.Error())
	}

	return tokenString, nil
}

func (jwtAuth *JwtAuth) CheckJwt(request *http.Request) (*jwt.MapClaims, error) {
	tokenCookie, err := request.Cookie("token")

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenCookie.Value, jwtAuth.decryptJwt)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if err := jwtAuth.isTokenExpired(claims); err != nil {
		return nil, err
	}

	if !ok && !token.Valid {
		return nil, helper.ErrorBuilder(helper.UnauthorizedConst)
	}

	return &claims, nil
}

func (jwtAuth *JwtAuth) CheckPermissionJWT(jwtToken string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, jwtAuth.decryptJwt)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if err := jwtAuth.isTokenExpired(claims); err != nil {
		return nil, err
	}

	if !ok && !token.Valid {
		return nil, helper.ErrorBuilder(helper.UnauthorizedConst)
	}

	return &claims, nil
}

func (jwtAuth *JwtAuth) CheckJwtWallet(request *http.Request) (*jwt.MapClaims, error) {
	tokenCookie := request.Header.Get("Authorization")

	token, err := jwt.Parse(tokenCookie, jwtAuth.decryptJwt)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if err := jwtAuth.isTokenExpired(claims); err != nil {
		return nil, err
	}

	if !ok && !token.Valid {
		return nil, helper.ErrorBuilder(helper.UnauthorizedConst)
	}

	return &claims, nil
}

func (jwtAuth *JwtAuth) isTokenExpired(claims jwt.MapClaims) error {
	tokenExpiration, ok := claims["expiration"].(float64)
	if !ok {
		return helper.ErrorBuilder(helper.ExpiredTokenConst)
	}

	expirationTime := time.Unix(int64(tokenExpiration), 0)
	currentTime := time.Now()

	if currentTime.After(expirationTime) {
		return helper.ErrorBuilder(helper.ExpiredTokenConst)
	}

	return nil
}

func (jwtAuth *JwtAuth) decryptJwt(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, helper.ErrorBuilder(helper.ErrorParsingTokenConst)
	}
	return jwtAuth.secretKey, nil
}

func (jwtAuth *JwtAuth) generateUnityJWT(tokenContent map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["u"] = tokenContent["u"]
	claims["expiration"] = jwtAuth.expiration

	tokenString, err := token.SignedString(jwtAuth.secretKey)

	if err != nil {
		return "", helper.ErrorBuilder(helper.SomethingWentWrongConst, err.Error())
	}

	return tokenString, nil
}
