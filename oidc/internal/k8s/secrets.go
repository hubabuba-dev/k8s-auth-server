package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetUser(username string) (map[string]string, error) {
	user := make(map[string]string)

	clientset, err := GetConfig()
	if err != nil {
		return nil, err
	}

	secret_name := "user" + "-" + username
	secret, err := clientset.CoreV1().Secrets("default").Get(
		context.TODO(),
		secret_name,
		metav1.GetOptions{},
	)

	if err != nil {
		return nil, err
	}

	userBytes := secret.Data["username"]
	passwordBytes := secret.Data["password"]
	user[string(userBytes)] = string(passwordBytes)

	return user, nil
}
