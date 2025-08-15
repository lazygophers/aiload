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
	SetLogger(log.Clone().SetOutput(io.Discard))

func Load() {

}
