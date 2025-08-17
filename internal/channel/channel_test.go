package channel

import (
	"reflect"
	"testing"

	"github.com/lazygophers/aiload"
)

func TestNewChannel(t *testing.T) {
	testCases := []struct {
		name      string
		channel   *aiload.ModelChannel
		wantType  reflect.Type
		expectErr bool
	}{
		{
			name: "Siliconflow Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform_Siliconflow,
			},
			wantType:  reflect.TypeOf(&SiliconFlowAdapter{}),
			expectErr: false,
		},
		{
			name: "Gemini Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform_Gemini,
			},
			wantType:  reflect.TypeOf(&GeminiAdapter{}),
			expectErr: false,
		},
		{
			name: "Ollama Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform_Ollama,
			},
			wantType:  reflect.TypeOf(&OllamaAdapter{}),
			expectErr: false,
		},
		{
			name: "OpenAI Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform_OpenAi,
			},
			wantType:  reflect.TypeOf(&OpenAiAdapter{}),
			expectErr: false,
		},
		{
			name: "OpenAI Compatible Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform_OpenAiCompatible,
			},
			wantType:  reflect.TypeOf(&OpenAiCompatibleAdapter{}),
			expectErr: false,
		},
		{
			name: "Unsupported Platform",
			channel: &aiload.ModelChannel{
				Platform: aiload.Platform(-1), // An unsupported platform
			},
			wantType:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter, err := NewChannel(tc.channel)

			if (err != nil) != tc.expectErr {
				t.Errorf("NewChannel() error = %v, expectErr %v", err, tc.expectErr)
				return
			}

			if !tc.expectErr && reflect.TypeOf(adapter) != tc.wantType {
				t.Errorf("NewChannel() = %v, want %v", reflect.TypeOf(adapter), tc.wantType)
			}
		})
	}
}
