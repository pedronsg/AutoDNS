package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"autodns/cmd/AutoDNS/config"
	"autodns/cmd/AutoDNS/providers"
	"autodns/cmd/AutoDNS/updater"

	"github.com/OrbitOS-org/sdk-go/v26/logger"
)

const logTag = "web"

type Server struct {
	cfg     *config.Manager
	updater *updater.Updater
	version string
	mux     *http.ServeMux
}

func NewServer(cfg *config.Manager, u *updater.Updater, version string) *Server {
	s := &Server{cfg: cfg, updater: u, version: version, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/status", s.handleStatus)
	s.mux.HandleFunc("/api/providers", s.handleProviders)
	s.mux.HandleFunc("/api/config", s.handleConfig)
	s.mux.HandleFunc("/api/entries", s.handleEntries)
	s.mux.HandleFunc("/api/entries/toggle", s.handleToggle)
	s.mux.HandleFunc("/api/update", s.handleUpdate)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(strings.ReplaceAll(indexHTML, "{{VERSION}}", "v"+s.version)))
}

type statusResponse struct {
	CurrentIP   string         `json:"current_ip"`
	CheckedAt   time.Time      `json:"checked_at"`
	Entries     []config.Entry `json:"entries"`
	IntervalMin int            `json:"interval_min"`
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	cfg := s.cfg.Get()
	writeJSON(w, statusResponse{
		CurrentIP:   s.updater.CurrentIP(),
		CheckedAt:   time.Now(),
		Entries:     cfg.Entries,
		IntervalMin: cfg.CheckIntervalMin,
	})
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	type providerInfo struct {
		Name      string            `json:"name"`
		Label     string            `json:"label"`
		ParamDefs []providers.ParamDef `json:"param_defs"`
	}
	list := providers.List()
	out := make([]providerInfo, len(list))
	for i, p := range list {
		out[i] = providerInfo{Name: p.Name(), Label: p.Label(), ParamDefs: p.ParamDefs()}
	}
	writeJSON(w, out)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg := s.cfg.Get()
		writeJSON(w, map[string]interface{}{
			"check_interval_min": cfg.CheckIntervalMin,
			"ip_detector_url":    cfg.IPDetectorURL,
			"web_port":           cfg.WebPort,
		})
	case http.MethodPost:
		var body struct {
			CheckIntervalMin int    `json:"check_interval_min"`
			IPDetectorURL    string `json:"ip_detector_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := s.cfg.Update(func(c *config.Config) {
			if body.CheckIntervalMin > 0 {
				c.CheckIntervalMin = body.CheckIntervalMin
			}
			if body.IPDetectorURL != "" {
				c.IPDetectorURL = body.IPDetectorURL
			}
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body struct {
			Name     string            `json:"name"`
			Provider string            `json:"provider"`
			Params   map[string]string `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Provider == "" {
			writeError(w, "name and provider are required", http.StatusBadRequest)
			return
		}
		if _, ok := providers.Get(body.Provider); !ok {
			writeError(w, "unknown provider: "+body.Provider, http.StatusBadRequest)
			return
		}
		entry := config.Entry{
			ID:       config.NewID(),
			Name:     body.Name,
			Provider: body.Provider,
			Enabled:  true,
			Params:   body.Params,
		}
		if err := s.cfg.Update(func(c *config.Config) {
			c.Entries = append(c.Entries, entry)
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Infof(logTag, "added entry %q (%s)", entry.Name, entry.Provider)
		writeJSON(w, entry)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "id required", http.StatusBadRequest)
			return
		}
		if err := s.cfg.Update(func(c *config.Config) {
			filtered := c.Entries[:0]
			for _, e := range c.Entries {
				if e.ID != id {
					filtered = append(filtered, e)
				}
			}
			c.Entries = filtered
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Infof(logTag, "deleted entry %s", id)
		writeJSON(w, map[string]string{"status": "ok"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "id required", http.StatusBadRequest)
		return
	}
	var enabled bool
	if err := s.cfg.UpdateEntry(id, func(e *config.Entry) {
		e.Enabled = !e.Enabled
		enabled = e.Enabled
	}); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]bool{"enabled": enabled})
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.updater.ForceUpdate()
	writeJSON(w, map[string]string{"status": "update triggered"})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
