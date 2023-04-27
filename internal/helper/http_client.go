package helper

import (
	"crypto/tls"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTelegramBotHTTPClient() tgbotapi.HTTPClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // nolint: gosec
			ServerName:         "api.telegram.org",
		},
	}
	client := &http.Client{Transport: transport}

	return &telegramBotHTTPClientImpl{
		httpClient: client,
	}
}

type telegramBotHTTPClientImpl struct {
	httpClient *http.Client
}

func (impl *telegramBotHTTPClientImpl) Do(req *http.Request) (*http.Response, error) {
	// req.Header.Add("Host", "api.telegram.org")
	req.Host = "api.telegram.org"

	return impl.httpClient.Do(req)
}
