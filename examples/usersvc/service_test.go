package usersvc_test

import (
	"context"
	"fmt"

	"github.com/protogodev/validate/examples/usersvc"
)

func Example() {
	var svc usersvc.Service = &usersvc.UserSvc{}
	svc = usersvc.ValidateMiddleware(nil)(svc)

	created, err := svc.CreateUser(context.Background(), usersvc.User{
		Name:  "Tracey",
		Age:   10,
		Email: "tracey@example.com",
	})
	fmt.Printf("created: %+v, err: %v\n", created, err)

	created, err = svc.CreateUser(context.Background(), usersvc.User{
		Name:  "!Tracey",
		Age:   -1,
		Email: "tracey#example.com",
	})
	fmt.Printf("created: %+v, err: %v\n", created, err)

	// Output:
	// created: {Name:Tracey Age:10 Email:tracey@example.com}, err: <nil>
	// created: {Name: Age:0 Email:}, err: user.name: INVALID(does not match the given regular expression), user.age: INVALID(is not between the given range), user.email: INVALID(invalid email)
}
