package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	DeviceHost       string  `json:"device_host"`
	CheckIntervalMin int     `json:"check_interval_min"`
	IPDetectorURL    string  `json:"ip_detector_url"`
	WebPort          int     `json:"web_port"`
	Entries          []Entry `json:"entries"`
}

type Entry struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Provider   string            `json:"provider"`
	Enabled    bool              `json:"enabled"`
	Params     map[string]string `json:"params"`
	LastIP     string            `json:"last_ip,omitempty"`
	LastUpdate time.Time         `json:"last_update,omitempty"`
	LastError  string            `json:"last_error,omitempty"`
	LastStatus string            `json:"last_status,omitempty"`
}

func defaults() *Config {
	return &Config{
		DeviceHost:       "192.168.5.226",
		CheckIntervalMin: 5,
		IPDetectorURL:    "https://api.ipify.org",
		WebPort:          9033,
		Entries:          []Entry{},
	}
}

type Manager struct {
	path string
	mu   sync.RWMutex
	cfg  *Config
}

func NewManager(path string) (*Manager, error) {
	m := &Manager{path: path}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		m.cfg = defaults()
		return m, m.save()
	}
	if err != nil {
		return nil, err
	}
	m.cfg = defaults()
	if err := json.Unmarshal(data, m.cfg); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) Get() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp := *m.cfg
	cp.Entries = make([]Entry, len(m.cfg.Entries))
	copy(cp.Entries, m.cfg.Entries)
	return cp
}

func (m *Manager) Update(fn func(*Config)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	fn(m.cfg)
	return m.save()
}

func (m *Manager) UpdateEntry(id string, fn func(*Entry)) error {
	return m.Update(func(c *Config) {
		for i := range c.Entries {
			if c.Entries[i].ID == id {
				fn(&c.Entries[i])
				return
			}
		}
	})
}

func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0600)
}

func NewID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
