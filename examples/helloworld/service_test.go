package helloworld_test

import (
	"context"
	"fmt"

	"github.com/protogodev/validate/examples/helloworld"
)

func Example() {
	var svc helloworld.Service = &helloworld.Greeter{}
	svc = helloworld.ValidateMiddleware(nil)(svc)

	message, err := svc.SayHello(context.Background(), "Tracey")
	fmt.Printf("message: %q, err: %v\n", message, err)

	message, err = svc.SayHello(context.Background(), "!Tracey")
	fmt.Printf("message: %q, err: %v\n", message, err)

	message, err = svc.SayHello(context.Background(), "Traaaaaaaacey")
	fmt.Printf("message: %q, err: %v\n", message, err)

	// Output:
	// message: "Hello Tracey", err: <nil>
	// message: "", err: name: INVALID(invalid format)
	// message: "", err: name: INVALID(bad length)
}
