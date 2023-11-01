package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type TokenDetails struct {
	Token     *string
	TokenUuid string
	USerID    string
	ExpiresIn *int64
}

// CreateToken uses the private ket to sign JSON Web token through the RS256 algorithm.
func CreateToken(userId string, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := time.Now().UTC()

	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.TokenUuid = uuid.NewV4().String()
	td.USerID = userId

	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userId
	atClaims["token_uuid"] = td.TokenUuid
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	return td, nil
}

// ValidateToken will utilize the public keys to authenticate JSON Web Tokens.
func ValidateToken(token, publicKey string) (*TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse token public key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		USerID:    fmt.Sprint(claims["sub"]),
	}, nil
}
