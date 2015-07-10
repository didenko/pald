package palc

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	paldPort uint16 = 49200
)

type Palc struct {
	svr  string
	port uint16
}

func New(srv string, port uint16) *Palc {
	p := new(Palc)
	p.svr = srv
	p.port = port
	return p
}

func (p *Palc) request(verb, param, value string) (string, error) {

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s?%s=%s", p.svr, p.port, verb, param, value))
	defer resp.Body.Close()

	if resp == nil {
		return "", newFromError(p.svr, p.port, err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", newFromError(p.svr, p.port, err)
	}

	body := string(bodyBytes)

	if err = newFromResp(p.svr, p.port, resp.StatusCode, body); err != nil {
		return "", err
	}

	return body, nil
}

func (p *Palc) Del(port string) error {
	_, err := p.request("del", "port", port)
	return err
}

func (p *Palc) byName(verb string, name string) (uint16, error) {

	body, err := p.request(verb, "service", name)

	if err != nil {
		return 0, err
	}

	var port uint16
	_, err = fmt.Sscan(body, &port)
	return port, newFromError(p.svr, p.port, err)
}

func (p *Palc) Get(service string) (uint16, error) {
	return p.byName("get", service)
}

func (p *Palc) Set(service string) (uint16, error) {
	return p.byName("set", service)
}
