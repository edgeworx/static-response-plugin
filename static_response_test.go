package static_response_plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	staticresponse "github.com/edgeworx/static-response-plugin"
)

func TestStaticResponse(t *testing.T) {
	ctx := context.Background()

	// Create an empty configuration
	cfg := staticresponse.CreateConfig()

	// Create a no-op handler
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	// Ensure that the plugin fails to initialize with an empty configuration
	handler, err := staticresponse.New(ctx, next, cfg, "staticresponse-plugin")
	if err == nil {
		t.Fatal("expected error for empty configuration, got nil")
	}

	// Add path configurations
	cfg.Paths = append(cfg.Paths, []staticresponse.Path{
		{
			Path:    "/",
			Content: "Hello World!",
		},
		{
			PathRegex: "^/regex/(.*)",
			Content:   "Hello Regex!",
		},
		{
			Path:    "/template",
			Content: `{{define "T"}}Hello {{.}}!{{end}}{{template "T" "Template"}}`,
		},
	}...)

	// Ensure that the plugin initializes successfully
	handler, err = staticresponse.New(ctx, next, cfg, "staticresponse-plugin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tt := []struct {
		name             string
		req              *http.Request
		expectedStatus   int
		expectedResponse string
	}{
		// A simple configured path should return the configured content
		{
			name:             "root",
			req:              httptest.NewRequest("GET", "/", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "Hello World!",
		},
		// A non-configured path should move to the next handler
		{
			name:             "not found",
			req:              httptest.NewRequest("GET", "/not-found", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		// A configured path with a regex should return the configured content
		{
			name:             "regex path foo",
			req:              httptest.NewRequest("GET", "/regex/foo", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "Hello Regex!",
		},
		{
			name:             "regex path bar",
			req:              httptest.NewRequest("GET", "/regex/bar", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "Hello Regex!",
		},
		{
			name:             "regex path baz",
			req:              httptest.NewRequest("GET", "/regex/baz", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "Hello Regex!",
		},
		// A configured path with a template should return the configured content
		{
			name:             "template",
			req:              httptest.NewRequest("GET", "/template", nil),
			expectedStatus:   http.StatusOK,
			expectedResponse: "Hello Template!",
		},
		// TODO: Test templates with request parameters and paths with JSON data
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			handler.ServeHTTP(rw, tc.req)
			if rw.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expectrecorder.Result().Sed status %d, got %d", tc.expectedStatus, rw.Result().StatusCode)
			}
			if rw.Body.String() != tc.expectedResponse {
				t.Errorf("expected response %s, got %s", tc.expectedResponse, rw.Body.String())
			}
		})
	}
}
