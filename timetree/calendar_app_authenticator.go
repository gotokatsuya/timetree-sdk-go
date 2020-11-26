package timetree

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gotokatsuya/timetree-sdk-go/timetree/api"
)

const DEFAULT_ACCESS_TOKEN_LIFETIME = 600

type CalendarAppAuthenticator struct {
	applicationID            string
	privateKeyPath           string
	accessTokenLifetimeInSec int

	client *api.Client
}

func NewCalendarAppAuthenticator(applicationID, privateKeyPath string) (*CalendarAppAuthenticator, error) {
	c := &CalendarAppAuthenticator{
		applicationID:            applicationID,
		privateKeyPath:           privateKeyPath,
		accessTokenLifetimeInSec: DEFAULT_ACCESS_TOKEN_LIFETIME,
	}
	cli, err := api.NewClientWithoutAccessToken(http.DefaultClient)
	if err != nil {
		return nil, err
	}
	c.client = cli
	return c, nil
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

// AccessTokenResponse type
type AccessTokenResponse struct {
	api.ErrorResponse
	AccessToken string `json:"access_token"`
	ExpireAt    int64  `json:"expire_at"`
	TokenType   string `json:"token_type"`
}

func (c *CalendarAppAuthenticator) generateToken() (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(c.accessTokenLifetimeInSec) * time.Second).Unix(),
		Issuer:    c.applicationID,
	})
	key, err := c.readPrivateKey()
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

func (c *CalendarAppAuthenticator) readPrivateKey() (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(c.privateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("invalid private key data")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

type CalendarAppWebhook struct {
	secret string
}

func NewCalendarAppWebhook(secret string) *CalendarAppWebhook {
	return &CalendarAppWebhook{
		secret: secret,
	}
}

func (c CalendarAppWebhook) Verify(httpRequest *http.Request) bool {
	sha := strings.TrimPrefix(httpRequest.Header.Get("X-Timetree-Signature"), "sha1=")
	actualMac := []byte(sha)

	mac := hmac.New(sha1.New, []byte(c.secret))
	requestBody, _ := ioutil.ReadAll(httpRequest.Body)
	httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	mac.Write(requestBody)
	macSum := mac.Sum(nil)
	expectedMac := []byte(base64.StdEncoding.EncodeToString(macSum))

	return hmac.Equal(actualMac, expectedMac)
}
