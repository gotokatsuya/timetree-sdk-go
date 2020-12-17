package timetree

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gotokatsuya/timetree-sdk-go/timetree/api"
)

const defaultAccessTokenLifetime = 600

type CalendarAppAuthenticator struct {
	applicationID            string
	accessTokenLifetimeInSec int
	privateKey               *rsa.PrivateKey

	client *api.Client
}

func NewCalendarAppAuthenticator(applicationID string, privateKey []byte) (*CalendarAppAuthenticator, error) {
	c := &CalendarAppAuthenticator{
		applicationID:            applicationID,
		accessTokenLifetimeInSec: defaultAccessTokenLifetime,
	}
	// PrivateKey
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("invalid private key data")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	c.privateKey = key

	cli, err := api.NewAuthenticatorClient(http.DefaultClient)
	if err != nil {
		return nil, err
	}
	c.client = cli
	return c, nil
}

// AccessTokenResponse type
type AccessTokenResponse struct {
	api.ErrorResponse
	AccessToken string `json:"access_token"`
	ExpireAt    int64  `json:"expire_at"`
	TokenType   string `json:"token_type"`
}

// AccessToken アクセストークンの取得
func (c *CalendarAppAuthenticator) AccessToken(ctx context.Context, installationID string) (*AccessTokenResponse, *http.Response, error) {
	path := fmt.Sprintf("/installations/%s/access_tokens", installationID)
	httpReq, err := c.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}
	token, err := c.generateToken()
	if err != nil {
		return nil, nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	resp := new(AccessTokenResponse)
	httpResp, err := c.client.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

func (c *CalendarAppAuthenticator) generateToken() (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(c.accessTokenLifetimeInSec) * time.Second).Unix(),
		Issuer:    c.applicationID,
	})
	return token.SignedString(c.privateKey)
}
