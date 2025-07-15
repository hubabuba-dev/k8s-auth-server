package cli

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"kubectl-login/internal/config"
	"kubectl-login/internal/kubectl"
)

type LoginRequset struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"token_type"`
}

type TokenSet struct {
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Type         string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	IssuedAt     time.Time `json:"issued_at"`
	Username     string    `json:"username"`
	Server       string    `json:"server"`
}

func LoginExecute(cfg *config.Config) error {
	err := checkCredentials(cfg)
	if err != nil {
		fmt.Println("Error checking creds")
		return err
	}

	resp, err := GetToken(cfg)
	if err != nil {
		fmt.Println("Error getting token")
		return err
	}

	clnt := kubectl.NewClient(cfg)

	if err := clnt.CreateConfig(resp.Token, resp.RefreshToken); err != nil {
		fmt.Println("Error Creating config")
		return err
	}

	return nil
}

func checkCredentials(cfg *config.Config) error {
	user := cfg.Username
	pass := cfg.Password
	interactive := cfg.Interactive
	if interactive || user == "" {
		fmt.Print("username: ")
		reader := bufio.NewReader(os.Stdin)
		user, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		cfg.Username = strings.TrimSpace(user)
	}

	if interactive || pass == "" {
		fmt.Print("password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		cfg.Password = string(password)
	}
	fmt.Print('\n')
	return nil
}

func GetToken(cfg *config.Config) (*LoginResponse, error) {
	client := &http.Client{Timeout: cfg.Timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	req_body := LoginRequset{Username: cfg.Username,
		Password: cfg.Password}

	json_body, err := json.Marshal(req_body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(cfg.Server+"/token", "application/json", bytes.NewBuffer(json_body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result LoginResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
