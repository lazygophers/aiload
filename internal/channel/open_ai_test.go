package channel

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOpenAITest(t *testing.T) *OpenAI {
	httpmock.ActivateNonDefault(client.GetClient())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	channel := &aiload.ModelChannel{
		BaseUrl: "https://api.openai.com",
	}
	return NewOpenAI(channel)
}

func TestOpenAICreateChatCompletion(t *testing.T) {
	req := &ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	t.Run("success", func(t *testing.T) {
		o := setupOpenAITest(t)
		mockResp := `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-3.5-turbo-0613",
			"choices": [{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello there!"
				},
				"finish_reason": "stop"
			}],
			"usage": {
				"prompt_tokens": 9,
				"completion_tokens": 12,
				"total_tokens": 21
			}
		}`
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, mockResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		resp, err := o.CreateChatCompletion(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "chatcmpl-123", resp.ID)
		assert.Len(t, resp.Choices, 1)
		assert.Equal(t, "assistant", resp.Choices[0].Message.Role)
		assert.Equal(t, "Hello there!", resp.Choices[0].Message.Content)
	})

	t.Run("http error", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
			httpmock.NewErrorResponder(errors.New("network error")),
		)

		resp, err := o.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		o := setupOpenAITest(t)
		mockErr := `{
			"error": {
				"message": "The model ` + "`gpt-3.5`" + ` does not exist",
				"type": "invalid_request_error",
				"param": null,
				"code": "model_not_found"
			}
		}`
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
			httpmock.NewStringResponder(404, mockErr),
		)

		resp, err := o.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)

		var apiErr *APIError
		require.True(t, errors.As(err, &apiErr))
		assert.Equal(t, 404, apiErr.StatusCode)
		assert.Equal(t, "invalid_request_error", apiErr.Type)
		assert.Equal(t, "model_not_found", apiErr.Code)
	})

	t.Run("bad response", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
			httpmock.NewStringResponder(200, `{"id": "chatcmpl-123",`),
		)

		resp, err := o.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOpenAICreateEmbedding(t *testing.T) {
	req := &EmbeddingRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"Hello", "world"},
	}

	t.Run("success", func(t *testing.T) {
		o := setupOpenAITest(t)
		mockResp := `{
			"object": "list",
			"data": [
				{
					"object": "embedding",
					"embedding": [0.1, 0.2, 0.3],
					"index": 0
				}
			],
			"model": "text-embedding-ada-002-v2",
			"usage": {
				"prompt_tokens": 8,
				"total_tokens": 8
			}
		}`
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, mockResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		resp, err := o.CreateEmbedding(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "list", resp.Object)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, "embedding", resp.Data[0].Object)
		assert.Equal(t, []float64{0.1, 0.2, 0.3}, resp.Data[0].Embedding)
	})

	t.Run("http error", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
			httpmock.NewErrorResponder(errors.New("network error")),
		)

		resp, err := o.CreateEmbedding(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
			httpmock.NewStringResponder(500, `{"error": {"message": "Internal Server Error"}}`),
		)

		resp, err := o.CreateEmbedding(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/embeddings",
			httpmock.NewStringResponder(200, `{"object": "list",`),
		)

		resp, err := o.CreateEmbedding(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOpenAICreateImage(t *testing.T) {
	req := &ImageGenerationRequest{
		Prompt: "a white siamese cat",
		N:      1,
		Size:   "1024x1024",
	}

	t.Run("success", func(t *testing.T) {
		o := setupOpenAITest(t)
		mockResp := `{
			"created": 1589478378,
			"data": [
				{
					"url": "https://oaidalleapiprodscus.blob.core.windows.net/private/org-..."
				}
			]
		}`
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/images/generations",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, mockResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		resp, err := o.CreateImage(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1589478378), resp.Created)
		assert.Len(t, resp.Data, 1)
		assert.NotEmpty(t, resp.Data[0].URL)
	})

	t.Run("http error", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/images/generations",
			httpmock.NewErrorResponder(errors.New("network error")),
		)

		resp, err := o.CreateImage(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/images/generations",
			httpmock.NewStringResponder(500, `{"error": {"message": "Internal Server Error"}}`),
		)

		resp, err := o.CreateImage(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/images/generations",
			httpmock.NewStringResponder(200, `{"created": 1589478378,`),
		)

		resp, err := o.CreateImage(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOpenAICreateAudioTranscription(t *testing.T) {
	// Create a temporary directory and a dummy file for testing
	testDir := "testdata"
	testFile := testDir + "/test.mp3"
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(testFile, []byte("dummy content"), 0644)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})

	t.Run("success", func(t *testing.T) {
		file, err := os.Open(testFile)
		require.NoError(t, err)
		defer file.Close()

		req := &AudioTranscriptionRequest{
			File:     file,
			FileName: "test.mp3",
			Model:    "whisper-1",
		}
		o := setupOpenAITest(t)
		mockResp := `{"text": "Hello, world."}`
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/audio/transcriptions",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20)
				require.NoError(t, err)

				file, _, err := req.FormFile("file")
				require.NoError(t, err)
				defer file.Close()

				assert.Equal(t, "whisper-1", req.FormValue("model"))

				resp := httpmock.NewStringResponse(200, mockResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		resp, err := o.CreateAudioTranscription(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "Hello, world.", resp.Text)
	})

	t.Run("http error", func(t *testing.T) {
		file, err := os.Open(testFile)
		require.NoError(t, err)
		defer file.Close()

		req := &AudioTranscriptionRequest{
			File:     file,
			FileName: "test.mp3",
			Model:    "whisper-1",
		}

		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/audio/transcriptions",
			httpmock.NewErrorResponder(errors.New("network error")),
		)

		resp, err := o.CreateAudioTranscription(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		file, err := os.Open(testFile)
		require.NoError(t, err)
		defer file.Close()

		req := &AudioTranscriptionRequest{
			File:     file,
			FileName: "test.mp3",
			Model:    "whisper-1",
		}

		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/audio/transcriptions",
			httpmock.NewStringResponder(500, `{"error": {"message": "Internal Server Error"}}`),
		)

		resp, err := o.CreateAudioTranscription(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		file, err := os.Open(testFile)
		require.NoError(t, err)
		defer file.Close()

		req := &AudioTranscriptionRequest{
			File:     file,
			FileName: "test.mp3",
			Model:    "whisper-1",
		}

		o := setupOpenAITest(t)
		httpmock.RegisterResponder("POST", "https://api.openai.com/v1/audio/transcriptions",
			httpmock.NewStringResponder(200, `{"text": "Hello, world."`),
		)

		resp, err := o.CreateAudioTranscription(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("file open error", func(t *testing.T) {
		req := &AudioTranscriptionRequest{
			File:     nil, // No file
			FileName: "test.mp3",
			Model:    "whisper-1",
		}
		o := setupOpenAITest(t)
		// This test does not require a mock responder as it should fail before the HTTP request.
		_, err := o.CreateAudioTranscription(req)
		require.Error(t, err)
	})
}
