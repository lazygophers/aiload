package channel

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// 禁用测试期间的日志输出，保持测试结果的清洁
	log.SetOutput(io.Discard)
	// 运行包中的所有测试
	exitCode := m.Run()
	// 恢复日志输出到标准错误
	log.SetOutput(os.Stderr)
	// 退出并返回测试结果的退出码
	os.Exit(exitCode)
}

func TestSiliconFlow_CreateChatCompletion(t *testing.T) {
	// Correctly initialize the provider using its constructor
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	// Since the code uses a global client, we must mock the global client.
	// This is not ideal, but necessary to test the current implementation.
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// Mock the API response
	mockResponse := SiliconFlowCreateChatCompletionRsp{
		ID:      "chatcmpl-mock-id",
		Object:  "chat.completion",
		Created: 1677652288,
		Model:   "deepseek-coder",
		Choices: []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{
					Role:    "assistant",
					Content: "\n\nHello there, how may I assist you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     9,
			CompletionTokens: 12,
			TotalTokens:      21,
		},
	}
	respBody, err := json.Marshal(mockResponse)
	assert.NoError(t, err)

	httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/chat/completions",
		func(req *http.Request) (*http.Response, error) {
			// Check for the Authorization header
			if req.Header.Get("Authorization") != "Bearer test-token" {
				return httpmock.NewStringResponse(401, "Unauthorized"), nil
			}
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Create a request
	req := &SiliconFlowCreateChatCompletionReq{
		Model: "deepseek-coder",
		Messages: []SiliconFlowCreateChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	// Call the function
	rsp, err := p.CreateChatCompletion(req)

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, "chatcmpl-mock-id", rsp.ID)
	assert.NotEmpty(t, rsp.Choices)
	assert.Equal(t, "assistant", rsp.Choices[0].Message.Role)
	assert.Equal(t, "\n\nHello there, how may I assist you today?", rsp.Choices[0].Message.Content)
	assert.Equal(t, "stop", rsp.Choices[0].FinishReason)
}

func TestSiliconFlow_CreateChatCompletion_Error(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// Mock the API error response
	httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/chat/completions",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Create a request
	req := &SiliconFlowCreateChatCompletionReq{
		Model: "deepseek-coder",
		Messages: []SiliconFlowCreateChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	// Call the function
	rsp, err := p.CreateChatCompletion(req)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, rsp)
	var siliconFlowError *SiliconFlowErrorRsp
	require.ErrorAs(t, err, &siliconFlowError)
	assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
}

func TestSiliconFlow_CreateEmbedding(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateEmbeddingRsp{
			Object: "list",
			Data: []struct {
				Object    string    `json:"object"`
				Embedding []float64 `json:"embedding"`
				Index     int       `json:"index"`
			}{
				{
					Object:    "embedding",
					Embedding: []float64{0.1, 0.2, 0.3},
					Index:     0,
				},
			},
			Model: "bge-large-zh-v1.5",
			Usage: struct {
				PromptTokens int `json:"prompt_tokens"`
				TotalTokens  int `json:"total_tokens"`
			}{
				PromptTokens: 10,
				TotalTokens:  10,
			},
		}
		respBody, err := json.Marshal(mockResponse)
		assert.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/embeddings",
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateEmbeddingReq{
			Model: "bge-large-zh-v1.5",
			Input: []string{"hello"},
		}

		rsp, err := p.CreateEmbedding(req)

		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, []float64{0.1, 0.2, 0.3}, rsp.Data[0].Embedding)
	})

	t.Run("error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/embeddings",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateEmbeddingReq{
			Model: "bge-large-zh-v1.5",
			Input: []string{"hello"},
		}

		rsp, err := p.CreateEmbedding(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})
}

