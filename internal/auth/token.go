package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

// sign calculeaza semnatura HMAC-SHA256 peste date si o returneaza
// in format Base64URL
func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// GenerateToken creeaza si semneaza un token JWT valid 24 de ore
func GenerateToken(userID int64, role, secret string) (string, error) {
	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256", "typ":"JWT"}`))

	claims := Claims{
		UserID: userID,
		Role:   role,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	payloadBytes, err := json.Marshal(claims)

	if err != nil {
		return "", err
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := headerB64 + "." + payloadB64
	sig := sign(unsigned, secret)

	return unsigned + "." + sig, nil
}

// VerifyToken verifica semnatura si expirarea unui JWT si returneaza
// datele din payload
func VerifyToken(tokenStr, secret string) (Claims, error) {
	parts := strings.Split(tokenStr, ".")

	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	expectedSig := sign(parts[0]+"."+parts[1], secret)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return Claims{}, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])

	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Exp {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
