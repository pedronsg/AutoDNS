package main

import (
	_ "embed"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"autodns/cmd/AutoDNS/config"
	_ "autodns/cmd/AutoDNS/providers" // register all providers via init()
	"autodns/cmd/AutoDNS/updater"
	"autodns/cmd/AutoDNS/web"

	eventsvcv26 "github.com/OrbitOS-org/sdk-go/v26/api/event_service/v26"
	"github.com/OrbitOS-org/sdk-go/v26/client"
	"github.com/OrbitOS-org/sdk-go/v26/logger"
	"github.com/OrbitOS-org/sdk-go/v26/metadata"
)

const logTag = "main"

//go:embed metadata.json
var metadataJSON []byte

var appManifest = metadata.MustParseAppManifestJSON(metadataJSON)

func main() {
	configPath := flag.String("config", resolveDataPath("config.json"), "Config file path")
	hostOverride := flag.String("host", "", "OrbitOS device IP (overrides config)")
	flag.Parse()

	meta := metadata.Build(appManifest)
	logger.Init(meta.Name, "INFO", true)
	logger.Infof(logTag, "Starting %s %s", meta.Name, meta.Version)

	cfg, err := config.NewManager(*configPath)
	if err != nil {
		logger.Fatalf(logTag, "load config: %v", err)
		os.Exit(1)
	}
	if *hostOverride != "" {
		_ = cfg.Update(func(c *config.Config) { c.DeviceHost = *hostOverride })
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	u := updater.New(cfg)
	go u.Run(ctx)

	appCfg := cfg.Get()
	webAddr := fmt.Sprintf("127.0.0.1:%d", appCfg.WebPort)
	srv := web.NewServer(cfg, u, meta.Version)

	httpServer := &http.Server{Addr: webAddr, Handler: srv}
	go func() {
		logger.Infof(logTag, "web UI → http://%s", webAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf(logTag, "web server: %v", err)
		}
	}()

	go connectDevice(ctx, cfg, u, webAddr)

	<-ctx.Done()
	logger.Infof(logTag, "shutting down")
	_ = httpServer.Shutdown(context.Background())
}

func connectDevice(ctx context.Context, cfg *config.Manager, u *updater.Updater, webAddr string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		host := cfg.Get().DeviceHost
		c, err := client.NewClientAuto(host)
		if err != nil {
			logger.Warnf(logTag, "device %s unreachable: %v — retrying in 30s", host, err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(30 * time.Second):
				continue
			}
		}

		hwModel, _ := c.SystemManager.GetHardwareModel()
		logger.Infof(logTag, "connected to device %s (%s)", host, hwModel)

		if err := c.AppHubManager.RegisterWebUI(webAddr, "/autodns"); err != nil {
			logger.Warnf(logTag, "AppHub register: %v", err)
		} else {
			logger.Infof(logTag, "registered with AppHub at /autodns")
		}

		err = c.EventManager.Subscribe(ctx, func(_ *eventsvcv26.Event) {
			logger.Infof(logTag, "network up — triggering DNS update")
			u.ForceUpdate()
		}, client.EVENT_NET_UP)

		_ = c.AppHubManager.UnregisterService()
		c.Close()

		if ctx.Err() != nil {
			return
		}
		logger.Warnf(logTag, "device stream ended (%v) — reconnecting in 30s", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
		}
	}
}

// resolveDataPath returns the path to a file inside orb/data, following the
// OrbitOS convention: gravity sets the process CWD to app.DataDir (orb/data),
// so "orb/data" relative to CWD works on device. When running from the repo
// the CWD is typically the project root, so "cmd/AutoDNS/orb/data" is tried first.
func resolveDataPath(file string) string {
	repoBased := filepath.Clean("cmd/AutoDNS/orb/data")
	if st, err := os.Stat(repoBased); err == nil && st.IsDir() {
		return filepath.Join(repoBased, file)
	}
	localBased := filepath.Clean("orb/data")
	if st, err := os.Stat(localBased); err == nil && st.IsDir() {
		return filepath.Join(localBased, file)
	}
	return filepath.Join(repoBased, file)
}
