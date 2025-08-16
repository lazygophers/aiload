package channel

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAICompatible(t *testing.T) {
	channel := &aiload.ModelChannel{
		Token: "test-token",
	}

	// 测试 NewOpenAICompatible
	p := NewOpenAICompatible(channel)
	require.NotNil(t, p)
	assert.Equal(t, channel, p.channel)

	// 测试 GetRequest
	req := p.GetRequest()
	assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))

	// 激活 httpmock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 为测试设置 BaseURL
	client.SetBaseURL("http://testhost/api/v1")

	t.Run("GetModelList Success", func(t *testing.T) {
		// 模拟成功响应
		mockRsp := &OpenAICompatibleGetModelListRsp{
			Data: []struct {
				ID      string `json:"id"`
				Object  string `json:"object"`
				Created int    `json:"created"`
				OwnedBy string `json:"owned_by"`
			}{
				{ID: "model-1", Object: "model"},
			},
			Object: "list",
		}
		httpmock.RegisterResponder("GET", "http://testhost/api/v1/models",
			httpmock.NewJsonResponderOrPanic(http.StatusOK, mockRsp))

		rsp, err := p.GetModelList()
		require.NoError(t, err)
		require.NotNil(t, rsp)
		assert.Equal(t, "list", rsp.Object)
		assert.Len(t, rsp.Data, 1)
		assert.Equal(t, "model-1", rsp.Data[0].ID)
	})

	t.Run("GetModelList Network Error", func(t *testing.T) {
		// 模拟网络错误
		httpmock.RegisterResponder("GET", "http://testhost/api/v1/models",
			httpmock.NewErrorResponder(errors.New("network error")))

		rsp, err := p.GetModelList()
		require.Error(t, err)
		assert.Nil(t, rsp)
		assert.Contains(t, err.Error(), "network error")
	})
}