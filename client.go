package aiload

import (
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/core"
)

const ServerName = "aiload"

func AddUserAdmin(ctx *lrpc.Ctx, req *AddUserAdminReq) (*AddUserAdminRsp, error) {
	var rsp AddUserAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathAddUserAdmin,
	}, req, &rsp)
}

func GetUserAdmin(ctx *lrpc.Ctx, req *GetUserAdminReq) (*GetUserAdminRsp, error) {
	var rsp GetUserAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathGetUserAdmin,
	}, req, &rsp)
}

func ListUserAdmin(ctx *lrpc.Ctx, req *ListUserAdminReq) (*ListUserAdminRsp, error) {
	var rsp ListUserAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathListUserAdmin,
	}, req, &rsp)
}

func SetUser(ctx *lrpc.Ctx, req *SetUserReq) (*SetUserRsp, error) {
	var rsp SetUserRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathSetUser,
	}, req, &rsp)
}

func SetUserAdmin(ctx *lrpc.Ctx, req *SetUserAdminReq) (*SetUserAdminRsp, error) {
	var rsp SetUserAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathSetUserAdmin,
	}, req, &rsp)
}

func DelUserAdmin(ctx *lrpc.Ctx, req *DelUserAdminReq) (*DelUserAdminRsp, error) {
	var rsp DelUserAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathDelUserAdmin,
	}, req, &rsp)
}

func Login(ctx *lrpc.Ctx, req *LoginReq) (*LoginRsp, error) {
	var rsp LoginRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathLogin,
	}, req, &rsp)
}

func SetChannelAdmin(ctx *lrpc.Ctx, req *SetChannelAdminReq) (*SetChannelAdminRsp, error) {
	var rsp SetChannelAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathSetChannelAdmin,
	}, req, &rsp)
}

func DelChannelAdmin(ctx *lrpc.Ctx, req *DelChannelAdminReq) (*DelChannelAdminRsp, error) {
	var rsp DelChannelAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathDelChannelAdmin,
	}, req, &rsp)
}

func AddChannelAdmin(ctx *lrpc.Ctx, req *AddChannelAdminReq) (*AddChannelAdminRsp, error) {
	var rsp AddChannelAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathAddChannelAdmin,
	}, req, &rsp)
}

func GetChannelAdmin(ctx *lrpc.Ctx, req *GetChannelAdminReq) (*GetChannelAdminRsp, error) {
	var rsp GetChannelAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathGetChannelAdmin,
	}, req, &rsp)
}

func ListChannelAdmin(ctx *lrpc.Ctx, req *ListChannelAdminReq) (*ListChannelAdminRsp, error) {
	var rsp ListChannelAdminRsp
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPathListChannelAdmin,
	}, req, &rsp)
}
