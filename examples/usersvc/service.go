package usersvc

import (
	"context"
	"regexp"

	v "github.com/RussellLuo/validating/v3"
	"github.com/RussellLuo/vext"
)

//go:generate protogo validate ./service.go Service

type User struct {
	Name  string
	Age   int
	Email string
}

func (u User) Schema() v.Schema {
	return v.Schema{
		v.F("name", u.Name): v.All(
			v.LenString(0, 10),
			v.Match(regexp.MustCompile(`^\w+$`)),
		),
		v.F("age", u.Age):     v.Range(0, 100),
		v.F("email", u.Email): vext.Email(),
	}
}

type Service interface {
	// CreateUser creates a user with the given attributes.
	//
	// @schema:
	//   user: _
	CreateUser(ctx context.Context, user User) (result User, err error)
}

type UserSvc struct{}

func (us UserSvc) CreateUser(ctx context.Context, user User) (result User, err error) {
	return user, nil
}
