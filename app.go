package hulk

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"git.mysre.cn/ggcp-golang/go-hulk/config"
	"git.mysre.cn/ggcp-golang/go-hulk/logger"
	"git.mysre.cn/ggcp-golang/go-hulk/util"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Application struct {
	Name    string
	Type    int32
	LogPath string

	Config *config.Config
}

type GRPCApplication struct {
	Application
	GRPCServer      *grpc.Server
	GatewayServeMux *runtime.ServeMux
	HTTPServeMux    *http.ServeMux

	isOpenGateway bool
	isSharePort   bool

	RegisterGRPCServer func(*grpc.Server)
	RegisterGateway    func(context.Context, *runtime.ServeMux) error
}

type GinApplication struct {
	Application
	GinEngin *gin.Engine

	RegisterRoute func(*gin.Engine) error
}

// Run 启动并运行一个 gRPC 服务
func (app *GRPCApplication) Run() error {
	if err := app.executeRegisterFunc(); err != nil {
		return err
	}

	errChan := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	logger.Debug(app.Name, "服务启动...")
	go func() {
		if app.isOpenGateway {
			go func() {
				if err := app.runGatewayServer(); err != nil {
					errChan <- fmt.Errorf("Run gatewayServer err: %v", err)
				}
			}()
		}

		if err := app.runGRPCServer(); err != nil {
			errChan <- fmt.Errorf("Run gRPCServer err: %v", err)
		}
	}()

	// gRPC优雅退出
	defer func() {
		app.GRPCServer.GracefulStop()
	}()

	select {
	case err := <-errChan:
		return err
	case <-quit:
	}

	return nil
}

// executeRegisterFunc 执行应用下相关的注册函数
func (app *GRPCApplication) executeRegisterFunc() error {
	if app.RegisterGRPCServer != nil {
		app.RegisterGRPCServer(app.GRPCServer)
	}
	if app.isOpenGateway && app.RegisterGateway != nil {
		if err := app.RegisterGateway(context.Background(), app.GatewayServeMux); err != nil {
			return fmt.Errorf("GRPCApplication.RegisterGateway err: %v", err)
		}
	}
	return nil
}

// runGRPCServer 运行 gRPC Server 端服务
func (app *GRPCApplication) runGRPCServer() error {
	port := app.Config.GetInt64("grpc.port")
	if port == 0 {
		return fmt.Errorf("gRPC 监听端口异常")
	}

	conn, err := net.Listen("tcp", util.GetPortString(port))
	if err != nil {
		return fmt.Errorf("TCP Listen err: %v", err)
	}

	logger.Debug("gRPC API 启动... 监听端口:", port)

	if err := app.GRPCServer.Serve(conn); err != nil {
		return fmt.Errorf("gRPCServer.server 启动异常: %v", err)
	}
	return nil
}

// runGatewayServer 运行用于提供 HTTP 接口服务的 gRPC-Gateway Server 端服务
func (app *GRPCApplication) runGatewayServer() error {
	port := app.Config.GetInt64("http.port")
	if app.isSharePort {
		port = app.Config.GetInt64("grpc.port")
	}
	if port == 0 {
		return fmt.Errorf("gRPC-Gateway 监听端口异常")
	}

	logger.Debug("HTTP API 启动... 监听端口:", port)

	if err := http.ListenAndServe(util.GetPortString(port), app.HTTPServeMux); err != nil {
		return fmt.Errorf("http.Server 启动异常: %v", err)
	}
	return nil
}

func (app *GRPCApplication) OpenGateway()  { app.isOpenGateway = true }
func (app *GRPCApplication) CloseGateway() { app.isOpenGateway = false }

func (app *GRPCApplication) OpenSharePort()  { app.isSharePort = true }
func (app *GRPCApplication) CloseSharePort() { app.isSharePort = false }

func (app *GinApplication) Run() error {
	return nil
}
