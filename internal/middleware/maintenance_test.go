package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testConfig struct {
	enabled bool
}

func (t *testConfig) MaintenanceMode() bool {
	return t.enabled
}

func TestHandle_Enabled(t *testing.T) {
	t.Parallel()

	responder := ProcessMaintenance(&testConfig{true})

	r := &http.Request{}
	w := httptest.NewRecorder()

	responder(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("handler was invoked")
	})).ServeHTTP(w, r)

	w.Flush()

	if got, want := w.Code, 429; got != want {
		t.Errorf("expected %d to be %d", got, want)
	}
}

func TestHandle_Disabled(t *testing.T) {
	t.Parallel()

	responder := ProcessMaintenance(&testConfig{false})

	r := &http.Request{}
	w := httptest.NewRecorder()

	responder(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w, r)

	w.Flush()

	if got, want := w.Code, http.StatusOK; got != want {
		t.Errorf("expected %d to be %d", got, want)
	}
}
