package ip

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	IpConfigCoAddr = "https://ifconfig.co/json"
)

type IPResponse struct {
	IP string `json:"ip"`
}

type IPGetter interface {
	GetIP() (string, error)
}

type IfconfigCo struct {
}

func NewIfconfigCo() *IfconfigCo {
	return new(IfconfigCo)
}

func (i *IfconfigCo) GetIP() (string, error) {
	resp, err := http.DefaultClient.Get(IpConfigCoAddr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ret IPResponse
	err = json.Unmarshal(s, &ret)
	if err != nil {
		return "", err
	}
	return ret.IP, nil
}
