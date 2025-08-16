package channel

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = time.Time{}

func setupOllamaTest(t *testing.T) *Ollama {
	mockClient := resty.New()
	httpmock.ActivateNonDefault(mockClient.GetClient())

	originalClient := client
	client = mockClient
	t.Cleanup(func() {
		client = originalClient
		httpmock.DeactivateAndReset()
	})

	channel := &aiload.ModelChannel{}
	return NewOllama(channel)
}

func TestOllama_GetLocalModelList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{
			"models":[
				{
					"name":"llama2:latest",
					"model":"llama2:latest",
					"modified_at":"2023-11-20T14:56:22.2891361Z",
					"size":3826793677,
					"digest":"a2f7e2f1d8a25a385a1e80921aa91e027048a12ed084ac825c8a4d46d1522030",
					"details":{
						"parent_model":"",
						"format":"gguf",
						"family":"llama",
						"families":["llama"],
						"parameter_size":"7B",
						"quantization_level":"Q4_0"
					}
				}
			]
		}`
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/tags", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.GetLocalModelList()
		require.NoError(t, err)
		require.NotNil(t, resp)
		if assert.Len(t, resp.Models, 1) {
			assert.Equal(t, "llama2:latest", resp.Models[0].Name)
		}
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/tags", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.GetLocalModelList()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/tags", httpmock.NewStringResponder(404, ""))

		resp, err := ollama.GetLocalModelList()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/tags", httpmock.NewStringResponder(200, `{"models":[`))

		resp, err := ollama.GetLocalModelList()
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_GenerateCompletion(t *testing.T) {
	req := &OllamaGenerateReq{Model: "test-model", Prompt: "why is the sky blue?"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"model":"test-model","created_at":"2023-11-20T15:00:00Z","response":"because of science","done":true}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.GenerateCompletion(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "test-model", resp.Model)
		assert.Equal(t, "because of science", resp.Response)
		assert.True(t, resp.Done)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.GenerateCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate", httpmock.NewStringResponder(500, ""))

		resp, err := ollama.GenerateCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate", httpmock.NewStringResponder(200, `{"model":"`))

		resp, err := ollama.GenerateCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_CreateChatCompletion(t *testing.T) {
	req := &OllamaChatReq{Model: "test-model", Messages: []OllamaMessage{{Role: "user", Content: "hello"}}}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"model":"test-model","created_at":"2023-11-20T15:05:00Z","message":{"role":"assistant","content":"world"},"done":true}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/chat", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.CreateChatCompletion(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "test-model", resp.Model)
		assert.Equal(t, "assistant", resp.Message.Role)
		assert.Equal(t, "world", resp.Message.Content)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/chat", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/chat", httpmock.NewStringResponder(400, ""))

		resp, err := ollama.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/chat", httpmock.NewStringResponder(200, `{"model":"`))

		resp, err := ollama.CreateChatCompletion(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_CreateModel(t *testing.T) {
	req := &OllamaCreateModelReq{Name: "custom-model", Modelfile: "FROM llama2"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"status":"success"}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/create", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.CreateModel(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "success", resp.Status)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/create", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.CreateModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/create", httpmock.NewStringResponder(500, ""))

		resp, err := ollama.CreateModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/create", httpmock.NewStringResponder(200, `{"status":`))

		resp, err := ollama.CreateModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_ShowModelInfo(t *testing.T) {
	req := &OllamaShowModelReq{Name: "test-model"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"modelfile":"FROM test","license":"MIT","parameters":"stop [INST]","template":"[INST] {{ .Prompt }} [/INST]"}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/show", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.ShowModelInfo(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "FROM test", resp.Modelfile)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/show", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.ShowModelInfo(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/show", httpmock.NewStringResponder(404, ""))

		resp, err := ollama.ShowModelInfo(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/show", httpmock.NewStringResponder(200, `{"modelfile":`))

		resp, err := ollama.ShowModelInfo(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_CopyModel(t *testing.T) {
	req := &OllamaCopyModelReq{Source: "model-a", Destination: "model-b"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/copy", httpmock.NewStringResponder(200, ""))

		err := ollama.CopyModel(req)
		require.NoError(t, err)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/copy", httpmock.NewErrorResponder(errors.New("network error")))

		err := ollama.CopyModel(req)
		require.Error(t, err)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/copy", httpmock.NewStringResponder(404, ""))

		err := ollama.CopyModel(req)
		require.Error(t, err)
	})
}

func TestOllama_DeleteModel(t *testing.T) {
	req := &OllamaDeleteModelReq{Name: "test-model"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("DELETE", "http://localhost:11434/api/delete", httpmock.NewStringResponder(200, ""))

		err := ollama.DeleteModel(req)
		require.NoError(t, err)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("DELETE", "http://localhost:11434/api/delete", httpmock.NewErrorResponder(errors.New("network error")))

		err := ollama.DeleteModel(req)
		require.Error(t, err)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("DELETE", "http://localhost:11434/api/delete", httpmock.NewStringResponder(404, ""))

		err := ollama.DeleteModel(req)
		require.Error(t, err)
	})
}

func TestOllama_PullModel(t *testing.T) {
	req := &OllamaPullModelReq{Name: "test-model"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"status":"success"}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/pull", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.PullModel(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "success", resp.Status)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/pull", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.PullModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/pull", httpmock.NewStringResponder(404, ""))

		resp, err := ollama.PullModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/pull", httpmock.NewStringResponder(200, `{"status":`))

		resp, err := ollama.PullModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_PushModel(t *testing.T) {
	req := &OllamaPushModelReq{Name: "test-model"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"status":"success"}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/push", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.PushModel(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "success", resp.Status)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/push", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.PushModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/push", httpmock.NewStringResponder(404, ""))

		resp, err := ollama.PushModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/push", httpmock.NewStringResponder(200, `{"status":`))

		resp, err := ollama.PushModel(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_GenerateEmbeddings(t *testing.T) {
	req := &OllamaEmbeddingsReq{Model: "test-model", Prompt: "hello"}
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"embedding":[0.1, 0.2, 0.3]}`
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/embeddings", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.GenerateEmbeddings(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, []float64{0.1, 0.2, 0.3}, resp.Embedding)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/embeddings", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.GenerateEmbeddings(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/embeddings", httpmock.NewStringResponder(500, ""))

		resp, err := ollama.GenerateEmbeddings(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("POST", "http://localhost:11434/api/embeddings", httpmock.NewStringResponder(200, `{"embedding":`))

		resp, err := ollama.GenerateEmbeddings(req)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_ListRunningModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{
			"models": [
				{
					"name": "test-model:latest",
					"model": "test-model:latest",
					"size": 12345,
					"digest": "abcdef123456",
					"details": {},
					"expires_at": "2024-01-01T00:00:00Z",
					"size_vram": 12345678
				}
			]
		}`
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/ps", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.ListRunningModels()
		require.NoError(t, err)
		require.NotNil(t, resp)
		if assert.Len(t, resp.Models, 1) {
			assert.Equal(t, "test-model:latest", resp.Models[0].Name)
		}
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/ps", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.ListRunningModels()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/ps", httpmock.NewStringResponder(500, ""))

		resp, err := ollama.ListRunningModels()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/ps", httpmock.NewStringResponder(200, `{"models":[`))

		resp, err := ollama.ListRunningModels()
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_GetVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		mockResp := `{"version":"0.1.15"}`
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/version", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, mockResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		resp, err := ollama.GetVersion()
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "0.1.15", resp.Version)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/version", httpmock.NewErrorResponder(errors.New("network error")))

		resp, err := ollama.GetVersion()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("api error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/version", httpmock.NewStringResponder(500, ""))

		resp, err := ollama.GetVersion()
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("bad response", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("GET", "http://localhost:11434/api/version", httpmock.NewStringResponder(200, `{"version":`))

		resp, err := ollama.GetVersion()
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOllama_CheckBlobExists(t *testing.T) {
	digest := "sha256:12345"
	t.Run("exists", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("HEAD", "http://localhost:11434/api/blobs/"+digest, httpmock.NewStringResponder(http.StatusOK, ""))

		exists, err := ollama.CheckBlobExists(digest)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("HEAD", "http://localhost:11434/api/blobs/"+digest, httpmock.NewStringResponder(http.StatusNotFound, ""))

		exists, err := ollama.CheckBlobExists(digest)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("http error", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("HEAD", "http://localhost:11434/api/blobs/"+digest, httpmock.NewErrorResponder(errors.New("network error")))

		exists, err := ollama.CheckBlobExists(digest)
		require.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("other status code", func(t *testing.T) {
		ollama := setupOllamaTest(t)
		httpmock.RegisterResponder("HEAD", "http://localhost:11434/api/blobs/"+digest, httpmock.NewStringResponder(500, ""))

		exists, err := ollama.CheckBlobExists(digest)
		require.Error(t, err)
		assert.False(t, exists)
	})
}
