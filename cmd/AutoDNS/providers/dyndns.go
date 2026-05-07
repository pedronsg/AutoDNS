package providers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type dynDNS struct{}

func init() { register(&dynDNS{}) }

func (d *dynDNS) Name() string  { return "dyndns" }
func (d *dynDNS) Label() string { return "DynDNS" }

func (d *dynDNS) ParamDefs() []ParamDef {
	return []ParamDef{
		{Key: "username", Label: "Username", Placeholder: "your username"},
		{Key: "password", Label: "Password", Placeholder: "your password", Secret: true},
		{Key: "hostname", Label: "Hostname", Placeholder: "myhome.dyndns.org"},
	}
}

func (d *dynDNS) Update(ip string, params map[string]string) error {
	username, err := param(params, "username")
	if err != nil {
		return err
	}
	password, err := param(params, "password")
	if err != nil {
		return err
	}
	hostname, err := param(params, "hostname")
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://members.dyndns.org/nic/update?hostname=%s&myip=%s",
		url.QueryEscape(hostname), url.QueryEscape(ip))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("User-Agent", "AutoDNS/1.0 admin@example.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("dyndns request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := strings.TrimSpace(string(body))

	switch {
	case strings.HasPrefix(result, "good"), strings.HasPrefix(result, "nochg"):
		return nil
	case result == "nohost":
		return fmt.Errorf("dyndns: hostname not found")
	case result == "badauth":
		return fmt.Errorf("dyndns: invalid credentials")
	case result == "abuse":
		return fmt.Errorf("dyndns: account blocked for abuse")
	default:
		return fmt.Errorf("dyndns: %s", result)
	}
}
