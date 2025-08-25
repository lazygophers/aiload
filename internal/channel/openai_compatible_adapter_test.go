package channel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAiCompatibleAdapter(t *testing.T) {
	channel := &aiload.ModelChannel{}
	adapter := NewOpenAiCompatibleAdapter(channel)
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.client)
}

func TestOpenAiCompatibleAdapter_GetModelList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/models", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(OpenAICompatibleGetModelListRsp{
				Object: "list",
				Data: []struct {
					Id      string `json:"id"`
					Object  string `json:"object"`
					Created int    `json:"created"`
					OwnedBy string `json:"owned_by"`
				}{
					{Id: "model-1", Object: "model", Created: 123, OwnedBy: "test"},
				},
			})
		}))
		defer mockServer.Close()

		channel := &aiload.ModelChannel{BaseUrl: mockServer.URL}
		adapter := NewOpenAiCompatibleAdapter(channel)

		rsp, err := adapter.GetModelList()

		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, "model-1", rsp.Data[0].Id)
	})

	t.Run("API Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer mockServer.Close()

		channel := &aiload.ModelChannel{BaseUrl: mockServer.URL}
		adapter := NewOpenAiCompatibleAdapter(channel)

		rsp, err := adapter.GetModelList()

		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}
