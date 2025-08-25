package channel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
)

func TestNewSiliconFlowAdapter(t *testing.T) {
	channel := &aiload.ModelChannel{}
	adapter := NewSiliconFlowAdapter(channel)
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.client)
}

func TestSiliconFlowAdapter_GetModelList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SiliconFlowGetModelListRsp{
				Object: "list",
				Data: []Model{
					{Id: "model-1", Object: "model", Created: 123, OwnedBy: "test"},
				},
			})
		}))
		defer mockServer.Close()

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", mockServer.URL+"/v1/models",
			httpmock.NewJsonResponderOrPanic(http.StatusOK, SiliconFlowGetModelListRsp{
				Object: "list",
				Data: []Model{
					{Id: "model-1", Object: "model", Created: 123, OwnedBy: "test"},
				},
			}))

		channel := &aiload.ModelChannel{BaseUrl: mockServer.URL}
		adapter := NewSiliconFlowAdapter(channel)

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
		}))
		defer mockServer.Close()

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", mockServer.URL+"/v1/models",
			httpmock.NewStringResponder(http.StatusInternalServerError, ""))

		channel := &aiload.ModelChannel{BaseUrl: mockServer.URL}
		adapter := NewSiliconFlowAdapter(channel)

		rsp, err := adapter.GetModelList()

		assert.Error(t, err)
		assert.Nil(t, rsp)
	})
}
