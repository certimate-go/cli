package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// apiEnvelope mirrors the certimate resp.Ok/resp.Err wrapper:
// {"code":0,"msg":"success","data":{...}}.
type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// newTestClient builds a Client pointed at the given test server.
func newTestClient(t *testing.T, server, token string) *Client {
	t.Helper()
	return &Client{
		httpClient: &http.Client{},
		server:     server,
		token:      token,
	}
}

func TestDownloadCertificate_Success(t *testing.T) {
	// Recognizable bytes standing in for a ZIP archive.
	wantZip := []byte("PK\x03\x04 fake zip contents \xde\xad\xbe\xef")

	var gotPath, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)

		// certimate wraps every custom-endpoint response in resp.Ok.
		// zipBytes is a []byte, so encoding/json base64-encodes it.
		type dataShape struct {
			ZipBytes []byte `json:"zipBytes"`
		}
		resp := struct {
			Code int       `json:"code"`
			Msg  string    `json:"msg"`
			Data dataShape `json:"data"`
		}{
			Code: 0,
			Msg:  "success",
			Data: dataShape{ZipBytes: wantZip},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := newTestClient(t, ts.URL, "test-token")

	got, err := client.DownloadCertificate(context.Background(), "cert-123", "PEM")
	if err != nil {
		t.Fatalf("DownloadCertificate returned unexpected error: %v", err)
	}

	if gotPath != "/api/certificates/cert-123/download" {
		t.Errorf("request path = %q, want /api/certificates/cert-123/download", gotPath)
	}

	// The server reads "fileFormat", not "format".
	if !strings.Contains(gotBody, `"fileFormat":"PEM"`) {
		t.Errorf("request body = %q, want it to contain \"fileFormat\":\"PEM\"", gotBody)
	}
	if strings.Contains(gotBody, `"format":"PEM"`) {
		t.Errorf("request body = %q, should not send the old \"format\" field", gotBody)
	}

	if string(got) != string(wantZip) {
		t.Errorf("downloaded bytes = %q, want %q", got, wantZip)
	}
}

func TestDownloadCertificate_ServerErrorWithHTTP200(t *testing.T) {
	// certimate returns HTTP 200 with a non-zero code even when the
	// request fails (resp.Err), so the client must inspect the envelope.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := apiEnvelope{
			Code: 404,
			Msg:  "the certificate was not found",
			Data: nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := newTestClient(t, ts.URL, "test-token")

	got, err := client.DownloadCertificate(context.Background(), "missing-id", "PEM")
	if err == nil {
		t.Fatalf("DownloadCertificate expected an error for code!=0, got nil; bytes=%q", got)
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error %q should surface the server message", err)
	}
	if len(got) != 0 {
		t.Errorf("on error, returned bytes should be empty, got %q", got)
	}
}
