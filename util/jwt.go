package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"time"

	"github.com/pkg/errors"
)

var encoder = base64.URLEncoding.WithPadding(base64.NoPadding)
var header = encoder.EncodeToString([]byte(`{"alg": "RS256", "typ": "JWT"}`))

func GetToken(appID int, key string) (string, error) {
	current := time.Now()
	payload := map[string]interface{}{
		"iat": current.Unix(),
		"exp": current.Add(time.Minute * 5).Unix(),
		"iss": appID,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", errors.WithStack(err)
	}

	p := encoder.EncodeToString(b)
	body := string(header) + "." + string(p)
	h := crypto.Hash.New(crypto.SHA256)
	h.Write([]byte(body))
	hashed := h.Sum(nil)

	block, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.WithStack(err)
	}

	data, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", errors.WithStack(err)
	}

	signature := encoder.EncodeToString(data)

	return body + "." + signature, nil
}