func TestSiliconFlow_CreateImageGenerations(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateImageGenerationsRsp{
			Created: 1677652288,
			Data: []struct {
				URL     string `json:"url,omitempty"`
				B64JSON string `json:"b64_json,omitempty"`
			}{
				{
					URL: "https://example.com/image.png",
				},
			},
		}
		respBody, err := json.Marshal(mockResponse)
		assert.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/generations",
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}
				var reqBody SiliconFlowCreateImageGenerationsReq
				err := json.NewDecoder(req.Body).Decode(&reqBody)
				assert.NoError(t, err)
				assert.Equal(t, "a cat", reqBody.Prompt)

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageGenerationsReq{
			Prompt: "a cat",
		}

		rsp, err := p.CreateImageGenerations(req)

		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, "https://example.com/image.png", rsp.Data[0].URL)
	})

	t.Run("error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/generations",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageGenerationsReq{
			Prompt: "a cat",
		}

		rsp, err := p.CreateImageGenerations(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})
}

func TestSiliconFlow_CreateImageEdits(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateImageEditsRsp{
			Created: 1677652288,
			Data: []struct {
				URL     string `json:"url,omitempty"`
				B64JSON string `json:"b64_json,omitempty"`
			}{
				{
					URL: "https://example.com/edited-image.png",
				},
			},
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/edits",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20) // 10 MB
				require.NoError(t, err)

				// Check form fields
				assert.Equal(t, "a cute cat", req.FormValue("prompt"))
				assert.Equal(t, "stable-diffusion-xl-1024-v1-0", req.FormValue("model"))
				assert.Equal(t, "1", req.FormValue("n"))
				assert.Equal(t, "1024x1024", req.FormValue("size"))
				assert.Equal(t, "url", req.FormValue("response_format"))

				// Check image file
				file, header, err := req.FormFile("image")
				require.NoError(t, err)
				defer file.Close()
				assert.Equal(t, "test-image.png", header.Filename)
				imgData, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, "fake-image-data", string(imgData))

				// Check mask file
				maskFile, maskHeader, err := req.FormFile("mask")
				require.NoError(t, err)
				defer maskFile.Close()
				assert.Equal(t, "test-mask.png", maskHeader.Filename)
				maskData, err := io.ReadAll(maskFile)
				require.NoError(t, err)
				assert.Equal(t, "fake-mask-data", string(maskData))

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageEditsReq{
			Prompt:         "a cute cat",
			Image:          strings.NewReader("fake-image-data"),
			ImageFileName:  "test-image.png",
			Mask:           strings.NewReader("fake-mask-data"),
			MaskFileName:   "test-mask.png",
			Model:          "stable-diffusion-xl-1024-v1-0",
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: "url",
		}

		rsp, err := p.CreateImageEdits(req)

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		require.Len(t, rsp.Data, 1)
		assert.Equal(t, "https://example.com/edited-image.png", rsp.Data[0].URL)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/edits",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageEditsReq{
			Prompt:        "a cute cat",
			Image:         strings.NewReader("fake-image-data"),
			ImageFileName: "test-image.png",
			Model:         "stable-diffusion-xl-1024-v1-0",
		}

		rsp, err := p.CreateImageEdits(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})

	t.Run("missing image", func(t *testing.T) {
		req := &SiliconFlowCreateImageEditsReq{
			Prompt: "a cute cat",
			// Image is nil
		}

		rsp, err := p.CreateImageEdits(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		assert.Contains(t, err.Error(), "image reader is required")
	})
}

func TestSiliconFlow_CreateImageVariations(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateImageVariationsRsp{
			Created: 1677652288,
			Data: []struct {
				URL     string `json:"url,omitempty"`
				B64JSON string `json:"b64_json,omitempty"`
			}{
				{
					URL: "https://example.com/variation-image.png",
				},
			},
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/variations",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20) // 10 MB
				require.NoError(t, err)

				// Check form fields
				assert.Equal(t, "stable-diffusion-xl-1024-v1-0", req.FormValue("model"))
				assert.Equal(t, "1", req.FormValue("n"))
				assert.Equal(t, "1024x1024", req.FormValue("size"))
				assert.Equal(t, "url", req.FormValue("response_format"))

				// Check image file
				file, header, err := req.FormFile("image")
				require.NoError(t, err)
				defer file.Close()
				assert.Equal(t, "test-variation-image.png", header.Filename)
				imgData, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, "fake-variation-image-data", string(imgData))

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageVariationsReq{
			Image:          strings.NewReader("fake-variation-image-data"),
			ImageFileName:  "test-variation-image.png",
			Model:          "stable-diffusion-xl-1024-v1-0",
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: "url",
		}

		rsp, err := p.CreateImageVariations(req)

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		require.Len(t, rsp.Data, 1)
		assert.Equal(t, "https://example.com/variation-image.png", rsp.Data[0].URL)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/images/variations",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateImageVariationsReq{
			Image:         strings.NewReader("fake-image-data"),
			ImageFileName: "test-image.png",
			Model:         "stable-diffusion-xl-1024-v1-0",
		}

		rsp, err := p.CreateImageVariations(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})

	t.Run("missing image", func(t *testing.T) {
		req := &SiliconFlowCreateImageVariationsReq{
			// Image is nil
		}

		rsp, err := p.CreateImageVariations(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		assert.Equal(t, "image reader is required", err.Error())
	})
}

