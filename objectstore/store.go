package objectstore

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type Store interface {
	Add(context.Context, string, io.Reader) (string, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
	IterateID(context.Context) (Iter, error)
}

type Iter interface {
	Next() (string, bool)
	Err() error
}

type Factory func(string) (Store, error)

var backends = make(map[string]Factory)
var backendsMu sync.RWMutex

func Register(backend string, factory Factory) {
	backendsMu.Lock()
	defer backendsMu.Unlock()
	backends[backend] = factory
}

func New(backend, conn string) (Store, error) {
	backendsMu.RLock()
	factory, ok := backends[backend]
	backendsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown storage backend '%s'", backend)
	}
	return factory(conn)
}
