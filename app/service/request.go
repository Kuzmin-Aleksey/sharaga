package service

import (
	"app/failure"
	"app/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (s *Service) client() *http.Client {
	return &http.Client{
		Timeout: time.Duration(s.cfg.Timeout) * time.Second,
	}
}

func (s *Service) request(dst any, method string, u string, body any, urlValues ...[2]string) error {
	fLoop := true

	breakAttempts := func() {
		fLoop = false
	}

	var err error

	for attempt := 1; fLoop; attempt++ {
		err = func() error {
			var buf io.Reader

			if body != nil {
				b := bytes.NewBuffer([]byte{})
				if err := json.NewEncoder(b).Encode(body); err != nil {
					breakAttempts()
					return err
				}
				buf = b
			}

			var urlValuesStr string
			if len(urlValues) > 0 {
				urlValuesStr = "?"
				for _, kv := range urlValues {
					if kv[0] == "" || kv[1] == "" {
						continue
					}

					urlValuesStr += url.QueryEscape(kv[0]) + "=" + url.QueryEscape(kv[1]) + "&"
				}
				urlValuesStr = urlValuesStr[:len(urlValuesStr)-1]
			}
			req, err := http.NewRequest(method, s.cfg.Address+u+urlValuesStr, buf)
			if err != nil {
				breakAttempts()
				return err
			}

			req.Header.Set("Authorization", "Bearer "+s.tokens.AccessToken)

			resp, err := s.client().Do(req)
			if err != nil {
				err = failure.NewNetworkError(err)
				if attempt > 2 {
					breakAttempts()
				}
				return err
			}

			if resp.StatusCode != http.StatusOK {
				switch resp.StatusCode {
				case http.StatusUnauthorized:
					err = failure.NewUnauthorizedError()
					s.RefreshTokens()

					if attempt > 1 {
						breakAttempts()
					}
					return err
				default:
					jsError := new(models.Error)
					_ = json.NewDecoder(resp.Body).Decode(jsError)
					breakAttempts()
					return failure.NewServerError(jsError.Error, resp.StatusCode)
				}
			}

			if dst != nil {
				if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
					breakAttempts()
					return failure.NewInvalidResponseError(err)
				}
			}

			return nil

		}()

		if err != nil {
			log.Printf("%s %s: %s, attempt #%d", method, u, err, attempt)
		} else {
			breakAttempts()
		}
	}

	return err
}

func (s *Service) Ping() error {
	resp, err := s.client().Get(s.cfg.Address + urlPing)
	if err != nil {
		return failure.NewNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return failure.NewServerError(http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	return nil
}
