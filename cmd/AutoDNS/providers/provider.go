package providers

import "fmt"

type Provider interface {
	Name() string
	Label() string
	Update(ip string, params map[string]string) error
	ParamDefs() []ParamDef
}

type ParamDef struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Placeholder string `json:"placeholder"`
	Secret      bool   `json:"secret"`
}

var registry = map[string]Provider{}

func register(p Provider) {
	registry[p.Name()] = p
}

func Get(name string) (Provider, bool) {
	p, ok := registry[name]
	return p, ok
}

func List() []Provider {
	out := make([]Provider, 0, len(registry))
	for _, p := range registry {
		out = append(out, p)
	}
	return out
}

func param(params map[string]string, key string) (string, error) {
	v, ok := params[key]
	if !ok || v == "" {
		return "", fmt.Errorf("missing parameter: %s", key)
	}
	return v, nil
}
