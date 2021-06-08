package hulk

import (
	"git.mysre.cn/yunlian-golang/go-hulk/boot"
	"git.mysre.cn/yunlian-golang/go-hulk/config"
	"git.mysre.cn/yunlian-golang/go-hulk/logger"
)

const (
	HULK_VERSION = "0.1.0"

	APP_TYPE_GRPC = 1
	APP_TYPE_GIN  = 2
)

// NewGRPCApplication 创建一个可以基于 gRPC 框架提供 RPC 接口服务的应用实例
// 基于 gRPC 不仅可以提供 gRPC 默认的 RPC 接口服务，还可以通过配置开启同时提供 HTTP 接口服务
func NewGRPCApplication(name string, opts ...boot.GRPCAppOption) *boot.GRPCApplication {
	app := &boot.GRPCApplication{
		Application: boot.Application{
			Name:    name,
			Type:    APP_TYPE_GRPC,
			LogPath: logger.DefaultLogSavePath,
			Config:  config.NewConfig(),
		},
		GRPCServer:      boot.NewGRPCServer(),
		GatewayServeMux: boot.NewGateway(),
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.Log == nil {
		// 如果未额外设置日志采集器，则使用内置的 ilog 实现
		app.Log = logger.Logger()
	}

	return app
}

// NewGinApplication 创建一个可以基于 Gin 框架提供 HTTP 接口服务的应用实例
func NewGinApplication(name string, opts ...boot.GinAppOption) *boot.GinApplication {
	app := &boot.GinApplication{
		Application: boot.Application{
			Name:    name,
			Type:    APP_TYPE_GIN,
			LogPath: logger.DefaultLogSavePath,
		},
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.Log == nil {
		// 如果未额外设置日志采集器，则使用内置的 ilog 实现
		app.Log = logger.Logger()
	}

	return app
}
