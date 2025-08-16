package channel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
)

func setupOpenAITest(t *testing.T) (*OpenAI, func()) {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal server error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIGetModelListRsp{
			Object: "list",
			Data: []struct {
				Id      string `json:"id"`
				Object  string `json:"object"`
				Created int    `json:"created"`
				OwnedBy string `json:"owned_by"`
			}{
				{Id: "gpt-4", Object: "model", Created: 1677610602, OwnedBy: "openai"},
			},
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAICreateChatCompletionRsp{
			Id:      "chatcmpl-123",
			Object:  "chat.completion",
			Created: 1677652288,
			Choices: []struct {
				Index        int               `json:"index"`
				Message      OpenAIChatMessage `json:"message"`
				FinishReason string            `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: OpenAIChatMessage{
						Role:    "assistant",
						Content: "Hello there, how may I assist you today?",
					},
					FinishReason: "stop",
				},
			},
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/completions", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAICreateCompletionRsp{
			Id:      "cmpl-123",
			Object:  "text_completion",
			Created: 1677652288,
			Model:   "gpt-3.5-turbo-instruct",
			Choices: []struct {
				Text         string      `json:"text"`
				Index        int         `json:"index"`
				LogProbs     interface{} `json:"logprobs"`
				FinishReason string      `json:"finish_reason"`
			}{
				{
					Text:         "\n\nThis is a test.",
					Index:        0,
					FinishReason: "length",
				},
			},
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/embeddings", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAICreateEmbeddingRsp{
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
			Model: "text-embedding-ada-002",
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/fine_tuning/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			rsp := OpenAIFineTuningJob{Id: "ftjob-abc123", Status: "succeeded"}
			json.NewEncoder(w).Encode(rsp)
			return
		}
		rsp := OpenAIListFineTuningJobsRsp{
			Object:  "list",
			Data:    []OpenAIFineTuningJobEvent{{Message: "job event"}},
			HasMore: false,
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/fine_tuning/jobs/ftjob-abc123", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIFineTuningJob{Id: "ftjob-abc123", Status: "succeeded"}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/images/generations", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIImageRsp{Created: 1677652288, Data: []struct {
			Url     string `json:"url"`
			B64Json string `json:"b64_json"`
		}{{Url: "https://example.com/image.png"}}}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			rsp := OpenAIFile{Id: "file-uploaded", Object: "file"}
			json.NewEncoder(w).Encode(rsp)
			return
		}
		rsp := OpenAIListFilesRsp{Object: "list", Data: []OpenAIFile{{Id: "file-abc123"}}}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/files/file-abc123", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodDelete {
			rsp := OpenAIDeleteFileRsp{Id: "file-abc123", Object: "file", Deleted: true}
			json.NewEncoder(w).Encode(rsp)
			return
		}
		rsp := OpenAIFile{Id: "file-abc123", Object: "file"}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/files/file-abc123/content", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("file content"))
	})

	mux.HandleFunc("/v1/audio/speech", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("dummy audio content"))
	})

	mux.HandleFunc("/v1/audio/transcriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAICreateTranscriptionRsp{Text: "This is a transcription."}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/audio/translations", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAICreateTranslationRsp{Text: "This is a translation."}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/fine_tuning/jobs/ftjob-abc123/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIFineTuningJob{Id: "ftjob-abc123", Status: "cancelled"}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/fine_tuning/jobs/ftjob-abc123/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIListFineTuningEventsRsp{
			Object: "list",
			Data:   []OpenAIFineTuningJobEvent{{Message: "event message"}},
		}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/images/edits", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIImageRsp{Created: 1677652288, Data: []struct {
			Url     string `json:"url"`
			B64Json string `json:"b64_json"`
		}{{Url: "https://example.com/edited_image.png"}}}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/images/variations", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIImageRsp{Created: 1677652288, Data: []struct {
			Url     string `json:"url"`
			B64Json string `json:"b64_json"`
		}{{Url: "https://example.com/variation_image.png"}}}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/models/gpt-4", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIModel{Id: "gpt-4", Object: "model"}
		json.NewEncoder(w).Encode(rsp)
	})

	mux.HandleFunc("/v1/models/ft:gpt-3.5-turbo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Error") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rsp := OpenAIDeleteModelRsp{Id: "ft:gpt-3.5-turbo", Deleted: true}
		json.NewEncoder(w).Encode(rsp)
	})

	server := httptest.NewServer(mux)

	openAIClient := NewOpenAI(&aiload.ModelChannel{
		Token:   "test-token",
		BaseUrl: server.URL,
	})

	return openAIClient, func() {
		server.Close()
	}
}

func TestOpenAI_GetModelList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		rsp, err := testClient.GetModelList()
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, "gpt-4", rsp.Data[0].Id)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Temporarily modify the client to send the error header
		// This is a bit of a hack. A better way would be to have separate test setups
		// or a more sophisticated mock server.
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "") // Clean up

		rsp, err := testClient.GetModelList()
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateChatCompletion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		req := &OpenAICreateChatCompletionReq{
			Model: "gpt-4",
			Messages: []OpenAIChatMessage{
				{Role: "user", Content: "Hello!"},
			},
		}
		rsp, err := testClient.CreateChatCompletion(req)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "chatcmpl-123", rsp.Id)
		assert.NotEmpty(t, rsp.Choices)
		assert.Equal(t, "Hello there, how may I assist you today?", rsp.Choices[0].Message.Content)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		req := &OpenAICreateChatCompletionReq{Model: "gpt-4"}
		rsp, err := testClient.CreateChatCompletion(req)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateCompletion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		req := &OpenAICreateCompletionReq{
			Model:  "gpt-3.5-turbo-instruct",
			Prompt: "This is a test.",
		}
		rsp, err := testClient.CreateCompletion(req)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "cmpl-123", rsp.Id)
		assert.NotEmpty(t, rsp.Choices)
		assert.Equal(t, "\n\nThis is a test.", rsp.Choices[0].Text)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		req := &OpenAICreateCompletionReq{Model: "gpt-3.5-turbo-instruct"}
		rsp, err := testClient.CreateCompletion(req)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateEmbedding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		req := &OpenAICreateEmbeddingReq{
			Model: "text-embedding-ada-002",
			Input: "The quick brown fox jumps over the lazy dog",
		}
		rsp, err := testClient.CreateEmbedding(req)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, []float64{0.1, 0.2, 0.3}, rsp.Data[0].Embedding)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		req := &OpenAICreateEmbeddingReq{Model: "text-embedding-ada-002"}
		rsp, err := testClient.CreateEmbedding(req)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateSpeech(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		req := &OpenAICreateSpeechReq{
			Model: "tts-1",
			Input: "Hello world",
			Voice: "alloy",
		}
		rsp, err := testClient.CreateSpeech(req)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, http.StatusOK, rsp.StatusCode())
		assert.Equal(t, "dummy audio content", string(rsp.Body()))
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		req := &OpenAICreateSpeechReq{Model: "tts-1"}
		rsp, err := testClient.CreateSpeech(req)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateTranscription(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Create a dummy file for upload
		tmpfile, err := os.CreateTemp("", "test.mp3")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("dummy audio content")
		tmpfile.Close()

		rsp, err := testClient.CreateTranscription(tmpfile.Name(), "whisper-1", "", "", "", 0)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "This is a transcription.", rsp.Text)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		rsp, err := testClient.CreateTranscription("dummy.mp3", "whisper-1", "", "", "", 0)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_CreateTranslation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Create a dummy file for upload
		tmpfile, err := os.CreateTemp("", "test.mp3")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("dummy audio content")
		tmpfile.Close()

		rsp, err := testClient.CreateTranslation(tmpfile.Name(), "whisper-1", "", "", 0)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "This is a translation.", rsp.Text)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		rsp, err := testClient.CreateTranslation("dummy.mp3", "whisper-1", "", "", 0)
		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}

func TestOpenAI_FineTuning(t *testing.T) {
	// Success cases are already grouped in one test, which is reasonable.
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Create
		createReq := &OpenAICreateFineTuningJobReq{
			TrainingFile: "file-abc",
			Model:        "gpt-3.5-turbo",
		}
		createRsp, err := testClient.CreateFineTuningJob(createReq)
		assert.NoError(t, err)
		assert.NotNil(t, createRsp)
		assert.Equal(t, "ftjob-abc123", createRsp.Id)

		// List
		listRsp, err := testClient.ListFineTuningJobs("", 0)
		assert.NoError(t, err)
		assert.NotNil(t, listRsp)
		assert.True(t, len(listRsp.Data) > 0)

		// Retrieve
		retrieveRsp, err := testClient.RetrieveFineTuningJob("ftjob-abc123")
		assert.NoError(t, err)
		assert.NotNil(t, retrieveRsp)
		assert.Equal(t, "ftjob-abc123", retrieveRsp.Id)

		// Cancel
		cancelRsp, err := testClient.CancelFineTuningJob("ftjob-abc123")
		assert.NoError(t, err)
		assert.NotNil(t, cancelRsp)
		assert.Equal(t, "cancelled", cancelRsp.Status)

		// List Events
		eventsRsp, err := testClient.ListFineTuningEvents("ftjob-abc123", "", 0)
		assert.NoError(t, err)
		assert.NotNil(t, eventsRsp)
		assert.True(t, len(eventsRsp.Data) > 0)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		_, err := testClient.CreateFineTuningJob(&OpenAICreateFineTuningJobReq{})
		assert.Error(t, err)

		_, err = testClient.ListFineTuningJobs("", 0)
		assert.Error(t, err)

		_, err = testClient.RetrieveFineTuningJob("ftjob-abc123")
		assert.Error(t, err)

		_, err = testClient.CancelFineTuningJob("ftjob-abc123")
		assert.Error(t, err)

		_, err = testClient.ListFineTuningEvents("ftjob-abc123", "", 0)
		assert.Error(t, err)
	})
}

func TestOpenAI_Images(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Create
		createReq := &OpenAICreateImageReq{
			Prompt: "a white siamese cat",
		}
		createRsp, err := testClient.CreateImage(createReq)
		assert.NoError(t, err)
		assert.NotNil(t, createRsp)
		assert.Len(t, createRsp.Data, 1)
		assert.Equal(t, "https://example.com/image.png", createRsp.Data[0].Url)

		// Create a dummy file for upload
		tmpfile, err := os.CreateTemp("", "test.png")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("dummy image content")
		tmpfile.Close()

		// Edit
		editRsp, err := testClient.CreateImageEdit(tmpfile.Name(), "", "a white siamese cat with a blue hat", "1", "1024x1024", "url", "")
		assert.NoError(t, err)
		assert.NotNil(t, editRsp)
		assert.Equal(t, "https://example.com/edited_image.png", editRsp.Data[0].Url)

		// Variation
		variationRsp, err := testClient.CreateImageVariation(tmpfile.Name(), "1", "1024x1024", "url", "")
		assert.NoError(t, err)
		assert.NotNil(t, variationRsp)
		assert.Equal(t, "https://example.com/variation_image.png", variationRsp.Data[0].Url)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		_, err := testClient.CreateImage(&OpenAICreateImageReq{})
		assert.Error(t, err)

		_, err = testClient.CreateImageEdit("dummy.png", "", "", "", "", "", "")
		assert.Error(t, err)

		_, err = testClient.CreateImageVariation("dummy.png", "", "", "", "")
		assert.Error(t, err)
	})
}

func TestOpenAI_Files(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// List
		listRsp, err := testClient.ListFiles("")
		assert.NoError(t, err)
		assert.NotNil(t, listRsp)
		assert.Equal(t, "list", listRsp.Object)
		assert.Len(t, listRsp.Data, 1)
		assert.Equal(t, "file-abc123", listRsp.Data[0].Id)

		// Create a dummy file for upload
		tmpfile, err := os.CreateTemp("", "test.jsonl")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString(`{"prompt": "Hello", "completion": "World"}`)
		tmpfile.Close()

		// Upload
		uploadRsp, err := testClient.UploadFile(tmpfile.Name(), "fine-tune")
		assert.NoError(t, err)
		assert.NotNil(t, uploadRsp)
		assert.Equal(t, "file-uploaded", uploadRsp.Id)

		// Retrieve
		retrieveRsp, err := testClient.RetrieveFile("file-abc123")
		assert.NoError(t, err)
		assert.NotNil(t, retrieveRsp)
		assert.Equal(t, "file-abc123", retrieveRsp.Id)

		// Retrieve Content
		content, err := testClient.RetrieveFileContent("file-abc123")
		assert.NoError(t, err)
		assert.Equal(t, "file content", content)

		// Delete
		deleteRsp, err := testClient.DeleteFile("file-abc123")
		assert.NoError(t, err)
		assert.NotNil(t, deleteRsp)
		assert.True(t, deleteRsp.Deleted)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		_, err := testClient.ListFiles("")
		assert.Error(t, err)

		_, err = testClient.UploadFile("dummy.jsonl", "fine-tune")
		assert.Error(t, err)

		_, err = testClient.RetrieveFile("file-abc123")
		assert.Error(t, err)

		_, err = testClient.RetrieveFileContent("file-abc123")
		assert.Error(t, err)

		_, err = testClient.DeleteFile("file-abc123")
		assert.Error(t, err)
	})
}

func TestOpenAI_Models(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()

		// Retrieve
		retrieveRsp, err := testClient.RetrieveModel("gpt-4")
		assert.NoError(t, err)
		assert.NotNil(t, retrieveRsp)
		assert.Equal(t, "gpt-4", retrieveRsp.Id)

		// Delete
		deleteRsp, err := testClient.DeleteModel("ft:gpt-3.5-turbo")
		assert.NoError(t, err)
		assert.NotNil(t, deleteRsp)
		assert.True(t, deleteRsp.Deleted)
	})

	t.Run("APIError", func(t *testing.T) {
		testClient, teardown := setupOpenAITest(t)
		defer teardown()
		client.SetHeader("X-Test-Error", "true")
		defer client.SetHeader("X-Test-Error", "")

		_, err := testClient.RetrieveModel("gpt-4")
		assert.Error(t, err)

		_, err = testClient.DeleteModel("ft:gpt-3.5-turbo")
		assert.Error(t, err)
	})
}