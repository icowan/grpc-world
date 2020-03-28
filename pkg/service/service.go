/**
 * @Time: 2020/3/28 15:27
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/icowan/grpc-world/pkg/repository"
)

type Service interface {
	Get(ctx context.Context, key string) (val string, err error)
	Put(ctx context.Context, key, val string) (err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func (s *service) Put(ctx context.Context, key, val string) (err error) {
	return s.repository.Put(key, val)
}

func (s *service) Get(ctx context.Context, key string) (val string, err error) {
	res, err := s.repository.Get(key)
	if err != nil {
		return
	}
	return res.Val, err
}

func New(logger log.Logger, repository repository.Repository) Service {
	return &service{logger: logger, repository: repository}
}
