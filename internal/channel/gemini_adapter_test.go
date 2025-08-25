package channel

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lazygophers/aiload"
)

func TestNewGeminiAdapter(t *testing.T) {
	channel := &aiload.ModelChannel{}
	adapter := NewGeminiAdapter(channel)

	if adapter == nil {
		t.Fatal("NewGeminiAdapter returned nil")
	}
}

func TestGeminiAdapter_GetModelList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/models" {
			http.NotFound(w, r)
			return
		}

		// Simulate an API error
		if r.Header.Get("X-Test-Error") == "true" {
			http.Error(w, `{"error": {"message": "internal server error"}}`, http.StatusInternalServerError)
			return
		}

		response := GeminiGetModelListRsp{
			Models: []GeminiModel{
				{
					Name:        "gemini-pro",
					DisplayName: "Gemini Pro",
					Version:     "v1",
				},
			},
			NextPageToken: "next-page-token",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	testCases := []struct {
		name          string
		channel       *aiload.ModelChannel
		wantErr       bool
		expectErr     error
		checkResponse func(t *testing.T, got *GetModelListRsp)
	}{
		{
			name: "Success",
			channel: &aiload.ModelChannel{
				BaseUrl: mockServer.URL,
			},
			wantErr: false,
			checkResponse: func(t *testing.T, got *GetModelListRsp) {
				if got == nil {
					t.Fatal("response should not be nil")
				}
				if got.NextPageToken != "next-page-token" {
					t.Errorf("got NextPageToken %q, want %q", got.NextPageToken, "next-page-token")
				}
				if len(got.Data) != 1 {
					t.Fatalf("got %d models, want 1", len(got.Data))
				}
				if got.Data[0].Name != "gemini-pro" {
					t.Errorf("got model name %q, want %q", got.Data[0].Name, "gemini-pro")
				}
			},
		},
		{
			name: "API Error",
			channel: &aiload.ModelChannel{
				BaseUrl: mockServer.URL,
			},
			wantErr:   true,
			expectErr: errors.New("Gemini api error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter := NewGeminiAdapter(tc.channel)

			if tc.name == "API Error" {
				// For the error case, we configure the mock server to respond with an error.
				// The test client will have a header that triggers this behavior.
				client.SetHeader("X-Test-Error", "true")
				defer client.SetHeader("X-Test-Error", "")
			}

			gotRsp, err := adapter.GetModelList()

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
