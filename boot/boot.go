package boot

import (
	"context"
	"net/http"

	"git.mysre.cn/yunlian-golang/go-hulk/config"
	"git.mysre.cn/yunlian-golang/go-hulk/logger"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Application 所有类型应用都需要具备的基础信息
type Application struct {
	Name    string
	Type    int32
	LogPath string

	Config *config.Config
	Log    logger.LogInterface
}

// GRPCApplication 基于 gRPC 实现的 RPC 服务应用类型
// 可以同时提供 RPC + HTTP 接口服务
// 但 HTTP 接口服务强依赖于 protobuf 的 IDL，以及对应的服务实现提供者 gRPC-Gateway
type GRPCApplication struct {
	Application
	GRPCServer      *grpc.Server
	GatewayServeMux *runtime.ServeMux
	HTTPServer      *http.Server

	isOpenGateway bool
	// isSharePort   bool

	RegisterGRPCServer func(*grpc.Server)
	RegisterGateway    func(context.Context, *runtime.ServeMux) error
}

// GinApplication 基于 Gin 实现的 HTTP 服务应用类型
// 可以简单的进行 Route 清单注册对应的 HTTP API 实现
// 无需处理繁琐的 protobuf 定义及相关 IDL 文件资源的维护
type GinApplication struct {
	Application
	GinEngin *gin.Engine

	RegisterRoute func(*gin.Engine) error
}

type GRPCAppOption func(*GRPCApplication)
type GinAppOption func(*GinApplication)

// Init 执行一些应用的初始化动作
// 1. 根据设置好的内容加载对应的配置文件内容
func (app *Application) Init() error {
	// 加载对应的配置文件内容
	app.Config.Load("app.yaml")
	return nil
}

// SetLogger 将应用的日志处理器设置为一个 LogInterface 接口的自定义实现
func (app *Application) SetLogger(lg logger.LogInterface) { app.Log = lg }
