package timetree

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

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
	expectedMac := []byte(hex.EncodeToString(macSum))

	return hmac.Equal(actualMac, expectedMac)
}
