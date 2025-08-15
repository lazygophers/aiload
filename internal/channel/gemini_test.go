package channel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGemini(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	g := NewGemini(mockChannel)

	require.NotNil(t, g)
	assert.Equal(t, mockChannel, g.channel)
}

func TestGemini_GetRequest(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	g := NewGemini(mockChannel)
	req := g.GetRequest()

	require.NotNil(t, req)
}

func TestGemini_GetModelList(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewGemini(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name          string
		giveReq       *GeminiGetModelListReq
		setupMock     func(t *testing.T, req *GeminiGetModelListReq)
		checkResponse func(t *testing.T, rsp *GeminiGetModelListRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiGetModelListReq{
				PageSize:  10,
				PageToken: "token-123",
			},
			setupMock: func(t *testing.T, req *GeminiGetModelListReq) {
				mockResponse := GeminiGetModelListRsp{
					Models: []GeminiModel{
						{
							Name:        "models/gemini-pro",
							DisplayName: "Gemini Pro",
						},
					},
					NextPageToken: "token-456",
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				httpmock.RegisterResponder("GET", "https://generativelanguage.googleapis.com/v1beta/models",
					func(r *http.Request) (*http.Response, error) {
						assert.Equal(t, "10", r.URL.Query().Get("pageSize"))
						assert.Equal(t, "token-123", r.URL.Query().Get("pageToken"))
						return httpmock.NewStringResponse(200, string(respBody)), nil
					},
				)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				require.Len(t, rsp.Models, 1)
				assert.Equal(t, "models/gemini-pro", rsp.Models[0].Name)
				assert.Equal(t, "token-456", rsp.NextPageToken)
			},
		},
		{
			name: "page size is zero",
			giveReq: &GeminiGetModelListReq{
				PageToken: "token-abc",
			},
			setupMock: func(t *testing.T, req *GeminiGetModelListReq) {
				httpmock.RegisterResponder("GET", "https://generativelanguage.googleapis.com/v1beta/models",
					func(r *http.Request) (*http.Response, error) {
						assert.Equal(t, "50", r.URL.Query().Get("pageSize")) // Default page size
						assert.Equal(t, "token-abc", r.URL.Query().Get("pageToken"))
						return httpmock.NewJsonResponse(200, GeminiGetModelListRsp{})
					},
				)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
			},
		},
		{
			name: "page size exceeds max",
			giveReq: &GeminiGetModelListReq{
				PageSize: 2000,
			},
			setupMock: func(t *testing.T, req *GeminiGetModelListReq) {
				httpmock.RegisterResponder("GET", "https://generativelanguage.googleapis.com/v1beta/models",
					func(r *http.Request) (*http.Response, error) {
						assert.Equal(t, "1000", r.URL.Query().Get("pageSize")) // Max page size
						return httpmock.NewJsonResponse(200, GeminiGetModelListRsp{})
					},
				)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
			},
		},
		{
			name:    "api error",
			giveReq: &GeminiGetModelListReq{},
			setupMock: func(t *testing.T, req *GeminiGetModelListReq) {
				httpmock.RegisterResponder("GET", "https://generativelanguage.googleapis.com/v1beta/models",
					httpmock.NewStringResponder(500, "Internal Server Error"),
				)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
			},
		},
		{
			name:    "network error",
			giveReq: &GeminiGetModelListReq{},
			setupMock: func(t *testing.T, req *GeminiGetModelListReq) {
				httpmock.RegisterResponder("GET", "https://generativelanguage.googleapis.com/v1beta/models",
					httpmock.NewErrorResponder(errors.New("connection failed")),
				)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "connection failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()
			tt.setupMock(t, tt.giveReq)

			rsp, err := p.GetModelList(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}