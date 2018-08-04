package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sysu-go-online/service-end/model"
)

// LogInHandler handler login event
func LogInHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if valid, err := CheckJWT(token); err == nil && valid {
		w.Header().Add("Authorization", token)
		w.WriteHeader(200)
		return nil
	}

	// check email and password
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return nil
	}
	user := UserController{}
	if err := json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(400)
		return nil
	}
	if ok := CheckUsername(user.Username); !ok {
		w.WriteHeader(400)
		return nil
	}
	user.User.Username = user.Username

	// create session to query user
	session := MysqlEngine.NewSession()
	has, err := user.User.GetWithUsername(session)
	if err != nil {
		session.Rollback()
		return err
	}
	if !has {
		w.WriteHeader(400)
		return nil
	}
	if !CompasePassword(user.Password, user.User.Password) {
		w.WriteHeader(400)
		return nil
	}

	// Generate token for user
	if token, err := GenerateJWT(user.User.Username); err != nil {
		return err
	} else {
		w.Header().Add("Authorization", token)
	}
	return nil
}

// LogOutHandler handler logout event
func LogOutHandler(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if valid, err := ValidateToken(token); err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return nil
	} else {
		if !valid {
			w.WriteHeader(401)
			return nil
		}
		if err = model.AddInvalidJWT(token, RedisClient); err != nil {
			return err
		}
		w.WriteHeader(200)
		return nil
	}
}
