package handlers

import (
	"encoding/base64"
	"fmt"
	"gokube/internal/utils"
	"log"
	"math/big"
	"net/http"
)

func (h *TokenHandler) HandleJWKS(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling JWKS")
	privateKey, currentId, err := h.signer.GetCurrentKey()
	publicKey := privateKey.PublicKey
	if err != nil {
		http.Error(w, "Error while getting current key", http.StatusInternalServerError)
	}

	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	eBytes := big.NewInt(int64(publicKey.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	jwk := map[string]interface{}{
		"kty": "RSA",
		"kid": currentId,
		"n":   n,
		"e":   e,
	}

	utils.ResponseSend(w, map[string]interface{}{
		"keys": []interface{}{jwk},
	})
	fmt.Println("JWKS was sent")
}
