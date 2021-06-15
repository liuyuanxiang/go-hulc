package boot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/liuyuanxiang/go-hulc/util"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// Run 启动并运行一个 gRPC 服务
func (app *GRPCApplication) Run() error {
	if err := app.Init(); err != nil {
		return fmt.Errorf("gRPC 应用初始化失败 err: %v", err)
	}
	if err := app.executeRegisterFunc(); err != nil {
		return fmt.Errorf("gRPC 执行预加载的注册函数失败 err: %v", err)
	}

	errChan := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	app.Log.Debug(app.Name, "服务启动...")

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

	// 应用优雅退出
	defer app.gracefulStop()

	select {
	case err := <-errChan:
		return err
	case <-quit:
	}

	return nil
}

func (app *GRPCApplication) OpenGateway()  { app.isOpenGateway = true }
func (app *GRPCApplication) CloseGateway() { app.isOpenGateway = false }

// WithGateway 设置开启 gRPC 服务的同时，是否开启 Gateway 额外提供 HTTP 接口服务
func WithGateway(yes bool) GRPCAppOption {
	return func(g *GRPCApplication) {
		if yes {
			g.OpenGateway()
		} else {
			g.CloseGateway()
		}
	}
}

// runGRPCServer 运行 gRPC Server 端服务
func (app *GRPCApplication) runGRPCServer() error {
	port := app.Config.GetInt64("grpc.port")
	if port == 0 {
		return fmt.Errorf("监听端口异常")
	}

	conn, err := net.Listen("tcp", util.GetPortString(port))
	if err != nil {
		return fmt.Errorf("TCP Listen err: %v", err)
	}

	app.Log.Debug("gRPC API 启动... 监听端口:", port)

	if err := app.GRPCServer.Serve(conn); err != nil {
		return fmt.Errorf("gRPCServer.server 启动异常: %v", err)
	}
	return nil
}

// runGatewayServer 运行用于提供 HTTP 接口服务的 gRPC-Gateway Server 端服务
func (app *GRPCApplication) runGatewayServer() error {
	port := app.Config.GetInt64("http.port")
	if port == 0 {
		return fmt.Errorf("监听端口异常")
	}

	app.HTTPServer = &http.Server{
		Addr:    util.GetPortString(port),
		Handler: NewGatewayServerMux(app.GatewayServeMux),
	}

	app.Log.Debug("HTTP API 启动... 监听端口:", port)

	if err := app.HTTPServer.ListenAndServe(); err != nil {
		return fmt.Errorf("http.Server 启动异常: %v", err)
	}
	return nil
}

// gracefulStop 应用优雅退出
// 如果仅开启了 gRPC 服务，则直接关闭 gRPC Server 即可
// 如果同时开启了 HTTP 服务，额外需要处理 HTTP Server 的优雅退出
func (app *GRPCApplication) gracefulStop() {
	if app.HTTPServer != nil {
		// 如果 gRPC 应用同时启动了 HTTP 接口服务，则额外处理 HTTPServer 的优雅退出
		// 设置一个定时上下文，如果 HTTPServer 超过 3s 都没能完全退出在 Server 上的 connection，则输出错误异常
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := app.HTTPServer.Shutdown(ctx); err != nil {
			app.Log.Error("HTTPServer shutdown err:", err)
		}
	}
	app.GRPCServer.GracefulStop()
}

// executeRegisterFunc 执行应用下相关的注册函数
func (app *GRPCApplication) executeRegisterFunc() error {
	if app.RegisterGRPCServer != nil {
		app.RegisterGRPCServer(app.GRPCServer)
	}
	if app.isOpenGateway && app.RegisterGateway != nil {
		if err := app.RegisterGateway(context.Background(), app.GatewayServeMux); err != nil {
			return err
		}
	}
	return nil
}
