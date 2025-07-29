package api

import (
	"encoding/xml"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lazygophers/aiload/internal/impl"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/i18n"
	app2 "github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/atexit"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/routine"
	"github.com/lazygophers/utils/runtime"
	"github.com/metacubex/mihomo/ntp"
	"github.com/valyala/fasthttp"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

var (
	app = fiber.New(fiber.Config{
		Immutable:            true,
		UnescapePath:         true,
		ReadTimeout:          time.Minute * 5,
		WriteTimeout:         time.Minute * 5,
		IdleTimeout:          time.Minute * 5,
		CompressedFileSuffix: ".gz",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(200).JSON(err.Error())
		},
		DisableStartupMessage:    true, // 打印fiber的监听信息
		AppName:                  app2.Name,
		StreamRequestBody:        true,
		ReduceMemoryUsage:        true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
		XMLEncoder:               xml.Marshal,
		EnableTrustedProxyCheck:  false,
		EnableIPValidation:       false,
		EnablePrintRoutes:        false,
		ColorScheme:              fiber.Colors{},
		EnableSplittingOnParsers: true,
	})

	api fiber.Router
)

func init() {
	app.Use(
		func(ctx *fiber.Ctx) error {
			ctx.Set(lrpc.HeaderTime, strconv.FormatInt(ntp.Now().Unix(), 10))
			return ctx.Next()
		},

		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				log.Errorf("panic:%v", e)
				runtime.PrintStack()
			},
		}),

		// 设置语言
		func(ctx *fiber.Ctx) error {
			i18n.SetLanguage(impl.GetLang(ctx))
			defer i18n.DelLanguage()

			return ctx.Next()
		},

		etag.New(etag.Config{
			Weak: true,
		}),

		cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return true
			},
			AllowCredentials: true,
			MaxAge:           86400,
		}),
	)

	app.Use(
		// 压缩
		compress.New(compress.Config{
			Next: func(c *fiber.Ctx) bool {
				if strings.HasPrefix(c.Path(), "/resource/") {
					return true
				}
				if strings.HasSuffix(c.Path(), ".tar.gz") {
					return true
				}

				return false
			},
			Level: compress.LevelBestCompression,
		}),

		func(ctx *fiber.Ctx) (err error) {
			ctx.Status(200)
			ctx.Set(fiber.HeaderServer, "ice-cream-heaven")

			reqId := ctx.Get(lrpc.HeaderTrance, "api:"+log.GenTraceId())
			log.SetTrace(reqId)
			defer log.DelTrace()

			defer runtime.CachePanicWithHandle(func(e interface{}) {
				log.Errorf("err:%v", e)
				ctx.Status(500)
			})

			ctx.Set(lrpc.HeaderTrance, log.GetTrace())

			addr, err := netip.ParseAddr(impl.RealIpFromFasthttp(ctx.Context()))
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}

			ctx.Context().SetUserValue("ip", addr)
			ctx.Set("X-Client-IP", addr.String())

			err = ctx.Next()

			ctx.Set(lrpc.HeaderTrance, string(log.GetTrace()))
			used := time.Since(ctx.Context().Time())
			if err != nil {
				ctx.Status(500)
				_ = ctx.JSON(map[string]any{
					"code": 500,
					"msg":  err.Error(),
					"hint": log.GetTrace(),
				})
			}

			ctx.Set("X-Used", used.String())

			return nil
		},

		// 保护
		//helmet.New(helmet.Config{
		//	ContentTypeNosniff:        "nosniff",
		//	XFrameOptions:             "SAMEORIGIN",
		//	HSTSPreloadEnabled:        true,
		//	ReferrerPolicy:            "no-referrer",
		//	CrossOriginEmbedderPolicy: "require-corp",
		//	CrossOriginOpenerPolicy:   "same-origin",
		//	CrossOriginResourcePolicy: "same-origin",
		//	OriginAgentCluster:        "?1",
		//	XDNSPrefetchControl:       "off",
		//	XDownloadOptions:          "noopen",
		//	XPermittedCrossDomain:     "none",
		//}),

		//expvarmw.New(),

		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				log.Errorf("panic:%v", e)
				runtime.PrintStack()
			},
		}),
	)

	app.Get("/ip", func(ctx *fiber.Ctx) error {
		return ctx.SendString(impl.RealIpFromFasthttp(ctx.Context()))
	})
	app.All("/generate_204", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(200)
	})

	api = app.Group("/api")
}

func RegisteApi(routes []*Route) {
	for _, route := range routes {
		api.Add(route.Method, route.Path, route.Handler)
	}
}

func Listen() error {
	app.Server().Logger = log.Clone().AppendPrefixMsg("[api]")

	// if state.Config().Api.Debug {
	//	app.Get("/dashboard", monitor.New(monitor.Config{
	//		Title:   "监控",
	//		Refresh: time.Second * 5,
	//	}))
	// }

	//app.Get("/metrics", monitor.New(monitor.Config{
	//	Title:      app2.Name,
	//	Refresh:    time.Second * 5,
	//	APIOnly:    false,
	//	Next:       nil,
	//	CustomHead: "",
	//	FontURL:    "",
	//	ChartJsURL: "",
	//}))

	app.Post("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋")
	})

	app.All("/generate_204", func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderCacheControl, "public,max-age=86400,s-maxage=86400,stale-if-error=86400,stale-while-revalidate=60")
		ctx.Set(fiber.HeaderAge, "86400")

		switch ctx.Method() {
		case fiber.MethodHead:
			return ctx.SendStatus(fasthttp.StatusNoContent)
		default:
			ctx.Status(fasthttp.StatusOK)
			return ctx.SendString("success")
		}
	})

	//app.Get("/ws", func(c *fiber.Ctx) error {
	//	c.Locals("conn_id", c.Context().ConnID())
	//	if websocket.IsWebSocketUpgrade(c) {
	//		return c.Next()
	//	}
	//	return fiber.ErrUpgradeRequired
	//},
	//	websocket.New(impl.Websocket),
	//)

	//app.Get("/sub", cover.HandlerFiber)

	app.Hooks().OnListen(func(info fiber.ListenData) error {
		log.Infof("executor listen on :%s", info.Port)

		for _, route := range app.GetRoutes(true) {
			switch route.Method {
			case fiber.MethodGet, fiber.MethodPost:
				log.Infof("register %s %s", route.Method, route.Path)
			}
		}

		return nil
	})

	app.Hooks().OnShutdown(func() error {
		log.Warn("api exit")
		return nil
	})

	atexit.Register(func() {
		log.Warn("api stop")
		_ = app.Shutdown()
	})

	routine.Go(func() (err error) {
		err = app.Listen(fmt.Sprintf(":%d", 14003))
		if err != nil {
			log.Errorf("err:%v", err)
			runtime.Exit()
			return err
		}

		return nil
	})

	return nil
}

func Shutdown() error {
	return app.Shutdown()
}

type Route struct {
	Method string
	Path   string

	Handler fiber.Handler

	Role string
}