func TestSiliconFlow_CreateAudioTranslation(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateAudioTranslationRsp{
			Text: "Hello, this is a test.",
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/audio/translations",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20) // 10 MB
				require.NoError(t, err)

				// Check form fields
				assert.Equal(t, "whisper-large-v3", req.FormValue("model"))
				assert.Equal(t, "Translate to English", req.FormValue("prompt"))
				assert.Equal(t, "json", req.FormValue("response_format"))
				assert.Equal(t, "0.7", req.FormValue("temperature"))

				// Check file
				file, header, err := req.FormFile("file")
				require.NoError(t, err)
				defer file.Close()
				assert.Equal(t, "test-audio.mp3", header.Filename)
				audioData, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, "fake-audio-data", string(audioData))

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateAudioTranslationReq{
			File:           strings.NewReader("fake-audio-data"),
			FileName:       "test-audio.mp3",
			Model:          "whisper-large-v3",
			Prompt:         "Translate to English",
			ResponseFormat: "json",
			Temperature:    0.7,
		}

		rsp, err := p.CreateAudioTranslation(req)

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "Hello, this is a test.", rsp.Text)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/audio/translations",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateAudioTranslationReq{
			File:     strings.NewReader("fake-audio-data"),
			FileName: "test-audio.mp3",
			Model:    "whisper-large-v3",
		}

		rsp, err := p.CreateAudioTranslation(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})

	t.Run("missing file", func(t *testing.T) {
		req := &SiliconFlowCreateAudioTranslationReq{
			Model: "whisper-large-v3",
			// File is nil
		}

		rsp, err := p.CreateAudioTranslation(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		assert.Equal(t, "file reader is required", err.Error())
	})
}

func TestSiliconFlow_CreateAudioTranscription(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowCreateAudioTranscriptionResp{
			Text: "This is a test transcription.",
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/audio/transcriptions",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20) // 10 MB
				require.NoError(t, err)

				// Check form fields
				assert.Equal(t, "whisper-large-v3", req.FormValue("model"))
				assert.Equal(t, "Transcribe this audio", req.FormValue("prompt"))
				assert.Equal(t, "json", req.FormValue("response_format"))
				assert.Equal(t, "0.8", req.FormValue("temperature"))

				// Check file
				file, header, err := req.FormFile("file")
				require.NoError(t, err)
				defer file.Close()
				assert.Equal(t, "test-audio.mp3", header.Filename)
				audioData, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, "fake-audio-data-for-transcription", string(audioData))

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateAudioTranscriptionReq{
			File:           strings.NewReader("fake-audio-data-for-transcription"),
			FileName:       "test-audio.mp3",
			Model:          "whisper-large-v3",
			Prompt:         "Transcribe this audio",
			ResponseFormat: "json",
			Temperature:    0.8,
		}

		rsp, err := p.CreateAudioTranscription(req)

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "This is a test transcription.", rsp.Text)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/audio/transcriptions",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowCreateAudioTranscriptionReq{
			File:     strings.NewReader("fake-audio-data-for-transcription"),
			FileName: "test-audio.mp3",
			Model:    "whisper-large-v3",
		}

		rsp, err := p.CreateAudioTranscription(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})

	t.Run("missing file", func(t *testing.T) {
		req := &SiliconFlowCreateAudioTranscriptionReq{
			Model: "whisper-large-v3",
			// File is nil
		}

		rsp, err := p.CreateAudioTranscription(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		assert.Equal(t, "file reader is required", err.Error())
	})
}

