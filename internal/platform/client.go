package channel

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/randx"
	"github.com/lazygophers/utils/runtime"
	"io"
	"strings"
	"time"
)

var client = resty.New().
	SetTimeout(time.Minute).
	SetRetryWaitTime(time.Second).
	SetRetryCount(3).
	SetJSONMarshaler(json.Marshal).
	SetJSONUnmarshaler(json.Unmarshal).
	OnPanic(func(request *resty.Request, err error) {
		log.Errorf("panic %s", err)
		runtime.PrintStack()
	}).
	SetHeaders(map[string]string{
		lrpc.HeaderKeepAlive:                "60",
		lrpc.HeaderAccessControlAllowOrigin: "*",
		lrpc.HeaderUserAgent:                fmt.Sprintf("%s/%s", app.Name, app.Version),
	}).
	OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		// 如果 ua 包含 random 的关键词，则随机生成一个 ua
		if strings.Contains(request.Header.Get(lrpc.HeaderUserAgent), "random") {
			request.Header.Set(lrpc.HeaderUserAgent, randx.UserAgentBrowser()+"; random")
		}

		return nil
	}).
	OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		if response.StatusCode() >= 200 && response.StatusCode() < 400 {
			return nil
		}

		return xerror.NewErrorWithMsg(int32(response.StatusCode()), response.Status())
	}).
	SetLogger(log.Clone().SetOutput(io.Discard))

func Load() {

}
