package objectstores

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

//lint:ignore ST1012 name doesn't start with Err
var Stop = errors.New("stop iterating")

type Store interface {
	Add(context.Context, string, io.Reader) (string, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
	IterateID(context.Context, func(string) error) error
}

type Factory func(string) (Store, error)

var factories = make(map[string]Factory)
var mu sync.RWMutex

func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	factories[name] = factory
}

func New(name, conn string) (Store, error) {
	mu.RLock()
	factory, ok := factories[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown store '%s'", name)
	}
	return factory(conn)
}
