package model

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

// JWTKey defines the token key
var JWTKey = "go-online"

// AddInvalidJWT add invalid jwt to the database
func AddInvalidJWT(jwtString string, client *redis.Client) error {
	// validate jwt
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return JWTKey, nil
	})

	// parse time from jwt
	var exp int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
			return nil
		}
		exp = claims["exp"].(int64)
	} else {
		return err
	}

	return client.Set(jwtString, "", time.Until(time.Unix(exp, 0))).Err()
}

// IsJWTExist judge if the token in redis
func IsJWTExist(tokenString string, client *redis.Client) (bool, error) {
	_, err := client.Get("key2").Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
