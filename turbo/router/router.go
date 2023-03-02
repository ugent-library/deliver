package router

import "github.com/ugent-library/deliver/turbo"

type Deserializer[T any] func([]byte) (string, T, error)

type HandlerFunc[T any] func(*turbo.Client, T)

type Router[T any] struct {
	routes       map[string]HandlerFunc[T]
	Deserializer Deserializer[T]
}

func (r *Router[T]) Respond(c *turbo.Client, msg []byte) {
	route, data, err := r.Deserializer(msg)
	// TODO handle error
	if err != nil {
		return
	}
	if handler, ok := r.routes[route]; ok {
		handler(c, data)
	}
}

func (r *Router[T]) Add(route string, handler HandlerFunc[T]) {
	r.routes[route] = handler
}
