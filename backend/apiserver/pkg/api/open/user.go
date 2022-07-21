package open

import (
	"context"
	"log"
	pb "lxxxxxxxx.github.com/applet/backend/proto/apiserver/public/open"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials/insecure"
)

type ImplementedUserServer struct {
	pb.UnimplementedUserServer
}

func (*ImplementedUserServer) RegisterGrpc(s *grpc.Server) {
	pb.RegisterUserServer(s, &ImplementedUserServer{})
}

func (*ImplementedUserServer) RegisterGrpcGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	err := pb.RegisterUserHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (*ImplementedUserServer) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRes, error) {
	log.Printf("reveive create user request,usernaem:%s,sex:%s", req.UserName, req.UserSex)
	return &pb.CreateUserRes{UserId: "111"}, nil
}
