/**
 * @Time: 2020/3/28 15:26
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package service

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/icowan/grpc-world/pkg/encode"
	ep "github.com/icowan/grpc-world/pkg/endpoint"
	"github.com/icowan/grpc-world/pkg/grpc"
	"github.com/icowan/grpc-world/pkg/grpc/pb"
	"github.com/icowan/grpc-world/pkg/http"
	"github.com/icowan/grpc-world/pkg/repository"
	"github.com/icowan/grpc-world/pkg/service"
	"github.com/oklog/oklog/pkg/group"
	"golang.org/x/time/rate"
	googleGrpc "google.golang.org/grpc"
	"net"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger log.Logger

const rateBucketNum = 20

var (
	fs       = flag.NewFlagSet("world", flag.ExitOnError)
	httpAddr = fs.String("http-addr", ":8080", "HTTP listen address")
	grpcAddr = fs.String("grpc-addr", ":8081", "gRPC listen address")
)

func Run() {

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	store := repository.New()

	svc := service.New(logger, store)
	svc = service.NewLoggingService(logger, svc)

	ems := []endpoint.Middleware{
		service.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 限流
	}

	eps := ep.NewEndpoint(svc, map[string][]endpoint.Middleware{
		"Get": ems,
		"Put": ems,
	})

	g := &group.Group{}
	initHttpHandler(eps, g)
	initGRPCHandler(eps, g)
	initCancelInterrupt(g)

	_ = level.Error(logger).Log("exit", g.Run())
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initHttpHandler(endpoints ep.Endpoints, g *group.Group) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(level.Error(logger))),
		kithttp.ServerErrorEncoder(encode.JsonError),
	}

	httpHandler := http.MakeHTTPHandler(endpoints, opts...)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		_ = level.Error(logger).Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		_ = level.Error(logger).Log("transport", "HTTP", "addr", *httpAddr)
		return netHttp.Serve(httpListener, httpHandler)
	}, func(error) {
		_ = level.Error(logger).Log("httpListener.Close", httpListener.Close())
	})

}

func initGRPCHandler(endpoints ep.Endpoints, g *group.Group) {
	grpcOpts := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}

	g.Add(func() error {
		baseServer := googleGrpc.NewServer()
		pb.RegisterServiceServer(baseServer, grpc.MakeGRPCHandler(endpoints, grpcOpts...))
		return baseServer.Serve(grpcListener)
	}, func(error) {
		_ = level.Error(logger).Log("grpcListener.Close", grpcListener.Close())
	})
}
