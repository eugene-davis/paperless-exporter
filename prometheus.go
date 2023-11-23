package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalDocs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "paperless_total_documents",
		Help: "The total number documents in Paperless NGX",
	})
)

var (
	inboxCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "paperless_inbox_count",
		Help: "The total number of documents in the Paperless NGX inbox",
	})
)

var (
	totalCharCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "paperless_total_char_count",
		Help: "The total number of characters in Paperless NGX",
	})
)

var (
	fileTypeCounts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "paperless_mime_type_count",
		Help: "The total number of docs with a given mime type",
	},
		[]string{
			"mime_type",
		})
)

var (
	fileTasks = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "paperless_file_tasks",
		Help: "The total number of file tasks.",
	},
		[]string{
			"status",
		})
)

// Retrieves the paperless stats and sets the associated prometheus metrics
func setPromStats(stats paperless_stats) {
	totalDocs.Set(float64(stats.TotalDocsCount))
	inboxCount.Set(float64(stats.InboxCount))
	totalCharCount.Set(float64(stats.TotalCharCount))
	for _, fileType := range stats.FileTypeCounts {
		fileTypeCounts.With(prometheus.Labels{"mime_type": fileType.Type}).Set(float64(fileType.Count))
	}
	var fileTasksStatus = make(map[string]int)
	for _, file_task := range stats.FileTaskStats {
		fileTasksStatus[file_task.Status]++
	}
	for taskResult, count := range fileTasksStatus {
		fileTasks.With(prometheus.Labels{"status": taskResult}).Set(float64(count))
	}
}

// Retrieves the paperless stats and sets the associated prometheus metrics in an infinite loop
func setPromStatsLoop(client Client, url string, token string, hostHeader string, refreshSecs int) {
	var failCount int
	for {
		stats := getPaperlessStats(client, url, token, hostHeader)

		slog.Debug("Got the following values from Prometheus", "TotalDocsCounts", stats.TotalCharCount, "InboxCount", stats.InboxCount, "TotalCharCount", stats.TotalCharCount)

		// if a returned value for TotalCharCount is zero, then the request must have failed
		if stats.TotalCharCount != 0 {
			setPromStats(stats)
		} else {
			slog.Debug("Skipping updating due to failure in getting the stats.")
			failCount++
		}

		if failCount > 5 {
			slog.Debug("Too many failures, exiting")
			os.Exit(1)
		}

		time.Sleep(time.Duration(refreshSecs) * time.Second)
	}
}
