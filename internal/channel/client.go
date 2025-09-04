package channel

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/randx"
	"github.com/lazygophers/utils/runtime"
)

// client 是一个经过预配置的 resty.Client 实例，用于向其他服务发起 HTTP 请求。
// 它被配置为具有默认一分钟的超时、三次重试机制，并使用自定义的 JSON 序列化/反序列化器。
var client = resty.New().
	SetTimeout(time.Minute).                          // 设置请求超时时间
	SetRetryWaitTime(time.Second).                    // 设置重试等待时间
	SetRetryCount(3).                                 // 设置重试次数
	SetJSONMarshaler(json.Marshal).                   // 设置 JSON 序列化器
	SetJSONUnmarshaler(json.Unmarshal).               // 设置 JSON 反序列化器
	OnPanic(func(request *resty.Request, err error) { // 注册 panic 处理器
		log.Errorf("error: %s", err)
		runtime.PrintStack()
	}).
	SetHeaders(map[string]string{ // 设置固定的请求头
		lrpc.HeaderKeepAlive:                "60",
		lrpc.HeaderAccessControlAllowOrigin: "*",
		lrpc.HeaderUserAgent:                fmt.Sprintf("%s/%s", app.Name, app.Version),
	}).
	OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		// 如果 User-Agent 头包含 "random" 关键字，则为其生成一个随机的浏览器 User-Agent。
		// 这在需要模拟不同客户端行为的场景中非常有用。
		if strings.Contains(request.Header.Get(lrpc.HeaderUserAgent), "random") {
			request.Header.Set(lrpc.HeaderUserAgent, randx.UserAgentBrowser()+"; random")
		}
		return nil
	}).
	SetLogger(log.Clone().SetOutput(io.Discard)) // 禁用默认日志记录，以避免不必要的输出
