package providers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type duckDNS struct{}

func init() { register(&duckDNS{}) }

func (d *duckDNS) Name() string  { return "duckdns" }
func (d *duckDNS) Label() string { return "DuckDNS" }

func (d *duckDNS) ParamDefs() []ParamDef {
	return []ParamDef{
		{Key: "token", Label: "Token", Placeholder: "your-duckdns-token", Secret: true},
		{Key: "domains", Label: "Domains", Placeholder: "home,office (comma-separated, without .duckdns.org)"},
	}
}

func (d *duckDNS) Update(ip string, params map[string]string) error {
	token, err := param(params, "token")
	if err != nil {
		return err
	}
	domains, err := param(params, "domains")
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://www.duckdns.org/update?domains=%s&token=%s&ip=%s&verbose=true",
		url.QueryEscape(domains), url.QueryEscape(token), url.QueryEscape(ip))

	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("duckdns request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	first := strings.SplitN(strings.TrimSpace(string(body)), "\n", 2)[0]

	if first != "OK" {
		return fmt.Errorf("duckdns: %s", first)
	}
	return nil
}
