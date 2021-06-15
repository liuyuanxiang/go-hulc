package boot

import (
	"context"
	"encoding/json"
	"net/http"

	"git.mysre.cn/liuyx02/go-hulc/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func NewGateway() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithErrorHandler(customHTTPError),
		runtime.WithForwardResponseOption(cors),
	)
}

func NewGatewayServerMux(gateway *runtime.ServeMux) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", gateway)
	return mux
}

type httpErrorResponse struct {
	ErrCode   int64       `json:"errcode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrDetail string      `json:"errDetail,omitempty"`
}

func customHTTPError(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	response := &httpErrorResponse{
		ErrCode: 10000,
		Message: "Unkonwn",
	}

	if s.Message() != "" {
		response.Message = s.Message()
		response.ErrDetail = s.Message()
		// if errors.As(err, &errcode.AppErr) {
		// 	response.Message = s.Message()
		// }
		// if config.Env == "test" || config.Env == "dev" {
		// 	response.ErrDetail = s.Message()
		// }
		logger.Error("gRPC-Gateway http err:", s.Message())
	}

	jsonMsg, _ := json.Marshal(response)
	w.Header().Set("Content-Type", marshaler.ContentType(s.Proto()))
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	if _, err = w.Write(jsonMsg); err != nil {
		logger.Error("gRPC-Gateway response write err:", err, s.Message())
	}
}

// HTTP 接口服务跨域处理
func cors(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	return nil
}
