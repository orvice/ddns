package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/weeon/log"
)

type Mu struct {
	apiUrl string
	nodeID int
}

func NewMu(apiUrl string, nodeID int) (*Mu, error) {
	return &Mu{
		apiUrl: apiUrl,
		nodeID: nodeID,
	}, nil
}

func (m *Mu) GetIP(ctx context.Context, domain string) (string, error) {
	return "", nil
}
func (m *Mu) UpdateIP(tx context.Context, domain, ip string) error {
	uri := fmt.Sprintf("%s/nodes/%d/ip", m.apiUrl, m.nodeID)
	ma := map[string]interface{}{
		"ip":     ip,
		"domain": domain,
	}

	body, err := json.Marshal(ma)
	if err != nil {
		return err
	}
	input := bytes.NewBuffer(body)

	resp, err := http.DefaultClient.Post(uri, "application/json", input)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	log.Debugw("update ip ",
		"response", string(b),
	)

	return nil
}
