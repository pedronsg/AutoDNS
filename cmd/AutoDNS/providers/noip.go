package providers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type noIP struct{}

func init() { register(&noIP{}) }

func (n *noIP) Name() string  { return "noip" }
func (n *noIP) Label() string { return "No-IP" }

func (n *noIP) ParamDefs() []ParamDef {
	return []ParamDef{
		{Key: "username", Label: "Username / Email", Placeholder: "user@example.com"},
		{Key: "password", Label: "Password", Placeholder: "your password", Secret: true},
		{Key: "hostname", Label: "Hostname", Placeholder: "myhome.ddns.net"},
	}
}

func (n *noIP) Update(ip string, params map[string]string) error {
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

	u := fmt.Sprintf("https://dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s",
		url.QueryEscape(hostname), url.QueryEscape(ip))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("User-Agent", "AutoDNS/1.0 admin@example.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("noip request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := strings.TrimSpace(string(body))

	switch {
	case strings.HasPrefix(result, "good"), strings.HasPrefix(result, "nochg"):
		return nil
	case result == "nohost":
		return fmt.Errorf("noip: hostname not found")
	case result == "badauth":
		return fmt.Errorf("noip: invalid credentials")
	case result == "abuse":
		return fmt.Errorf("noip: account blocked for abuse")
	default:
		return fmt.Errorf("noip: %s", result)
	}
}
