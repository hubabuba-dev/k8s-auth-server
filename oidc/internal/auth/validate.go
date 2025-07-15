package auth

import (
	"fmt"
	"gokube/internal/k8s"
)

type UserData struct {
	Username string
	Password string
}

func ValidateUser(user UserData) (bool, error) {
	secret, err := k8s.GetUser(user.Username)
	if err != nil {
		return false, err
	}

	pass, exist := secret[user.Username]
	if !exist {
		return false, fmt.Errorf("user not exists")
	}
	if user.Password != pass {
		return false, fmt.Errorf("wrong password")
	}
	return true, nil
}
