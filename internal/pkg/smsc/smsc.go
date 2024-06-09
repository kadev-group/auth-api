package smsc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SMSc struct {
	http        *http.Client
	login       string
	password    string
	host        string
	reserveHost string
}

func NewSMSc(login, password string) *SMSc {
	return &SMSc{
		http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		login:       login,
		password:    password,
		host:        "smsc.kz",
		reserveHost: "www2.smsc.kz",
	}
}

func (s *SMSc) Send(ctx context.Context, phone, message string) error {
	request, err := s.prepareRequest(ctx, s.host, phone, message)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	err = s.sendRequest(request)
	if err != nil {
		if s.noNeedToRetryErrors(err) {
			return err
		}

		err = s.tryWithReservedHost(ctx, phone, message)
		if err != nil {
			return fmt.Errorf("failed to send with reserved host: %w", err)
		}
	}

	return nil
}

func (s *SMSc) tryWithReservedHost(ctx context.Context, phone, message string) error {
	request, err := s.prepareRequest(ctx, s.reserveHost, phone, message)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	err = s.sendRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMSc) sendRequest(request *http.Request) error {
	response, err := s.http.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: %d %s", response.StatusCode, response.Status)
	}

	result := struct {
		ErrorMessage string `json:"error,omitempty"`
		ErrorCode    int    `json:"error_code,omitempty"`
		ID           int    `json:"id,omitempty"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if result.ErrorCode != 0 {
		return Parse(result.ErrorCode, result.ErrorMessage)
	}

	return nil
}

func (s *SMSc) prepareRequest(ctx context.Context, host, phone, message string) (*http.Request, error) {
	endpoint := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "sys/send.php",
	}

	bodyParams := &url.Values{}
	bodyParams.Set("login", s.login)
	bodyParams.Set("psw", s.password)
	bodyParams.Set("phones", phone)
	bodyParams.Set("mes", message)
	bodyParams.Set("fmt", "3")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), strings.NewReader(bodyParams.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

func (s *SMSc) noNeedToRetryErrors(err error) bool {
	for _, e := range possibleErrors {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}
