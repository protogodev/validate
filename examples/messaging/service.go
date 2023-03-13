package messaging

import (
	"context"
	"fmt"
)

//go:generate protogo validate --custom=./decl.go ./service.go Service

type Service interface {
	// GetMessage get the specified message.
	//
	// @schema:
	//   userID: len(1, 10)
	//   messageID: uuid
	GetMessage(ctx context.Context, userID string, messageID string) (text string, err error)
}

type Messaging struct{}

func (m *Messaging) GetMessage(ctx context.Context, userID string, messageID string) (string, error) {
	return fmt.Sprintf("user[%s]: message[%s]", userID, messageID), nil
}
