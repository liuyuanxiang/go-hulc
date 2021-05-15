package hulk

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	HULK_VERSION = "0.1.0"

	APP_TYPE_GRPC         = 1
	APP_TYPE_GIN          = 2
	APP_TYPE_GRPC_AND_GIN = 3
)

type Application struct {
	Name    string
	Port    int64
	Type    int32
	LogPath string
}

type GRPCApplication struct {
	App        *Application
	GRPCServer *grpc.Server
}

type GinApplication struct {
	App      *Application
	GinEngin *gin.Engine
}

type GRPCAndGinApplication struct {
	App        *Application
	GRPCServer *grpc.Server
	GinEngin   *gin.Engine
}
