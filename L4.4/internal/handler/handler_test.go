package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsEndpoint(t *testing.T) {

	handler := NewHandler()

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("unexpected content-type: %s", contentType)
	}

	body := recorder.Body.String()

	expectedMetrics := []string{
		"go_goroutines",
		"go_memstats_alloc_bytes",
		"go_memstats_num_gc",
		"go_gc_percent",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("expected metric %s not found", metric)
		}
	}

}

func TestPprofEndpoints(t *testing.T) {

	handler := NewHandler()

	endpoints := []string{
		"/debug/pprof/",
		"/debug/pprof/heap",
		"/debug/pprof/goroutine",
	}

	for _, ep := range endpoints {
		request := httptest.NewRequest(http.MethodGet, ep, nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Errorf("endpoint %s returned status %d", ep, recorder.Code)
		}
	}

}
