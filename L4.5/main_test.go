package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func Benchmark(b *testing.B) {

	gin.SetMode(gin.TestMode)

	handler := gin.New()
	handler.POST("/add", add)

	body, err := json.Marshal(map[string]int{"first": 2, "second": 2})
	if err != nil {
		b.Fatalf("failed to marshal: %v", err)
	}

	b.ReportAllocs()

	for b.Loop() {
		req := httptest.NewRequest("POST", "/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, req)
	}

}