func TestSiliconFlow_UploadFile(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := SiliconFlowFileObj{
			ID:        "file-mock-id",
			Bytes:     140,
			CreatedAt: 1677652288,
			Filename:  "test-file.jsonl",
			Object:    "file",
			Purpose:   "fine-tune",
			Status:    "processed",
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/files",
			func(req *http.Request) (*http.Response, error) {
				err := req.ParseMultipartForm(10 << 20) // 10 MB
				require.NoError(t, err)

				// Check form fields
				assert.Equal(t, "fine-tune", req.FormValue("purpose"))

				// Check file
				file, header, err := req.FormFile("file")
				require.NoError(t, err)
				defer file.Close()
				assert.Equal(t, "test-file.jsonl", header.Filename)
				fileData, err := io.ReadAll(file)
				require.NoError(t, err)
				assert.Equal(t, "fake-file-data", string(fileData))

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowUploadFileReq{
			File:     strings.NewReader("fake-file-data"),
			FileName: "test-file.jsonl",
			Purpose:  "fine-tune",
		}

		rsp, err := p.UploadFile(req)

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "file-mock-id", rsp.ID)
		assert.Equal(t, "fine-tune", rsp.Purpose)
		assert.Equal(t, "test-file.jsonl", rsp.Filename)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/files",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		req := &SiliconFlowUploadFileReq{
			File:     strings.NewReader("fake-file-data"),
			FileName: "test-file.jsonl",
			Purpose:  "fine-tune",
		}

		rsp, err := p.UploadFile(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})

	t.Run("missing file", func(t *testing.T) {
		req := &SiliconFlowUploadFileReq{
			Purpose: "fine-tune",
			// File is nil
		}

		rsp, err := p.UploadFile(req)

		assert.Error(t, err)
		assert.Nil(t, rsp)
		assert.Equal(t, "file reader is required", err.Error())
	})
}
func TestSiliconFlow_ListFiles(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := &SiliconFlowListFilesRsp{
			Object: "list",
			Data: []SiliconFlowFileObj{
				{
					ID:        "file-mock-id-1",
					Bytes:     140,
					CreatedAt: 1677652288,
					Filename:  "test-file-1.jsonl",
					Object:    "file",
					Purpose:   "fine-tune",
					Status:    "processed",
				},
				{
					ID:        "file-mock-id-2",
					Bytes:     256,
					CreatedAt: 1677652289,
					Filename:  "test-file-2.jsonl",
					Object:    "file",
					Purpose:   "fine-tune",
					Status:    "processed",
				},
			},
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/files",
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.ListFiles()

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		require.Len(t, rsp.Data, 2)
		assert.Equal(t, "file-mock-id-1", rsp.Data[0].ID)
		assert.Equal(t, "test-file-2.jsonl", rsp.Data[1].Filename)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/files",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.ListFiles()

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})
}

func TestSiliconFlow_DeleteFile(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := &SiliconFlowDeleteFileRsp{
			ID:      "file-mock-id-to-delete",
			Object:  "file",
			Deleted: true,
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("DELETE", "https://api.siliconflow.cn/v1/files/file-mock-id-to-delete",
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.DeleteFile("file-mock-id-to-delete")

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "file-mock-id-to-delete", rsp.ID)
		assert.True(t, rsp.Deleted)
		assert.Equal(t, "file", rsp.Object)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("DELETE", "https://api.siliconflow.cn/v1/files/file-mock-id-for-error",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.DeleteFile("file-mock-id-for-error")

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})
}
func TestSiliconflow_RetrieveFile(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("success", func(t *testing.T) {
		mockResponse := &SiliconFlowFileObj{
			ID:        "file-mock-id-1",
			Bytes:     140,
			CreatedAt: 1677652288,
			Filename:  "test-file-1.jsonl",
			Object:    "file",
			Purpose:   "fine-tune",
			Status:    "processed",
		}
		respBody, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/files/file-mock-id-1",
			func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.RetrieveFile("file-mock-id-1")

		assert.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "file-mock-id-1", rsp.ID)
		assert.Equal(t, "test-file-1.jsonl", rsp.Filename)
	})

	t.Run("api error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/files/file-mock-id-for-error",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
		)

		rsp, err := p.RetrieveFile("file-mock-id-for-error")

		assert.Error(t, err)
		assert.Nil(t, rsp)
		var siliconFlowError *SiliconFlowErrorRsp
		require.ErrorAs(t, err, &siliconFlowError)
		assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
	})
}

