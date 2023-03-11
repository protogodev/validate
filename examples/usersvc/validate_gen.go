// Code generated by validate; DO NOT EDIT.
// github.com/protogodev/validate

package usersvc

import (
	"context"

	v "github.com/RussellLuo/validating/v3"
)

func ValidateMiddleware(wrap func(error) error) func(Service) Service {
	return func(next Service) Service {
		if wrap == nil {
			wrap = func(err error) error { return err }
		}
		return validateMiddleware{
			next: next,
			wrap: wrap,
		}
	}
}

type validateMiddleware struct {
	next Service
	wrap func(error) error
}

func (mw validateMiddleware) CreateUser(ctx context.Context, user User) (User, error) {
	schema := v.Schema{
		v.F("user", user): user.Schema(),
	}

	if err := v.Validate(schema); err != nil {
		return User{}, mw.wrap(err)
	}

	return mw.next.CreateUser(ctx, user)
}
