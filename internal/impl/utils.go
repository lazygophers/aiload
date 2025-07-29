package impl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/utils/network"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"strings"
)

func ToHttpHeader(ctx *fasthttp.RequestCtx) http.Header {
	header := http.Header{}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		header.Set(string(key), string(value))
	})

	return header
}

func ToHeader(ctx *fasthttp.RequestCtx) map[string]string {
	header := map[string]string{}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		header[string(key)] = string(value)
	})
	return header
}

func RealIpFromFasthttp(ctx *fasthttp.RequestCtx) (ip string) {
	ip = RealIpFromHeader(ToHttpHeader(ctx))
	if ip == "" {
		ip = ctx.RemoteIP().String()
	}

	return
}

func RealIpFromHeader(header http.Header) string {
	val := header.Get("Cf-Connecting-Ip")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("Cf-Pseudo-Ipv4")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("Cf-Connecting-Ipv6")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("Cf-Pseudo-Ipv6")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("Fastly-Client-Ip")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("True-Client-Ip")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("X-Real-IP")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("X-Client-IP")
	if val != "" {
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("X-Original-Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if !network.IsLocalIp(v) {
				return v
			}
		}
		if !network.IsLocalIp(val) {
			return val
		}
	}

	val = header.Get("X-Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = header.Get("X-Forwarded")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = header.Get("Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = header.Get("Forwarded")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	return ""
}

func CountryCode(ctx *fasthttp.RequestCtx) string {
	header := ToHttpHeader(ctx)
	if country := header.Get("Cf-Ipcountry"); country != "" {
		return country
	}

	return ""
}

func GetLang(ctx *fiber.Ctx) string {
	lang := ctx.Get("X-Language")

	if lang == "" {
		lang = ctx.Get("Accept-Language")
	}

	if lang == "" {
		lang = state.State.I18n.DefaultLang()
	}

	return lang
}
