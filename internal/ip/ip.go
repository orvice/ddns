package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	ipConfigCoAddr = "https://ifconfig.co/json"
)

type Response struct {
	IP string `json:"ip"`
}

type Getter interface {
	GetIP() (string, error)
}

// httpDoer is the minimal HTTP surface IfconfigCo needs. It is satisfied by
// *http.Client and by retryablehttp's standard client wrapper, and is easy to
// fake in tests.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type IfconfigCo struct {
	client httpDoer
	url    string
	logger *slog.Logger
}

// NewIfGetter returns the production IfconfigCo backed by retryablehttp.
func NewIfGetter() Getter {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	return &IfconfigCo{
		client: retryClient.StandardClient(),
		url:    ipConfigCoAddr,
		logger: slog.Default(),
	}
}

// NewIfconfigCo builds an IfconfigCo with caller-supplied HTTP client and URL.
// Intended for tests; production code should use NewIfGetter.
func NewIfconfigCo(client httpDoer, url string, logger *slog.Logger) *IfconfigCo {
	if logger == nil {
		logger = slog.Default()
	}
	return &IfconfigCo{client: client, url: url, logger: logger}
}

func (i *IfconfigCo) GetIP() (string, error) {
	req, err := http.NewRequest(http.MethodGet, i.url, nil)
	if err != nil {
		i.logger.Error("build request error", "error", err)
		return "", err
	}
	resp, err := i.client.Do(req)
	if err != nil {
		i.logger.Error("get ip error", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ifconfig.co returned status %d", resp.StatusCode)
	}

	s, err := io.ReadAll(resp.Body)
	if err != nil {
		i.logger.Error("read body error", "error", err)
		return "", err
	}
	var ret Response
	if err := json.Unmarshal(s, &ret); err != nil {
		i.logger.Error("unmarshal error", "error", err)
		return "", err
	}
	i.logger.Info("get ip", "ip", ret.IP)
	return ret.IP, nil
}
