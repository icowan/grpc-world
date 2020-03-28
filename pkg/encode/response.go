/**
 * @Time: 2020/3/28 15:50
 * @Author: solacowa@gmail.com
 * @File: response
 * @Software: GoLand
 */

package encode

import (
	"context"
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type Failure interface {
	Failed() error
}

type Errorer interface {
	Error() error
}

func err2code(err error) int {
	return http.StatusOK
}

type errorWrapper struct {
	Error string `json:"error"`
}

func Error(ctx context.Context, err error, w http.ResponseWriter) {
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}

func JsonError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func JsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		JsonError(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Error == nil {
		resp.Success = true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}
