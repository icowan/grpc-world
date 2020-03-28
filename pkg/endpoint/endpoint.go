/**
 * @Time: 2020/3/28 15:39
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/grpc-world/pkg/encode"
	"github.com/icowan/grpc-world/pkg/service"
)

type GetRequest struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Endpoints struct {
	GetEndpoint endpoint.Endpoint
	PutEndpoint endpoint.Endpoint
}

func NewEndpoint(s service.Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(GetRequest)
			val, err := s.Get(ctx, req.Key)
			return encode.Response{
				Error: err,
				Data:  val,
			}, err
		},
		PutEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(GetRequest)
			err = s.Put(ctx, req.Key, req.Val)
			return encode.Response{
				Error: err,
			}, err
		},
	}

	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}

	return eps
}