func TestSiliconFlow_CreateFineTuningJob(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 定义测试用例
	tests := []struct {
		name          string
		giveRequest   *SiliconFlowFineTuningJobRequest
		mockResponder httpmock.Responder
		wantErr       bool
		wantErrAs     interface{}
		checkResponse func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error)
	}{
		{
			name: "success",
			giveRequest: &SiliconFlowFineTuningJobRequest{
				TrainingFile: "file-train-mock-id",
				Model:        "gpt-3.5-turbo",
			},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				mockResponse := SiliconFlowFineTuningJob{
					ID:             "ft-job-mock-id",
					Object:         "fine_tuning.job",
					Model:          "gpt-3.5-turbo",
					CreatedAt:      1677652288,
					Status:         "succeeded",
					FineTunedModel: "ft:gpt-3.5-turbo:my-org:custom-model-name:1",
					TrainingFile:   "file-train-mock-id",
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}

				var reqBody SiliconFlowFineTuningJobRequest
				err = json.NewDecoder(req.Body).Decode(&reqBody)
				require.NoError(t, err)
				assert.Equal(t, "file-train-mock-id", reqBody.TrainingFile)
				assert.Equal(t, "gpt-3.5-turbo", reqBody.Model)

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: false,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "ft-job-mock-id", rsp.ID)
				assert.Equal(t, "succeeded", rsp.Status)
				assert.Equal(t, "ft:gpt-3.5-turbo:my-org:custom-model-name:1", rsp.FineTunedModel)
			},
		},
		{
			name: "api error",
			giveRequest: &SiliconFlowFineTuningJobRequest{
				TrainingFile: "file-train-mock-id",
				Model:        "gpt-3.5-turbo",
			},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr:   true,
			wantErrAs: &SiliconFlowErrorRsp{},
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
			},
		},
		{
			name: "network error",
			giveRequest: &SiliconFlowFineTuningJobRequest{
				TrainingFile: "file-train-mock-id",
				Model:        "gpt-3.5-turbo",
			},
			mockResponder: httpmock.NewErrorResponder(errors.New("network connection error")),
			wantErr:       true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network connection error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset() // 每个子测试前重置 mock
			httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/fine_tuning/jobs", tt.mockResponder)

			rsp, err := p.CreateFineTuningJob(tt.giveRequest)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rsp, err)
			}
		})
	}
}
func TestSiliconFlow_ListFineTuningJobs(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 定义测试用例
	tests := []struct {
		name          string
		giveLimit     int
		giveAfter     string
		mockResponder httpmock.Responder
		wantErr       bool
		wantErrAs     interface{}
		checkResponse func(t *testing.T, rsp *SiliconFlowFineTuningJobList, err error)
	}{
		{
			name:      "success",
			giveLimit: 10,
			giveAfter: "job-id-123",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 验证查询参数
				assert.Equal(t, "10", req.URL.Query().Get("limit"))
				assert.Equal(t, "job-id-123", req.URL.Query().Get("after"))

				mockResponse := SiliconFlowFineTuningJobList{
					Object: "list",
					Data: []SiliconFlowFineTuningJob{
						{
							ID:             "ft-job-mock-id-1",
							Object:         "fine_tuning.job",
							Model:          "gpt-3.5-turbo",
							CreatedAt:      1677652288,
							Status:         "succeeded",
							FineTunedModel: "ft:gpt-3.5-turbo:my-org:custom-model-name:1",
							TrainingFile:   "file-train-mock-id-1",
						},
					},
					HasMore: false,
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: false,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJobList, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "list", rsp.Object)
				require.Len(t, rsp.Data, 1)
				assert.Equal(t, "ft-job-mock-id-1", rsp.Data[0].ID)
				assert.False(t, rsp.HasMore)
			},
		},
		{
			name:      "api error",
			giveLimit: 5,
			mockResponder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr:   true,
			wantErrAs: &SiliconFlowErrorRsp{},
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJobList, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
			},
		},
		{
			name:          "network error",
			mockResponder: httpmock.NewErrorResponder(errors.New("network connection error")),
			wantErr:       true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJobList, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network connection error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset() // 每个子测试前重置 mock
			httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/fine_tuning/jobs", tt.mockResponder)

			rsp, err := p.ListFineTuningJobs(tt.giveLimit, tt.giveAfter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rsp, err)
			}
		})
	}
}

