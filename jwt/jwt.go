package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

var mySecret []byte

const (
	TokenExpireDuration   = time.Hour * 24 * 1
	RefreshExpireDuration = time.Hour * 24 * 7
	Issuer                = "cygnus"
)

func Init(secret string) {
	mySecret = []byte(secret)
}

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

func GenToken(userID uint64) (aToken, rToken string, err error) {
	c := MyClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    Issuer,
		},
	}

	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	if err != nil {
		return "", "", err
	}
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshExpireDuration).Unix(),
		Issuer:    Issuer,
	}).SignedString(mySecret)
	return
}

func ParseToken(tokenString string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID)
	}
	return
}
