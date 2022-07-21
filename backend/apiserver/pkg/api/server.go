package api

import (
	"context"
	"log"
	"net"
	"net/http"

	"lxxxxxxxx.github.com/applet/backend/pkg/api/open"
	"lxxxxxxxx.github.com/applet/backend/pkg/common"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ServerGroup interface {
	RegisterGrpc(s *grpc.Server)
	RegisterGrpcGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption)
}

type Server struct {
	serverGroups          []ServerGroup
	grpc                  *grpc.Server
	gateway               *http.Server
	grpcUnaryInterceptors []grpc.UnaryServerInterceptor
}

func (s *Server) addServerGroup(server ServerGroup) {
	s.serverGroups = append(s.serverGroups, server)
}

func (s *Server) initGrpc() error {
	s.grpc = grpc.NewServer(
		//grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(make([]grpc.StreamServerIntercepter, 0)...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.grpcUnaryInterceptors...)),
	)
	for _, server := range s.serverGroups {
		server.RegisterGrpc(s.grpc)
	}
	return nil
}

func (s *Server) initGrpcGateway(ctx context.Context) error {
	serverConfig := common.GlobalConfig().Server
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, server := range s.serverGroups {
		server.RegisterGrpcGateway(ctx, mux, serverConfig.GrpcListenAddr, opts)
	}

	var handler http.Handler = mux

	s.gateway = &http.Server{
		Addr:    serverConfig.GatewayListenAddr,
		Handler: handler,
	}
	return nil
}

func (s *Server) RunGrpc(ctx context.Context) error {
	s.initGrpc()
	ch := make(chan error)
	go func() {
		lis, err := net.Listen("tcp", common.GlobalConfig().Server.GrpcListenAddr)
		if err != nil {
			panic(err)
		}

		if err := s.grpc.Serve(lis); err != nil {
			log.Printf("Run grpc failed,error:%w", err)
			ch <- err
		}
	}()
	return <-ch
}

func (s *Server) RunGrpcAndGateway(ctx context.Context) error {
	s.initGrpc()

	if err := s.initGrpcGateway(ctx); err != nil {
		return err
	}
	ch := make(chan error)
	log.Printf("Serve grpc on %s", common.GlobalConfig().Server.GrpcListenAddr)
	go func() {
		lis, err := net.Listen("tcp", common.GlobalConfig().Server.GrpcListenAddr)
		if err != nil {
			log.Printf("Listen grpc failed,err:%w", err)
			ch <- err
		}
		if err := s.grpc.Serve(lis); err != nil {
			log.Printf("Server grpc failed,err:%w", err)
			ch <- err
		}
	}()

	log.Printf("Serve http on %s", common.GlobalConfig().Server.GatewayListenAddr)
	go func() {
		if err := s.gateway.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("start gateway failed,err:%w", err)
			ch <- err
		}
	}()
	return <-ch
}

func NewServer() Server {
	server := Server{}
	server.addServerGroup(&open.ImplementedUserServer{})
	return server
}
