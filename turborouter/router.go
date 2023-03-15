package turborouter

import "github.com/ugent-library/deliver/turbo"

type Deserializer[T any] func([]byte) (string, T, error)

type HandlerFunc[T, TT any] func(*turbo.Client[T], TT)

type Router[T, TT any] struct {
	routes       map[string]HandlerFunc[T, TT]
	Deserializer Deserializer[TT]
}

func (r *Router[T, TT]) Respond(c *turbo.Client[T], msg []byte) {
	route, data, err := r.Deserializer(msg)
	// TODO handle error
	if err != nil {
		return
	}
	if handler, ok := r.routes[route]; ok {
		handler(c, data)
	}
}

func (r *Router[T, TT]) Add(route string, handler HandlerFunc[T, TT]) {
	if r.routes == nil {
		r.routes = make(map[string]HandlerFunc[T, TT])
	}
	r.routes[route] = handler
}
