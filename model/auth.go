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
		return []byte(JWTKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// parse time from jwt
	var exp int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expired := claims["exp"]
		if expired == nil {
			return nil
		}
		exp = int64(expired.(float64))
		if time.Now().Unix() > exp {
			return nil
		}
	} else {
		return nil
	}

	return client.Set(jwtString, "", time.Until(time.Unix(exp, 0))).Err()
}

// IsJWTExist judge if the token in redis
func IsJWTExist(tokenString string, client *redis.Client) (bool, error) {
	_, err := client.Get(tokenString).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
