package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type mime_type_stats struct {
	Type  string `json:"mime_type"`
	Count int    `json:"mime_type_count"`
}

type file_tasks_stat struct {
	Id              int    `json:"id"`
	TaskId          string `json:"task_id"`
	TaskFileName    string `json:"task_file_name"`
	DateCreated     string `json:"date_created"`
	DateDone        string `json:"date_done"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	Result          string `json:"result"`
	Acknowledged    bool   `json:"acknowledged"`
	RelatedDocument string `json:"related_document"`
}

type paperless_stats struct {
	TotalDocsCount int               `json:"documents_total"`
	InboxCount     int               `json:"documents_inbox"`
	InboxTags      int               `json:"inbox_tag"`
	TotalCharCount int               `json:"character_count"`
	FileTypeCounts []mime_type_stats `json:"document_file_type_counts"`
	FileTaskStats  []file_tasks_stat
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Makes a call to paperless and returns it in the provided results struct
func getPaperlessAPIInfo(client Client, url string, token string, hostHeader string, results any) error {
	req, _ := http.NewRequest("GET", url, nil)
	formatted_token := fmt.Sprintf("Token %s", token)
	req.Header.Set("Authorization", formatted_token)

	if hostHeader != "" {
		req.Header.Set("Host", hostHeader)
	}

	res, err := client.Do(req)

	if err != nil {
		slog.Error("error making http request: %s\n", err)
		return err
	}

	slog.Debug("response:", "status_code", res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to load response %s", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("error response:", "body", string(bodyBytes))
		return errors.New(fmt.Sprintf("Failing status code %d", res.StatusCode))
	}

	bodyString := string(bodyBytes)
	slog.Debug("Response", "body", bodyString)
	err = json.Unmarshal(bodyBytes, results)
	if err != nil {
		slog.Error("Failed to decode response body", "error", err, "body", bodyString)
		return err
	}
	return nil
}

// Gets the paperless stats
func getPaperlessStats(client Client, url string, token string, hostHeader string) paperless_stats {
	var stats paperless_stats
	statusUrl := fmt.Sprintf("%s/api/statistics/", url)
	slog.Debug("making reqest to", "url", statusUrl)
	err := getPaperlessAPIInfo(client, statusUrl, token, hostHeader, &stats)
	if err != nil {
		slog.Error("Failed to get statistics", "error", err)
		return stats
	}

	// Get file tasks info
	var fileTasksStats []file_tasks_stat
	fileTasksUrl := fmt.Sprintf("%s/api/tasks/", url)
	slog.Debug("making reqest to", "url", fileTasksUrl)
	err = getPaperlessAPIInfo(client, fileTasksUrl, token, hostHeader, &fileTasksStats)
	if err != nil {
		slog.Error("Failed to get file tasks", "error", err)
		return stats
	}

	stats.FileTaskStats = fileTasksStats

	return stats
}
