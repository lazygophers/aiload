package channel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGemini_GetModelList(t *testing.T) {
	tests := []struct {
		name          string
		giveReq       *GeminiGetModelListReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiGetModelListRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiGetModelListReq{
				PageSize:  10,
				PageToken: "token-123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "10", r.URL.Query().Get("pageSize"))
				assert.Equal(t, "token-123", r.URL.Query().Get("pageToken"))

				mockResponse := GeminiGetModelListRsp{
					Models: []GeminiModel{
						{
							Name:        "models/gemini-pro",
							DisplayName: "Gemini Pro",
						},
					},
					NextPageToken: "token-456",
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
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
			name:    "api error",
			giveReq: &GeminiGetModelListReq{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
			},
		},
		// The "network error" case is harder to simulate with httptest,
		// as it requires causing a transport-level error.
		// For this exercise, we'll focus on API and HTTP errors.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			mockChannel := &aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			}
			p := NewGemini(mockChannel)

			rsp, err := p.GetModelList(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestGemini_GetModel(t *testing.T) {
	tests := []struct {
		name          string
		giveReq       *GeminiGetModelReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiGetModelRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiGetModelReq{
				Model: "gemini-pro",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, "/v1/models/gemini-pro")

				mockResponse := GeminiModel{
					Name:        "models/gemini-pro",
					DisplayName: "Gemini Pro",
					Description: "The best model for text generation.",
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "models/gemini-pro", rsp.Name)
				assert.Equal(t, "Gemini Pro", rsp.DisplayName)
			},
		},
		{
			name: "api_error",
			giveReq: &GeminiGetModelReq{
				Model: "gemini-pro-non-existent",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": {"message": "Model not found"}}`))
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Model not found")
			},
		},
		{
			name: "http_error",
			giveReq: &GeminiGetModelReq{
				Model: "gemini-pro",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			checkResponse: func(t *testing.T, rsp *GeminiGetModelRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			rsp, err := p.GetModel(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestGemini_EmbedContent(t *testing.T) {
	tests := []struct {
		name          string
		giveReq       *GeminiEmbedContentReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiEmbedContentRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiEmbedContentReq{
				Model: "embedding-001",
				Content: &GeminiContent{
					Parts: GeminiCommonPart{Text: "Hello, world!"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, ":embedContent")
				mockResponse := GeminiEmbedContentRsp{
					Embedding: ContentEmbedding{
						Values: []float32{0.1, 0.2, 0.3},
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
			},
			checkResponse: func(t *testing.T, rsp *GeminiEmbedContentRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, []float32{0.1, 0.2, 0.3}, rsp.Embedding.Values)
			},
		},
		{
			name: "api_error",
			giveReq: &GeminiEmbedContentReq{
				Model: "embedding-001",
				Content: &GeminiContent{
					Parts: GeminiCommonPart{Text: "Invalid content"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": {"message": "Invalid content"}}`))
			},
			checkResponse: func(t *testing.T, rsp *GeminiEmbedContentRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Invalid content")
			},
		},
		{
			name: "http_error",
			giveReq: &GeminiEmbedContentReq{
				Model: "embedding-001",
				Content: &GeminiContent{
					Parts: GeminiCommonPart{Text: "Hello, world!"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Service Unavailable"))
			},
			checkResponse: func(t *testing.T, rsp *GeminiEmbedContentRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			rsp, err := p.EmbedContent(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}
func TestGemini_BatchEmbedContents(t *testing.T) {
	tests := []struct {
		name          string
		giveReq       *GeminiBatchEmbedContentsReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiBatchEmbedContentsRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiBatchEmbedContentsReq{
				Requests: []*GeminiEmbedContentReq{
					{
						Model: "embedding-001",
						Content: &GeminiContent{
							Parts: GeminiCommonPart{Text: "Hello"},
						},
					},
					{
						Model: "embedding-001",
						Content: &GeminiContent{
							Parts: GeminiCommonPart{Text: "World"},
						},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, ":batchEmbedContents")
				mockResponse := &GeminiBatchEmbedContentsRsp{
					Embeddings: []ContentEmbedding{
						{Values: []float32{0.1, 0.2}},
						{Values: []float32{0.3, 0.4}},
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
			},
			checkResponse: func(t *testing.T, rsp *GeminiBatchEmbedContentsRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				require.Len(t, rsp.Embeddings, 2)
				assert.Equal(t, []float32{0.1, 0.2}, rsp.Embeddings[0].Values)
				assert.Equal(t, []float32{0.3, 0.4}, rsp.Embeddings[1].Values)
			},
		},
		{
			name: "api_error",
			giveReq: &GeminiBatchEmbedContentsReq{
				Requests: []*GeminiEmbedContentReq{
					{Model: "embedding-001", Content: &GeminiContent{Parts: GeminiCommonPart{Text: "test"}}},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": {"message": "Invalid request"}}`))
			},
			checkResponse: func(t *testing.T, rsp *GeminiBatchEmbedContentsRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Invalid request")
			},
		},
		{
			name: "http_error",
			giveReq: &GeminiBatchEmbedContentsReq{
				Requests: []*GeminiEmbedContentReq{
					{Model: "embedding-001", Content: &GeminiContent{Parts: GeminiCommonPart{Text: "test"}}},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			checkResponse: func(t *testing.T, rsp *GeminiBatchEmbedContentsRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			rsp, err := p.BatchEmbedContents(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestGemini_GenerateContent(t *testing.T) {
	tests := []struct {
		name          string
		model         string
		giveReq       *GeminiGenerateContentReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiGenerateContentRsp, err error)
	}{
		{
			name:  "success",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "Hello"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, ":generateContent")
				mockResponse := &GeminiGenerateContentRsp{
					Candidates: []*Candidate{
						{
							Content: &GeminiContent{
								Parts: GeminiCommonPart{Text: "Hi there!"},
								Role:  "model",
							},
							FinishReason: "STOP",
						},
					},
				}
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(mockResponse)
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGenerateContentRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				require.Len(t, rsp.Candidates, 1)
				assert.Equal(t, "Hi there!", rsp.Candidates[0].Content.Parts.Text)
			},
		},
		{
			name:  "api_error_with_safety_ratings",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "A sensitive question"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				mockResponse := &GeminiGenerateContentRsp{
					Candidates: []*Candidate{},
					PromptFeedback: &PromptFeedback{
						BlockReason: "SAFETY",
						SafetyRatings: []*SafetyRating{
							{
								Category:    "HARM_CATEGORY_HARASSMENT",
								Probability: "HIGH",
							},
						},
					},
				}
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(mockResponse)
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGenerateContentRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Empty(t, rsp.Candidates)
				require.NotNil(t, rsp.PromptFeedback)
				assert.Equal(t, "SAFETY", rsp.PromptFeedback.BlockReason)
			},
		},
		{
			name:  "http_error",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "This will fail"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("Internal Server Error"))
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, rsp *GeminiGenerateContentRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Internal Server Error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			rsp, err := p.GenerateContent(context.Background(), tt.model, tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}
func TestGemini_CountTokens(t *testing.T) {
	tests := []struct {
		name          string
		giveReq       *GeminiCountTokensReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, rsp *GeminiCountTokensRsp, err error)
	}{
		{
			name: "success",
			giveReq: &GeminiCountTokensReq{
				Model: "gemini-pro",
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "Why is the sky blue?"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, ":countTokens")
				mockResponse := GeminiCountTokensRsp{
					TotalTokens: 123,
				}
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(mockResponse)
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, rsp *GeminiCountTokensRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, 123, rsp.TotalTokens)
			},
		},
		{
			name: "api_error",
			giveReq: &GeminiCountTokensReq{
				Model: "gemini-pro",
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "Some invalid input"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": {"message": "Invalid input provided"}}`))
			},
			checkResponse: func(t *testing.T, rsp *GeminiCountTokensRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Invalid input provided")
			},
		},
		{
			name: "http_error",
			giveReq: &GeminiCountTokensReq{
				Model: "gemini-pro",
				Contents: []*GeminiContent{
					{
						Parts: GeminiCommonPart{Text: "This will cause a server error"},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("Internal Server Error"))
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, rsp *GeminiCountTokensRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "Internal Server Error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			rsp, err := p.CountTokens(context.Background(), tt.giveReq)

			tt.checkResponse(t, rsp, err)
		})
	}
}
func TestGemini_StreamGenerateContent(t *testing.T) {
	tests := []struct {
		name          string
		model         string
		giveReq       *GeminiGenerateContentReq
		handler       http.HandlerFunc
		checkResponse func(t *testing.T, ch <-chan *GeminiGenerateContentRsp, err error)
	}{
		{
			name:  "success",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{Parts: GeminiCommonPart{Text: "Hello"}},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, ":streamGenerateContent")
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)

				// Simulate streaming response
				responses := []*GeminiGenerateContentRsp{
					{Candidates: []*Candidate{{Content: &GeminiContent{Parts: GeminiCommonPart{Text: "Hi "}}}}},
					{Candidates: []*Candidate{{Content: &GeminiContent{Parts: GeminiCommonPart{Text: "there!"}}}}},
				}

				for _, rsp := range responses {
					jsonData, err := json.Marshal(rsp)
					require.NoError(t, err)
					_, err = w.Write([]byte("data: " + string(jsonData) + "\n\n"))
					require.NoError(t, err)
					// Flush the data to the client
					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
				}
			},
			checkResponse: func(t *testing.T, ch <-chan *GeminiGenerateContentRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, ch)

				var responses []*GeminiGenerateContentRsp
				for rsp := range ch {
					responses = append(responses, rsp)
				}

				require.Len(t, responses, 2)
				assert.Equal(t, "Hi ", responses[0].Candidates[0].Content.Parts.Text)
				assert.Equal(t, "there!", responses[1].Candidates[0].Content.Parts.Text)
			},
		},
		{
			name:  "error_in_stream",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{Parts: GeminiCommonPart{Text: "This will error"}},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)

				// Send one successful chunk
				successRsp := &GeminiGenerateContentRsp{Candidates: []*Candidate{{Content: &GeminiContent{Parts: GeminiCommonPart{Text: "First part"}}}}}
				jsonData, err := json.Marshal(successRsp)
				require.NoError(t, err)
				_, err = w.Write([]byte("data: " + string(jsonData) + "\n\n"))
				require.NoError(t, err)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}

				// Send an error chunk. The current implementation in gemini.go will unmarshal this
				// into a zero-valued GeminiGenerateContentRsp and continue, which is what we'll test for.
				errorData := `{"error": {"code": 400, "message": "An error occurred", "status": "INVALID_ARGUMENT"}}`
				_, err = w.Write([]byte("data: " + errorData + "\n\n"))
				require.NoError(t, err)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			},
			checkResponse: func(t *testing.T, ch <-chan *GeminiGenerateContentRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, ch)

				var receivedResponses []*GeminiGenerateContentRsp
				for rsp := range ch {
					receivedResponses = append(receivedResponses, rsp)
				}

				// We expect two responses. The first is valid, the second is a zero-valued struct
				// because the error in the stream is not currently handled by populating an Error field.
				require.Len(t, receivedResponses, 2)
				assert.Equal(t, "First part", receivedResponses[0].Candidates[0].Content.Parts.Text)
				assert.Empty(t, receivedResponses[1].Candidates) // This should be a zero-valued response
			},
		},
		{
			name:  "http_error",
			model: "gemini-pro",
			giveReq: &GeminiGenerateContentReq{
				Contents: []*GeminiContent{
					{Parts: GeminiCommonPart{Text: "HTTP fail"}},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("Internal Server Error"))
				require.NoError(t, err)
			},
			checkResponse: func(t *testing.T, ch <-chan *GeminiGenerateContentRsp, err error) {
				assert.NoError(t, err) // The function itself doesn't return an error
				require.NotNil(t, ch)

				// The goroutine should exit on HTTP error, closing the channel.
				// We should receive no items.
				var responses []*GeminiGenerateContentRsp
				for rsp := range ch {
					responses = append(responses, rsp)
				}
				assert.Empty(t, responses)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			p := NewGemini(&aiload.ModelChannel{
				Token:   "test-token",
				BaseUrl: server.URL,
			})

			ch, err := p.StreamGenerateContent(context.Background(), tt.model, tt.giveReq)

			tt.checkResponse(t, ch, err)
		})
	}
}