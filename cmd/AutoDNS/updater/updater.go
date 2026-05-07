package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"autodns/cmd/AutoDNS/config"
	"autodns/cmd/AutoDNS/providers"

	"github.com/OrbitOS-org/sdk-go/v26/logger"
)

const logTag = "updater"

type Updater struct {
	cfg       *config.Manager
	mu        sync.RWMutex
	currentIP string
	forceCh   chan struct{}
}

func New(cfg *config.Manager) *Updater {
	return &Updater{
		cfg:     cfg,
		forceCh: make(chan struct{}, 1),
	}
}

func (u *Updater) CurrentIP() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.currentIP
}

func (u *Updater) ForceUpdate() {
	select {
	case u.forceCh <- struct{}{}:
	default:
	}
}

func (u *Updater) Run(ctx context.Context) {
	u.runCycle()

	cfg := u.cfg.Get()
	interval := time.Duration(cfg.CheckIntervalMin) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-u.forceCh:
			logger.Infof(logTag, "forced update triggered")
			u.runCycle()
			ticker.Reset(time.Duration(u.cfg.Get().CheckIntervalMin) * time.Minute)
		case <-ticker.C:
			newInterval := time.Duration(u.cfg.Get().CheckIntervalMin) * time.Minute
			if newInterval != interval {
				interval = newInterval
				ticker.Reset(interval)
			}
			u.runCycle()
		}
	}
}

func (u *Updater) runCycle() {
	cfg := u.cfg.Get()

	ip, err := detectIP(cfg.IPDetectorURL)
	if err != nil {
		logger.Errorf(logTag, "detect public IP: %v", err)
		return
	}
	logger.Infof(logTag, "public IP: %s", ip)

	u.mu.Lock()
	u.currentIP = ip
	u.mu.Unlock()

	for _, entry := range cfg.Entries {
		if !entry.Enabled {
			continue
		}
		if entry.LastIP == ip {
			logger.Infof(logTag, "entry %q: IP unchanged (%s), skipping", entry.Name, ip)
			continue
		}
		u.updateEntry(entry, ip)
	}
}

func (u *Updater) updateEntry(entry config.Entry, ip string) {
	p, ok := providers.Get(entry.Provider)
	if !ok {
		logger.Errorf(logTag, "entry %q: unknown provider %q", entry.Name, entry.Provider)
		_ = u.cfg.UpdateEntry(entry.ID, func(e *config.Entry) {
			e.LastError = fmt.Sprintf("unknown provider: %s", entry.Provider)
			e.LastStatus = "error"
		})
		return
	}

	logger.Infof(logTag, "updating %q via %s → %s", entry.Name, p.Label(), ip)
	err := p.Update(ip, entry.Params)

	_ = u.cfg.UpdateEntry(entry.ID, func(e *config.Entry) {
		e.LastIP = ip
		e.LastUpdate = time.Now()
		if err != nil {
			e.LastError = err.Error()
			e.LastStatus = "error"
			logger.Errorf(logTag, "entry %q failed: %v", entry.Name, err)
		} else {
			e.LastError = ""
			e.LastStatus = "ok"
			logger.Infof(logTag, "entry %q updated successfully", entry.Name)
		}
	})
}

func detectIP(detectorURL string) (string, error) {
	resp, err := http.Get(detectorURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty response from %s", detectorURL)
	}
	return ip, nil
}
