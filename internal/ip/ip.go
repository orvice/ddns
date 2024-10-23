package ip

import (
	"encoding/json"
	"io"
	"net/http"
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
	cli := http.Client{}
	defer cli.CloseIdleConnections()
	resp, err := cli.Get(ipConfigCoAddr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ret Response
	err = json.Unmarshal(s, &ret)
	if err != nil {
		return "", err
	}
	return ret.IP, nil
}
