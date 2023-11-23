package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

type test_parameters struct {
	Body       string
	StatusCode int
	Error      error
}

// Configures a mock HTTP client for tests
func getMockClient(test test_parameters) HTTPClientMock {
	var client = &HTTPClientMock{}
	client.DoFunc = func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			// create the custom body
			Body: io.NopCloser(strings.NewReader(test.Body)),
			// create the custom status code
			StatusCode: test.StatusCode,
		}, test.Error
	}
	return *client
}

// Basic test to ensure that a correct value can get loaded from json
func TestGetPaperlessAPIInfoHappyCase(t *testing.T) {
	var test_stats paperless_stats
	testParams := test_parameters{
		Body:       `{"documents_total":3}`,
		StatusCode: http.StatusOK,
		Error:      nil,
	}

	client := getMockClient(testParams)
	// Function under test
	err := getPaperlessAPIInfo(client, "test.url", "token", "header", &test_stats)

	if err != testParams.Error && err.Error() != testParams.Error.Error() {
		t.Fatalf("Expected error to be %v, got %v", testParams.Error, err)
	}

	if test_stats.TotalDocsCount != 3 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 3", test_stats.TotalDocsCount)
	}
}

// Tests results when the response fail
func TestGetPaperlessAPIInfoFailedRequest(t *testing.T) {
	var test_stats paperless_stats
	testParams := test_parameters{
		Body:       ``,
		StatusCode: http.StatusOK,
		Error:      fmt.Errorf(http.StatusText(http.StatusBadRequest)),
	}

	client := getMockClient(testParams)
	// Function under test
	err := getPaperlessAPIInfo(client, "test.url", "token", "header", &test_stats)

	if err != testParams.Error && err.Error() != testParams.Error.Error() {
		t.Fatalf("Expected error to be %v, got %v", testParams.Error, err)
	}

	if test_stats.TotalDocsCount != 0 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 0", test_stats.TotalDocsCount)
	}
}

// Tests the failure on bad json response
func TestGetPaperlessAPIInfoBadJsonResponse(t *testing.T) {
	var test_stats paperless_stats
	testParams := test_parameters{
		Body:       `aoeu`,
		StatusCode: http.StatusOK,
		Error:      nil,
	}

	client := getMockClient(testParams)
	// Function under test
	err := getPaperlessAPIInfo(client, "test.url", "token", "header", &test_stats)

	var badJsonUnmarshalErr *json.InvalidUnmarshalError
	if errors.As(err, &badJsonUnmarshalErr) {
		t.Fatalf("Expected error to be %v, got %v", testParams.Error, err)
	}

	if test_stats.TotalDocsCount != 0 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 0", test_stats.TotalDocsCount)
	}
}

// Basic test to ensure that a correct value can get loaded from json
func TestGetPaperlessStatsHappyCase(t *testing.T) {
	testParams := test_parameters{
		Body:       `{"documents_total":3}`,
		StatusCode: http.StatusOK,
		Error:      nil,
	}

	client := getMockClient(testParams)
	// Function under test
	test_stats := getPaperlessStats(client, "test.url", "token", "header")

	if test_stats.TotalDocsCount != 3 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 3", test_stats.TotalDocsCount)
	}
}

// Basic test to ensure that no response still gets stats returned but set to 0
func TestGetPaperlessStatsNoResponse(t *testing.T) {
	testParams := test_parameters{
		Body:       ``,
		StatusCode: http.StatusOK,
		Error:      nil,
	}

	client := getMockClient(testParams)
	// Function under test
	test_stats := getPaperlessStats(client, "test.url", "token", "header")

	if test_stats.TotalDocsCount != 0 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 3", test_stats.TotalDocsCount)
	}
}

// Basic test to ensure that responses that generate an error still gets stats returned but set to 0
func TestGetPaperlessStatsError(t *testing.T) {
	testParams := test_parameters{
		Body:       ``,
		StatusCode: http.StatusOK,
		Error:      fmt.Errorf(http.StatusText(http.StatusBadRequest)),
	}

	client := getMockClient(testParams)
	// Function under test
	test_stats := getPaperlessStats(client, "test.url", "token", "header")

	if test_stats.TotalDocsCount != 0 {
		t.Fatalf("Parsed %d for TotalDocsCount rather than 3", test_stats.TotalDocsCount)
	}
}
