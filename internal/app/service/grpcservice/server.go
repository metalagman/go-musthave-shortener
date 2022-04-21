package grpcservice

import (
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type Server struct {
	server *grpc.Server
}

type ServiceInit func(grpc.ServiceRegistrar)

func New(opts ...grpc.ServerOption) *Server {
	s := &Server{}
	s.server = grpc.NewServer(opts...)
	return s
}

func (s *Server) InitServices(init ...ServiceInit) {
	for _, f := range init {
		f(s.server)
	}
}

func (s *Server) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(
				r.Header.Get("Content-Type"), "application/grpc") {
				s.server.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
