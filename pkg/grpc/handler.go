/**
 * @Time: 2020/3/28 15:57
 * @Author: solacowa@gmail.com
 * @File: grpc
 * @Software: GoLand
 */

package grpc

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/icowan/grpc-world/pkg/encode"
	ep "github.com/icowan/grpc-world/pkg/endpoint"
	"github.com/icowan/grpc-world/pkg/grpc/pb"
)

type grpcServer struct {
	get grpc.Handler
	put grpc.Handler
}

func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.ServiceResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceResponse), nil
}

func (g *grpcServer) Put(ctx context.Context, req *pb.GetRequest) (*pb.ServiceResponse, error) {
	_, rep, err := g.put.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceResponse), nil
}

func MakeGRPCHandler(eps ep.Endpoints, opts ...grpc.ServerOption) pb.ServiceServer {
	return &grpcServer{
		get: grpc.NewServer(
			eps.GetEndpoint,
			decodeGetRequest,
			encodeResponse,
			opts...,
		),
		put: grpc.NewServer(
			eps.PutEndpoint,
			decodeGetRequest,
			encodeResponse,
			opts...,
		),
	}
}

func decodeGetRequest(_ context.Context, r interface{}) (interface{}, error) {
	return ep.GetRequest{
		Key: r.(*pb.GetRequest).Key,
		Val: r.(*pb.GetRequest).Val,
	}, nil
}

func encodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(encode.Response)
	var errStr, data string
	var err error
	if resp.Error != nil {
		errStr = resp.Error.Error()
		err = resp.Error
	}
	if resp.Data != nil {
		data = resp.Data.(string)
	}
	return &pb.ServiceResponse{
		Success: resp.Success,
		Code:    int64(resp.Code),
		Data:    data,
		Err:     errStr,
	}, err
}
