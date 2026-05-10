package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type cloudflare struct{}

func init() { register(&cloudflare{}) }

func (c *cloudflare) Name() string  { return "cloudflare" }
func (c *cloudflare) Label() string { return "Cloudflare" }

func (c *cloudflare) ParamDefs() []ParamDef {
	return []ParamDef{
		{Key: "api_token", Label: "API Token", Placeholder: "your Cloudflare API token", Secret: true},
		{Key: "zone_id", Label: "Zone ID", Placeholder: "zone ID from the domain Overview page"},
		{Key: "record_name", Label: "Record Name", Placeholder: "home.example.com"},
	}
}

func (c *cloudflare) Update(ip string, params map[string]string) error {
	token, err := param(params, "api_token")
	if err != nil {
		return err
	}
	zoneID, err := param(params, "zone_id")
	if err != nil {
		return err
	}
	recordName, err := param(params, "record_name")
	if err != nil {
		return err
	}

	recordID, err := c.findRecord(token, zoneID, recordName)
	if err != nil {
		return err
	}

	return c.patchRecord(token, zoneID, recordID, recordName, ip)
}

type cfResponse struct {
	Success bool              `json:"success"`
	Errors  []cfError         `json:"errors"`
	Result  []json.RawMessage `json:"result"`
}

type cfError struct {
	Message string `json:"message"`
}

func (c *cloudflare) do(method, url, token string, body interface{}) ([]byte, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *cloudflare) findRecord(token, zoneID, name string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A&name=%s", zoneID, name)
	data, err := c.do("GET", url, token, nil)
	if err != nil {
		return "", fmt.Errorf("cloudflare list records: %w", err)
	}

	var resp struct {
		Success bool `json:"success"`
		Errors  []cfError
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("cloudflare parse: %w", err)
	}
	if !resp.Success {
		if len(resp.Errors) > 0 {
			return "", fmt.Errorf("cloudflare: %s", resp.Errors[0].Message)
		}
		return "", fmt.Errorf("cloudflare: unknown error")
	}
	if len(resp.Result) == 0 {
		return "", fmt.Errorf("cloudflare: A record %q not found in zone", name)
	}
	return resp.Result[0].ID, nil
}

func (c *cloudflare) patchRecord(token, zoneID, recordID, name, ip string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	payload := map[string]interface{}{
		"type":    "A",
		"name":    name,
		"content": ip,
		"ttl":     1,
		"proxied": false,
	}
	data, err := c.do("PATCH", url, token, payload)
	if err != nil {
		return fmt.Errorf("cloudflare patch: %w", err)
	}

	var resp struct {
		Success bool      `json:"success"`
		Errors  []cfError `json:"errors"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("cloudflare parse: %w", err)
	}
	if !resp.Success {
		if len(resp.Errors) > 0 {
			return fmt.Errorf("cloudflare: %s", resp.Errors[0].Message)
		}
		return fmt.Errorf("cloudflare: unknown error")
	}
	return nil
}
