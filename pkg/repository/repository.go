/**
 * @Time: 2020/3/28 16:47
 * @Author: solacowa@gmail.com
 * @File: repository
 * @Software: GoLand
 */

package repository

import (
	"errors"
	"sync"
	"time"
)

type Store struct {
	Key       string
	Val       string
	CreatedAt time.Time
}

type StoreKey string

var ErrUnknown = errors.New("unknown store")

type Repository interface {
	Put(key, val string) error
	Get(key string) (res *Store, err error)
}

type store struct {
	mtx    sync.RWMutex
	stores map[StoreKey]*Store
}

func (s *store) Put(key, val string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.stores[StoreKey(key)] = &Store{
		Key:       key,
		Val:       val,
		CreatedAt: time.Now(),
	}
	return nil
}

func (s *store) Get(key string) (res *Store, err error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if val, ok := s.stores[StoreKey(key)]; ok {
		return val, nil
	}
	return nil, ErrUnknown
}

func New() Repository {
	return &store{
		stores: make(map[StoreKey]*Store),
	}
}