func TestSiliconFlow_RetrieveFineTuningJob(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 定义测试用例
	tests := []struct {
		name          string
		giveJobID     string
		mockResponder httpmock.Responder
		wantErr       bool
		checkResponse func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error)
	}{
		{
			name:      "success",
			giveJobID: "ft-job-mock-id-success",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				mockResponse := SiliconFlowFineTuningJob{
					ID:             "ft-job-mock-id-success",
					Object:         "fine_tuning.job",
					Model:          "gpt-3.5-turbo",
					CreatedAt:      1677652288,
					Status:         "succeeded",
					FineTunedModel: "ft:gpt-3.5-turbo:my-org:custom-model-name:1",
					TrainingFile:   "file-train-mock-id-1",
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: false,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "ft-job-mock-id-success", rsp.ID)
				assert.Equal(t, "succeeded", rsp.Status)
				assert.Equal(t, "ft:gpt-3.5-turbo:my-org:custom-model-name:1", rsp.FineTunedModel)
			},
		},
		{
			name:      "not found error",
			giveJobID: "ft-job-mock-id-not-found",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(404, `{"error":{"message":"Job not found"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Job not found", siliconFlowError.Err.Message)
			},
		},
		{
			name:      "internal server error",
			giveJobID: "ft-job-mock-id-server-error",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"error":{"message":"Internal server error"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Internal server error", siliconFlowError.Err.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset() // 每个子测试前重置 mock
			url := "https://api.siliconflow.cn/v1/fine_tuning/jobs/" + tt.giveJobID
			httpmock.RegisterResponder("GET", url, tt.mockResponder)

			rsp, err := p.RetrieveFineTuningJob(tt.giveJobID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rsp, err)
			}
		})
	}
}

func TestSiliconFlow_CancelFineTuningJob(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 定义测试用例
	tests := []struct {
		name          string
		giveJobID     string
		mockResponder httpmock.Responder
		wantErr       bool
		checkResponse func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error)
	}{
		{
			name:      "success",
			giveJobID: "ft-job-mock-id-success",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 模拟一个成功取消的响应
				mockResponse := SiliconFlowFineTuningJob{
					ID:             "ft-job-mock-id-success",
					Object:         "fine_tuning.job",
					Model:          "gpt-3.5-turbo",
					CreatedAt:      1677652288,
					Status:         "cancelled", // 关键状态：已取消
					FineTunedModel: "",
					TrainingFile:   "file-train-mock-id-1",
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				if req.Header.Get("Authorization") != "Bearer test-token" {
					return httpmock.NewStringResponse(401, "Unauthorized"), nil
				}

				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: false,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "ft-job-mock-id-success", rsp.ID)
				assert.Equal(t, "cancelled", rsp.Status)
			},
		},
		{
			name:      "not found error",
			giveJobID: "ft-job-mock-id-not-found",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 模拟任务不存在的 404 错误
				resp := httpmock.NewStringResponse(404, `{"error":{"message":"Job not found"}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Job not found", siliconFlowError.Err.Message)
			},
		},
		{
			name:      "cannot cancel error due to status",
			giveJobID: "ft-job-mock-id-already-completed",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 模拟任务因已完成而无法取消的错误
				resp := httpmock.NewStringResponse(400, `{"error":{"message":"Job has already completed, cannot cancel."}}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantErr: true,
			checkResponse: func(t *testing.T, rsp *SiliconFlowFineTuningJob, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Job has already completed, cannot cancel.", siliconFlowError.Err.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset() // 每个子测试前重置 mock
			url := "https://api.siliconflow.cn/v1/fine_tuning/jobs/" + tt.giveJobID + "/cancel"
			httpmock.RegisterResponder("POST", url, tt.mockResponder)

			rsp, err := p.CancelFineTuningJob(tt.giveJobID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rsp, err)
			}
		})
	}
}

func TestGetModelSuccess(t *testing.T) {
	// Initialize the provider
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	// Activate the HTTP mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// Define the model ID and the mock response
	modelID := "deepseek-coder"
	mockResponse := Model{
		Id:      modelID,
		Object:  "model",
		Created: 1696931010,
		OwnedBy: "deepseek",
	}

	// Create a JSON responder
	responder, err := httpmock.NewJsonResponder(200, mockResponse)
	require.NoError(t, err)

	// Register the responder for the specific URL
	url := "https://api.siliconflow.cn/v1/models/" + modelID
	httpmock.RegisterResponder("GET", url, responder)

	// Call the method under test
	model, err := p.GetModel(modelID)

	// Assert the results
	assert.NoError(t, err)
	require.NotNil(t, model)
	assert.Equal(t, mockResponse.Id, model.Id)
	assert.Equal(t, mockResponse.Object, model.Object)
	assert.Equal(t, mockResponse.Created, model.Created)
	assert.Equal(t, mockResponse.OwnedBy, model.OwnedBy)
}
func TestGetModelFailure(t *testing.T) {
	// Initialize the provider
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	// Activate the HTTP mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// Define the model ID that does not exist
	modelID := "non-existent-model"
	errorMessage := "model not found"

	// Mock the API error response
	mockErrorResponse := SiliconFlowErrorRsp{
		Err: struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Param   string `json:"param"`
			Code    string `json:"code"`
		}{
			Message: errorMessage,
			Type:    "invalid_request_error",
		},
	}

	// Create a JSON responder for the error
	responder, err := httpmock.NewJsonResponder(http.StatusNotFound, mockErrorResponse)
	require.NoError(t, err)

	// Register the responder for the specific URL
	url := "https://api.siliconflow.cn/v1/models/" + modelID
	httpmock.RegisterResponder("GET", url, responder)

	// Call the method under test
	model, err := p.GetModel(modelID)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, model)

	// Check if the error is of the expected type
	var siliconFlowError *SiliconFlowErrorRsp
	require.ErrorAs(t, err, &siliconFlowError)
	assert.Equal(t, errorMessage, siliconFlowError.Err.Message)
}
func TestSiliconFlow_GetInfoList(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name          string
		setupMock     func()
		giveChannel   *aiload.ModelChannel
		checkResponse func(t *testing.T, rsp *SiliconFlowGetInfoListRsp, err error)
	}{
		{
			name: "success",
			setupMock: func() {
				mockResponse := SiliconFlowGetInfoListRsp{
					Code:    200,
					Message: "Success",
					Status:  true,
					Data: struct {
						Id            string `json:"id"`
						Name          string `json:"name"`
						Image         string `json:"image"`
						Email         string `json:"email"`
						IsAdmin       bool   `json:"isAdmin"`
						Balance       string `json:"balance"`
						Status        string `json:"status"`
						Introduction  string `json:"introduction"`
						Role          string `json:"role"`
						ChargeBalance string `json:"chargeBalance"`
						TotalBalance  string `json:"totalBalance"`
					}{
						Id:      "user-123",
						Name:    "Test User",
						Email:   "test@example.com",
						Balance: "100.00",
					},
				}
				respBody, err := json.Marshal(mockResponse)
				require.NoError(t, err)

				httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/user/info",
					func(req *http.Request) (*http.Response, error) {
						if req.Header.Get("Authorization") != "Bearer test-token" {
							return httpmock.NewStringResponse(401, "Unauthorized"), nil
						}
						resp := httpmock.NewStringResponse(200, string(respBody))
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				)
			},
			giveChannel: mockChannel,
			checkResponse: func(t *testing.T, rsp *SiliconFlowGetInfoListRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.True(t, rsp.Status)
				assert.Equal(t, "user-123", rsp.Data.Id)
				assert.Equal(t, "Test User", rsp.Data.Name)
			},
		},
		{
			name: "api error",
			setupMock: func() {
				httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/user/info",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(400, `{"error":{"message":"Bad request"}}`)
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				)
			},
			giveChannel: mockChannel,
			checkResponse: func(t *testing.T, rsp *SiliconFlowGetInfoListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "Bad request", siliconFlowError.Err.Message)
			},
		},
		{
			name: "no token",
			setupMock: func() {
				// No mock needed as it should fail before the request
			},
			giveChannel: &aiload.ModelChannel{Token: ""}, // No token
			checkResponse: func(t *testing.T, rsp *SiliconFlowGetInfoListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.True(t, errors.Is(err, ErrAuthTokenNil))
			},
		},
		{
			name: "http error",
			setupMock: func() {
				httpmock.RegisterResponder("GET", "https://api.siliconflow.cn/v1/user/info",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewStringResponse(500, "Internal Server Error"), nil
					},
				)
			},
			giveChannel: mockChannel,
			checkResponse: func(t *testing.T, rsp *SiliconFlowGetInfoListRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				_, isSiliconFlowError := err.(*SiliconFlowErrorRsp)
				assert.False(t, isSiliconFlowError, "error should not be a SiliconFlowErrorRsp")
				assert.Contains(t, err.Error(), "Internal Server Error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()
			tt.setupMock()

			// Use the channel provided by the test case
			provider := NewSiliconFlow(tt.giveChannel)
			rsp, err := provider.GetInfoList()

			tt.checkResponse(t, rsp, err)
		})
	}
}
func TestSiliconFlow_RetrieveFileContent(t *testing.T) {
	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name          string
		giveFileID    string
		mockResponder httpmock.Responder
		checkResponse func(t *testing.T, body io.ReadCloser, err error)
	}{
		{
			name:       "success",
			giveFileID: "file-id-success",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 验证请求头
				assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
				// 返回成功的响应
				return httpmock.NewStringResponse(200, "file content"), nil
			},
			checkResponse: func(t *testing.T, body io.ReadCloser, err error) {
				assert.NoError(t, err)
				require.NotNil(t, body)
				defer body.Close()
				content, readErr := io.ReadAll(body)
				assert.NoError(t, readErr)
				assert.Equal(t, "file content", string(content))
			},
		},
		{
			name:       "api error - structured",
			giveFileID: "file-id-api-error",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 模拟一个结构化的错误响应
				errRsp := SiliconFlowErrorRsp{
					Err: struct {
						Message string `json:"message"`
						Type    string `json:"type"`
						Param   string `json:"param"`
						Code    string `json:"code"`
					}{
						Message: "File not found",
						Type:    "invalid_request_error",
						Code:    "not_found",
					},
				}
				return httpmock.NewJsonResponse(404, errRsp)
			},
			checkResponse: func(t *testing.T, body io.ReadCloser, err error) {
				assert.Nil(t, body)
				require.Error(t, err)
				var siliconFlowError *SiliconFlowErrorRsp
				require.ErrorAs(t, err, &siliconFlowError)
				assert.Equal(t, "File not found", siliconFlowError.Error())
			},
		},
		{
			name:       "api error - unstructured",
			giveFileID: "file-id-unstructured-error",
			mockResponder: func(req *http.Request) (*http.Response, error) {
				// 模拟一个非 JSON 格式的错误响应
				return httpmock.NewStringResponse(500, "Internal Server Error"), nil
			},
			checkResponse: func(t *testing.T, body io.ReadCloser, err error) {
				assert.Nil(t, body)
				require.Error(t, err)
				// 错误不应是 SiliconFlowErrorRsp 类型
				_, isSiliconFlowError := err.(*SiliconFlowErrorRsp)
				assert.False(t, isSiliconFlowError)
				assert.Equal(t, "Internal Server Error", err.Error())
			},
		},
		{
			name:          "network error",
			giveFileID:    "file-id-network-error",
			mockResponder: httpmock.NewErrorResponder(errors.New("connection refused")),
			checkResponse: func(t *testing.T, body io.ReadCloser, err error) {
				assert.Nil(t, body)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "connection refused")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()
			url := "https://api.siliconflow.cn/v1/files/" + tt.giveFileID + "/content"
			httpmock.RegisterResponder("GET", url, tt.mockResponder)

			body, err := p.RetrieveFileContent(tt.giveFileID)

			tt.checkResponse(t, body, err)
		})
	}
}

