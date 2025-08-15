package channel

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestNewOpenAI(t *testing.T) {
	// 表驱动测试
	tests := []struct {
		name    string
		give    *aiload.ModelChannel
		want    *OpenAI
		wantErr bool
	}{
		{
			name: "正常情况",
			give: &aiload.ModelChannel{
				Token: "test-token",
			},
			want: &OpenAI{
				channel: &aiload.ModelChannel{
					Token: "test-token",
				},
			},
			wantErr: false,
		},
		{
			name:    "传入 nil channel",
			give:    nil,
			want:    &OpenAI{channel: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOpenAI(tt.give)
			assert.Equal(t, tt.want, got, "NewOpenAI() 返回的对象与期望不符")
		})
	}
}

func TestOpenAI_GetRequest(t *testing.T) {
	// 准备测试用的 OpenAI 实例
	p := NewOpenAI(&aiload.ModelChannel{
		Token: "test-auth-token",
	})

	// 执行被测试的方法
	req := p.GetRequest()

	// 断言
	require.NotNil(t, req, "GetRequest() 不应返回 nil")
	assert.Equal(t, "test-auth-token", req.Token, "Request 中的 Token 与期望不符")
}

func TestOpenAI_GetModelList(t *testing.T) {
	// 激活 httpmock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	// 准备测试用的 OpenAI 实例
	p := NewOpenAI(&aiload.ModelChannel{
		Token: "fake-token",
	})

	// 定义模拟的成功响应
	mockSuccessResponse := OpenAIGetModelListRsp{
		Object: "list",
		Data: []struct {
			Id      string `json:"id"`
			Object  string `json:"object"`
			Created int    `json:"created"`
			OwnedBy string `json:"owned_by"`
		}{
			{Id: "gpt-3.5-turbo", Object: "model", Created: 1677610602, OwnedBy: "openai"},
			{Id: "gpt-4", Object: "model", Created: 1687882411, OwnedBy: "openai"},
		},
	}
	mockSuccessBody, _ := json.Marshal(mockSuccessResponse)

	// 表驱动测试
	tests := []struct {
		name           string
		mockSetup      func()
		want           *OpenAIGetModelListRsp
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "成功获取模型列表",
			mockSetup: func() {
				httpmock.RegisterResponder("GET", "https://api.openai.com/v1/models",
					httpmock.NewBytesResponder(http.StatusOK, mockSuccessBody))
			},
			want:    &mockSuccessResponse,
			wantErr: false,
		},
		{
			name: "API 返回非 200 状态码",
			mockSetup: func() {
				httpmock.RegisterResponder("GET", "https://api.openai.com/v1/models",
					httpmock.NewStringResponder(http.StatusUnauthorized, `{"error": "unauthorized"}`))
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: "unauthorized",
		},
		{
			name: "网络请求失败",
			mockSetup: func() {
				httpmock.RegisterResponder("GET", "https://api.openai.com/v1/models",
					httpmock.NewErrorResponder(errors.New("network error")))
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: "network error",
		},
		{
			name: "响应 Body 为无效 JSON",
			mockSetup: func() {
				httpmock.RegisterResponder("GET", "https://api.openai.com/v1/models",
					httpmock.NewStringResponder(http.StatusOK, `{"object": "list", "data": [xxx]}`))
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock 状态并应用当前测试用例的 mock 设置
			httpmock.Reset()
			tt.mockSetup()

			// 执行被测试的方法
			got, err := p.GetModelList()

			// 断言错误
			if tt.wantErr {
				require.Error(t, err, "期望出现错误，但实际没有")
				if tt.wantErrMessage != "" {
					assert.Contains(t, err.Error(), tt.wantErrMessage, "错误信息与期望不符")
				}
			} else {
				require.NoError(t, err, "不期望出现错误，但实际发生了")
			}

			// 断言返回值
			assert.Equal(t, tt.want, got, "GetModelList() 的返回值与期望不符")
		})
	}
}
