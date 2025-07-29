package impl

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/core"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils"
	"github.com/lazygophers/utils/json"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
	"strings"
)

func ParseBody(ctx *fiber.Ctx, v any) error {
	switch ctx.Method() {
	case fiber.MethodPost:
		if len(ctx.Body()) == 0 {
			return nil
		}

		contentType := ctx.Get(fiber.HeaderContentType)
		if strings.Contains(contentType, "protobuf") {
			if x, ok := v.(proto.Message); ok {
				err := proto.Unmarshal(ctx.Body(), x)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				log.Debugf("%s %s req body:%v", ctx.Method(), ctx.Path(), v)

				err = utils.Validate(v)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				return nil
			}
		}

		// 默认按照 json 处理
		err := json.Unmarshal(ctx.Body(), v)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = utils.Validate(v)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Debugf("%s %s req body:%v", ctx.Method(), ctx.Path(), v)

		return nil

	default:
		log.Errorf("unsupported method %s", ctx.Method())
	}

	return nil
}

func sendData(ctx *fiber.Ctx, v reflect.Value) error {
	if ctx.Response().IsBodyStream() {
		return nil
	}

	if len(ctx.Response().Body()) > 0 {
		return nil
	}

	var vv any
	if v.IsValid() {
		vv = v.Interface()
	}

	log.Debugf("%s %s resp body:%v", ctx.Method(), ctx.Path(), vv)

	contentType := ctx.Get(fiber.HeaderContentType)
	if strings.Contains(contentType, "protobuf") {
		if x, ok := vv.(proto.Message); ok {
			ctx.Set(fiber.HeaderContentType, "application/protobuf")

			a, err := anypb.New(x)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			buffer, err := proto.Marshal(&core.BaseResponse{
				Data: a,
				Hint: log.GetTrace(),
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return ctx.Send(buffer)
		}
	}

	ctx.Set(fiber.HeaderContentType, "application/json")

	return ctx.JSON(&lrpc.BaseResponse{
		Data: vv,
		Hint: log.GetTrace(),
	})
}

func sendError(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	if x, ok := err.(*xerror.Error); ok {
		contentType := ctx.Get(fiber.HeaderContentType)
		if x.Msg == "" {
			x.Msg = i18n.Localize(fmt.Sprintf("error.%d", x.Code))
		}

		if strings.Contains(contentType, "protobuf") {
			ctx.Set(fiber.HeaderContentType, "application/protobuf")

			buffer, err := proto.Marshal(&core.BaseResponse{
				Code:    x.Code,
				Message: x.Msg,
				Hint:    log.GetTrace(),
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return ctx.Send(buffer)
		}

		ctx.Set(fiber.HeaderContentType, "application/json")

		return ctx.JSON(&lrpc.BaseResponse{
			Code:    x.Code,
			Message: x.Msg,
			Hint:    log.GetTrace(),
		})
	}

	contentType := ctx.Get(fiber.HeaderContentType)
	if strings.Contains(contentType, "application/protobuf") {
		ctx.Set(fiber.HeaderContentType, "application/protobuf")

		buffer, err := proto.Marshal(&core.BaseResponse{
			Code:    -1,
			Message: err.Error(),
			Hint:    log.GetTrace(),
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return ctx.Send(buffer)
	}

	ctx.Set(fiber.HeaderContentType, "application/json")

	return ctx.JSON(&lrpc.BaseResponse{
		Code:    -1,
		Message: err.Error(),
		Hint:    log.GetTrace(),
	})
}

func roleCheck(ctx *fiber.Ctx, pathRole aiload.ModelUser_Role) error {
	// 先通过 token 获取信息
	ctx.Set(fiber.HeaderRetryAfter, "30")

	//info, err := session.Get(ctx)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	if !xerror.CheckCode(err, xerror.ErrNoAuth) {
	//		log.Errorf("err:%v", err)
	//		return err
	//	}
	//
	//	if ctx.Get(lrpc.HeaderToken) != "" {
	//		return err
	//	}
	//
	//	if pathRole != aiload.ModelUser_Role_Public {
	//		ctx.Status(401)
	//		return err
	//	}
	//
	//	info = &session.Info{
	//		Role: aiload.ModelUser_Role_Public,
	//		User: &barbecue.ModelUser{
	//			Type:     aiload.ModelUser_Role_Public,
	//			Language: GetLang(ctx),
	//		},
	//	}
	//}
	//
	//if info.User != nil && info.User.Language != "" {
	//	i18n.SetLanguage(info.User.Language)
	//}
	//
	//if !pathRole.IsInclude(info.Role) {
	//	return xerror.NewError(xerror.ErrNoAuth)
	//}
	//
	//log.Debugf("user rule %s", info.Role)
	//ctx.Context().SetUserValue("info", info)
	//ctx.Context().SetUserValue("role", info.Role)

	return nil
}

func parseRole(role string) aiload.ModelUser_Role {
	switch strings.ToLower(role) {
	case "user", "u", "":
		return aiload.ModelUser_User

	case "p", "public":
		return aiload.ModelUser_Public

	case "admin", "a":
		return aiload.ModelUser_Admin

	default:
		log.Panicf("unsuppost user type %s", role)
	}

	return aiload.ModelUser_User
}

func ToHandler(logic any, role string) fiber.Handler {
	r := parseRole(role)

	lt := reflect.TypeOf(logic)
	lv := reflect.ValueOf(logic)
	if lt.Kind() != reflect.Func {
		panic("parameter is not func")
	}

	// 不管怎么样，第一个参数都是一定要存在的
	x := lt.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Name() != "Ctx" {
		panic("first in is must *github.com/gofiber/fiber/v2/ws.Ctx")
	}

	if x.PkgPath() != "github.com/gofiber/fiber/v2" {
		panic("first in is must *github.com/gofiber/fiber/v2.Ctx")
	}

	// 两个入参，一个出参，那就是需要解析请求参数的
	if lt.NumIn() == 2 && lt.NumOut() == 1 {
		// 先判断出参是否是error
		x = lt.Out(0)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Name() != "error" {
			panic("out is must error")
		}

		// 处理一下入参
		in := lt.In(1)
		for in.Kind() == reflect.Ptr {
			in = in.Elem()
		}

		if in.Kind() != reflect.Struct {
			panic("2rd in is must struct")
		}

		return func(ctx *fiber.Ctx) (err error) {
			err = roleCheck(ctx, r)
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, err)
			}

			req := reflect.New(in)
			err = ParseBody(ctx, req.Interface())
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, xerror.NewInvalidParam(err.Error()))
			}

			err = utils.Validate(req.Interface())
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, xerror.NewInvalidParam(err))
			}

			out := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req})
			if !out[0].IsNil() {
				log.Errorf("err:%v", out[0].Interface())
				return sendError(ctx, out[0].Interface().(error))
			}

			return sendData(ctx, reflect.Value{})
		}
	}

	// 一个入参，两个出参，那就是需要返回数据的
	if lt.NumIn() == 1 && lt.NumOut() == 2 {

		// 先判断第二个出参是否是error
		x = lt.Out(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Name() != "error" {
			panic("out 1 is must error")
		}

		// 处理一下返回
		x = lt.Out(0)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			panic("out 0 is must struct")
		}

		return func(ctx *fiber.Ctx) (err error) {
			err = roleCheck(ctx, r)
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, err)
			}

			out := lv.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if out[1].IsNil() {
				return sendData(ctx, out[0])
			}

			return sendError(ctx, out[1].Interface().(error))
		}
	}

	// 两个入参，两个出参，那就是需要解析请求参数，返回数据的
	if lt.NumIn() == 2 && lt.NumOut() == 2 {
		// 先判断第二个出参是否是error
		x = lt.Out(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Name() != "error" {
			panic("out 1 is must error")
		}

		// 处理一下入参
		in := lt.In(1)
		for in.Kind() == reflect.Ptr {
			in = in.Elem()
		}
		if in.Kind() != reflect.Struct {
			panic("2rd in is must struct")
		}

		// 处理一下返回
		out := lt.Out(0)
		for out.Kind() == reflect.Ptr {
			out = out.Elem()
		}
		if out.Kind() != reflect.Struct {
			panic("out 0 is must struct")
		}

		return func(ctx *fiber.Ctx) (err error) {
			err = roleCheck(ctx, r)
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, err)
			}

			req := reflect.New(in)
			err = ParseBody(ctx, req.Interface())
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, xerror.NewInvalidParam(err.Error()))
			}

			err = utils.Validate(req.Interface())
			if err != nil {
				return sendError(ctx, xerror.NewInvalidParam(err.Error()))
			}

			out := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req})
			if out[1].IsNil() {
				return sendData(ctx, out[0])
			}

			return sendError(ctx, out[1].Interface().(error))
		}
	}

	// 一个入参，一个出参，那就不需要解析请求参数，也不需要返回数据的
	if lt.NumIn() == 1 && lt.NumOut() == 1 {
		// 先判断出参是否是error
		x = lt.Out(0)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Name() != "error" {
			panic("out is must error")
		}

		return func(ctx *fiber.Ctx) (err error) {
			err = roleCheck(ctx, r)
			if err != nil {
				log.Errorf("err:%v", err)
				return sendError(ctx, err)
			}

			out := lv.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if !out[0].IsNil() {
				log.Errorf("err:%v", out[0].Interface())
				return sendError(ctx, out[0].Interface().(error))
			}

			return sendData(ctx, reflect.Value{})
		}
	}

	panic("func is not support")
}

type Route struct {
	Method string
	Path   string

	Handler fiber.Handler

	Role string
}
