package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Gets an environmental variable with a default
func getEnvWithDefault(key string, default_value string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = default_value
	}
	return value
}

func main() {
	// get log level
	log_level := getEnvWithDefault("VERBOSITY", "INFO")
	var handlerOpts slog.HandlerOptions
	switch strings.ToUpper(log_level) {
	case "DEBUG":
		handlerOpts = slog.HandlerOptions{Level: slog.LevelDebug}
	case "INFO":
		handlerOpts = slog.HandlerOptions{Level: slog.LevelInfo}
	case "WARN":
		handlerOpts = slog.HandlerOptions{Level: slog.LevelWarn}
	case "ERROR":
		handlerOpts = slog.HandlerOptions{Level: slog.LevelError}
	default:
		slog.Error("Invalid log level set", "level", log_level)
		handlerOpts = slog.HandlerOptions{Level: slog.LevelInfo}
	}

	// set log level
	handler := slog.NewTextHandler(os.Stderr, &handlerOpts)
	slog.SetDefault(slog.New(handler))

	// get base URL
	requestURL := getEnvWithDefault("URL", "http://localhost:8000")
	slog.Debug("set request url", "url", requestURL)

	// set refresh frequency
	refreshSecs, err := strconv.Atoi(getEnvWithDefault("REFRESH_SECS", "1800"))
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to load REFRESH_SECS %d, defaulting to 1800 seconds", refreshSecs))
		refreshSecs = 1800
	}
	slog.Debug("set refresh seconds", "seconds", refreshSecs)

	// get host header if set
	hostHeader := getEnvWithDefault("HOST_HEADER", "")
	slog.Debug("set host header", "host", hostHeader)

	// get paperless token
	token := os.Getenv("PAPERLESS_TOKEN")
	token_file := os.Getenv("PAPERLESS_TOKEN_FILE")
	if token == "" {
		slog.Debug("PAPERLESS_TOKEN environmental variable not available, checking PAPERLESS_TOKEN_FILE")
		if token_file == "" {
			slog.Debug("PAPERLESS_TOKEN_FILE environmental variable not available, defaulting to ./secrets/token")
			token_file = "./secrets/token"
		}
	}
	token_bytes, err := os.ReadFile(token_file)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to read tokenfile %s", token_file))
	}
	token = string(token_bytes)

	if token == "" {
		slog.Error("No Paperless token available, exiting")
		os.Exit(1)
	}

	port := getEnvWithDefault("METRICS_PORT", "8001")

	client := &http.Client{Timeout: 2 * time.Second}
	go setPromStatsLoop(client, requestURL, token, hostHeader, refreshSecs)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
