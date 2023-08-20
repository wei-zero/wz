package jwtauthority

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wei-zero/wz/auth"
)

var (
	// ErrTokenExpired indicate expired token
	ErrTokenExpired = errors.New("token is either expired or not active yet")
)

// authority encode/decode token
type authority struct {
	algorithm string // signature and hash algorithm
	secret    string // secret for signature signing and verification. can be replaced with certificate.
	expireIn  int64
}

func (a *authority) Issue(identity *auth.Identity, duration time.Duration) (string, int64, error) {
	if duration == 0 {
		duration = time.Duration(a.expireIn) * time.Second
	}
	expiresAt := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"DeviceID": identity.DeviceID,
		"UserID":   identity.UserID,
		"Roles":    identity.Roles,
		"Internal": identity.Internal,
		"exp":      expiresAt,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(a.algorithm), claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

func (a *authority) Verify(tokenStr string) (identity *auth.Identity, expireAt int64, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return
	}

	if !token.Valid {
		err = fmt.Errorf("token invalid")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identity = &auth.Identity{}
		if claims["DeviceID"] != nil {
			identity.DeviceID = claims["DeviceID"].(string)
		}

		var uid float64
		uid, ok = claims["UserID"].(float64)

		if claims["Internal"] != nil {
			identity.Internal = claims["Internal"].(string)
		} else {
			if !ok {
				err = errors.New("UserID not valid")
				return
			}
		}

		identity.UserID = int64(uid)
		var _roles []interface{}
		_roles, ok = claims["Roles"].([]interface{})
		if ok {
			for _, role := range _roles {
				identity.Roles = append(identity.Roles, role.(string))
			}
		}
		expireAt = int64(claims["exp"].(float64))
	} else {
		err = errors.New("get claims failed or token invalid")
		return
	}
	return
}

// New create new codec
func New(algorithm string, secret string, expireIn int64) auth.Authority {
	return &authority{
		algorithm: algorithm,
		secret:    secret,
		expireIn:  expireIn,
	}
}
