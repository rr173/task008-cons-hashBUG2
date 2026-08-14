package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestBatchOwnersRejectsEmptyKey verifies that POST /owners rejects
// requests containing empty or whitespace-only keys, consistent with
// how GET /owner validates the key parameter.
func TestBatchOwnersRejectsEmptyKey(t *testing.T) {
	srv := NewService()
	srv.AddNode("node1", 10)
	mux := buildMux(srv)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	cases := []struct {
		name string
		body string
	}{
		{"empty string key", `{"keys":["valid","","also-valid"]}`},
		{"whitespace-only key", `{"keys":["ok","   ","fine"]}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Post(ts.URL+"/owners", "application/json",
				strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				t.Fatalf("batch owners should reject %s but got 200: %s", tc.name, string(body))
			}

			var errResp struct {
				Error string `json:"error"`
			}
			json.Unmarshal(body, &errResp)
			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("want 400, got %d: %s", resp.StatusCode, errResp.Error)
			}
		})
	}
}
