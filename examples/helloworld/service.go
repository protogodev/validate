package helloworld

import (
	"context"
)

//go:generate protogo validate ./service.go Service

// Service is used for saying hello.
type Service interface {
	// SayHello says hello to the given name.
	//
	// @schema:
	//   name: len(0, 10).msg("bad length") && match(`^\w+$`).msg("invalid format")
	SayHello(ctx context.Context, name string) (message string, err error)
}

type Greeter struct{}

func (g *Greeter) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
