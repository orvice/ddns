package ip

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIfconfigCo_GetIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ip":"1.2.3.4"}`))
	}))
	defer server.Close()

	g := NewIfconfigCo(server.Client(), server.URL, nil)
	got, err := g.GetIP()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "1.2.3.4" {
		t.Fatalf("want 1.2.3.4, got %q", got)
	}
}

func TestIfconfigCo_GetIP_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	g := NewIfconfigCo(server.Client(), server.URL, nil)
	if _, err := g.GetIP(); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestIfconfigCo_GetIP_BadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	g := NewIfconfigCo(server.Client(), server.URL, nil)
	if _, err := g.GetIP(); err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
}
