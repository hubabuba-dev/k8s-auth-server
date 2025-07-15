package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gokube/internal/auth"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"gokube/internal/utils"
)

type UserData struct {
	Username string
	Password string
}

type TokenHandler struct {
	signer auth.TokenSigner
}

func NewTokenHandler(signer auth.TokenSigner) *TokenHandler {
	log.Println("Creating token hAandler")
	return &TokenHandler{
		signer: signer,
	}
}

func (h *TokenHandler) HandleToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating token")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading body: ", err)
	}
	var user auth.UserData
	if err = json.Unmarshal(body, &user); err != nil {
		log.Println("Error while unmarshalling: ", err)
	}

	exists, err := auth.ValidateUser(user)
	if err != nil {
		log.Println("Error while validating user: ", err)
		http.Error(w, "Error while validating user", http.StatusInternalServerError)
	} else if !exists {
		log.Println("User not exists: ", user)
		http.Error(w, "User not exists", http.StatusUnauthorized)
	} else {
		log.Println("Creating jwt token")
		jwtStr, err := h.createJWT(user.Username)
		if err != nil {
			log.Println("Error while creating JWT: ", err)
			http.Error(w, "Error", http.StatusInternalServerError)
		} else {
			utils.ResponseSend(w, map[string]string{
				"access_token": jwtStr,
				"token_type":   "Bearer",
			})
		}
	}

}

func (h *TokenHandler) createJWT(username string) (string, error) {
	privateKey, currentId, err := h.signer.GetCurrentKey()

	if err != nil {
		return "", errors.New("key not found")
	}

	claims := jwt.MapClaims{
		"iss": "https://oidc-server.example.com",
		"sub": username,
		"aud": "kubernetes",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = currentId

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error while signing key: %w", err)
	}

	return tokenString, nil
}

func HandleDiscovery(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseSend(w, map[string]interface{}{
		"issuer":                 "https://oidc-server.example.com",
		"jwks_uri":               "https://oidc-server.example.com/keys",
		"token_endpoint":         "https://oidc-server.example.com/token",
		"authorization_endpoint": "https://oidc-server.example.com/auth",
		"scopes_supported":       []string{"openid", "email", "groups"},
	})
}
