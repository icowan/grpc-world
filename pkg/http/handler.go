/**
 * @Time: 2020/3/28 15:52
 * @Author: solacowa@gmail.com
 * @File: http
 * @Software: GoLand
 */

package http

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/grpc-world/pkg/encode"
	ep "github.com/icowan/grpc-world/pkg/endpoint"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHTTPHandler(eps ep.Endpoints, opts ...kithttp.ServerOption) http.Handler {

	r := mux.NewRouter()

	r.Handle("/get/{key}", kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/put/{key}", kithttp.NewServer(
		eps.PutEndpoint,
		decodePutRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		return nil, errors.New("route bad")
	}
	return ep.GetRequest{Key: key}, nil
}

func decodePutRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ep.GetRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		return nil, errors.New("route bad")
	}
	req.Key = key
	return req, nil
}
