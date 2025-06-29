package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CoreKitMDK/corekit-service-core/v2/pkg/core"
	"github.com/CoreKitMDK/corekit-service-logger/v2/pkg/logger"
)

func TestCore(t *testing.T) {
	Core, err := core.NewCore()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(2 * time.Second)

	Core.Logger.Log(logger.INFO, "Test message")

	time.Sleep(2 * time.Second)

	Core.Metrics.Record("test", 1)

	time.Sleep(2 * time.Second)

	Core.Events.Emit("TEST", "Test data")

	time.Sleep(4 * time.Second)
}

func TestTracingConfiguration(t *testing.T) {
	Core, err := core.NewCore()
	if err != nil {
		t.Error(err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := Core.Tracing.TraceHttpRequest(r).Start()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
		tracer.TraceHttpResponseWriter(w).End()
	}))
	defer ts.Close()

	// Make test request
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	time.Sleep(2 * time.Second)
}
