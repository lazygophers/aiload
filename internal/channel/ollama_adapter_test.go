package channel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lazygophers/aiload"
)

func TestNewOllamaAdapter(t *testing.T) {
	channel := &aiload.ModelChannel{}
	adapter := NewOllamaAdapter(channel)

	if adapter == nil {
		t.Fatal("NewOllamaAdapter returned nil")
	}
	if adapter.client == nil {
		t.Error("adapter.client is nil")
	}
}

func TestOllamaAdapter_GetModelList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			http.NotFound(w, r)
			return
		}

		if r.Header.Get("X-Test-Error") == "true" {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		response := OllamaGetLocalModelListRsp{
			Models: []OllamaModel{
				{
					Name:       "llama2",
					Model:      "llama2:latest",
					ModifiedAt: time.Now(),
					Size:       12345678,
					Digest:     "abcdef123456",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	testCases := []struct {
		name          string
		channel       *aiload.ModelChannel
		req           *GetModelListReq
		wantErr       bool
		expectErr     error
		checkResponse func(t *testing.T, got *GetModelListRsp)
	}{
		{
			name: "Success",
			channel: &aiload.ModelChannel{
				BaseUrl: mockServer.URL,
			},
			req:     &GetModelListReq{},
			wantErr: false,
			checkResponse: func(t *testing.T, got *GetModelListRsp) {
				if got == nil {
					t.Fatal("response should not be nil")
				}
				if len(got.Data) != 1 {
					t.Fatalf("got %d models, want 1", len(got.Data))
				}
				if got.Data[0].Name != "llama2" {
					t.Errorf("got model name %q, want %q", got.Data[0].Name, "llama2")
				}
				if got.Data[0].Size != 12345678 {
					t.Errorf("got model size %d, want %d", got.Data[0].Size, 12345678)
				}
			},
		},
		{
			name: "API Error",
			channel: &aiload.ModelChannel{
				BaseUrl: mockServer.URL,
			},
			req:       &GetModelListReq{},
			wantErr:   true,
			expectErr: errors.New("ollama api error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter := NewOllamaAdapter(tc.channel)

			if tc.name == "API Error" {
				// Temporarily set a header on the global client to trigger the mock error
				client.SetHeader("X-Test-Error", "true")
				defer client.Header.Del("X-Test-Error") // Clean up after the test
			}

			gotRsp, err := adapter.GetModelList(context.Background(), tc.req)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetModelList() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if tc.checkResponse != nil {
					tc.checkResponse(t, gotRsp)
				}
			} else if tc.expectErr != nil {
				if !strings.Contains(err.Error(), tc.expectErr.Error()) {
					t.Errorf("GetModelList() error = %v, want error containing %v", err, tc.expectErr)
				}
			}
		})
	}
}
