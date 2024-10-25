package ip

import (
	"encoding/json"
	"io"
	"log/slog"

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

type IfconfigCo struct {
}

func NewIfconfigCo() *IfconfigCo {
	return new(IfconfigCo)
}

func (i *IfconfigCo) GetIP() (string, error) {
	logger := slog.Default()
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	resp, err := retryClient.Get(ipConfigCoAddr)
	if err != nil {
		logger.Error("get ip error", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	s, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error", "error", err)
		return "", err
	}
	var ret Response
	err = json.Unmarshal(s, &ret)
	if err != nil {
		logger.Error("unmarshal error", "error", err)
		return "", err
	}
	logger.Info("get ip", "ip", ret.IP)
	return ret.IP, nil
}
