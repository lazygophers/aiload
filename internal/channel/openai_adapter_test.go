package channel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lazygophers/aiload"
)

func TestNewOpenAiAdapter(t *testing.T) {
	channel := &aiload.ModelChannel{}
	adapter := NewOpenAiAdapter(channel)

	if adapter == nil {
		t.Fatal("NewOpenAiAdapter returned nil")
	}
	if adapter.client == nil {
		t.Error("adapter.client is nil")
	}
}

func TestOpenAiAdapter_GetModelList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}

		if r.Header.Get("X-Test-Error") == "true" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: APIError{
					Message: "internal server error",
					Type:    "server_error",
				},
			})
			return
		}

		response := ModelListResponse{
			Object: "list",
			Data: []ModelData{
				{
					Id:      "gpt-4",
					Object:  "model",
					Created: 1677610602,
					OwnedBy: "openai",
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
				if got.Object != "list" {
					t.Errorf("got Object %q, want %q", got.Object, "list")
				}
				if len(got.Data) != 1 {
					t.Fatalf("got %d models, want 1", len(got.Data))
				}
				if got.Data[0].Id != "gpt-4" {
					t.Errorf("got model id %q, want %q", got.Data[0].Id, "gpt-4")
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
			expectErr: errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter := NewOpenAiAdapter(tc.channel)

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
