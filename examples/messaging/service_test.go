package messaging_test

import (
	"context"
	"fmt"

	"github.com/protogodev/validate/examples/messaging"
)

func Example() {
	var svc messaging.Service = &messaging.Messaging{}
	svc = messaging.ValidateMiddleware(nil)(svc)

	text, err := svc.GetMessage(context.Background(), "123", "00000000-1111-2222-3333-001122334455")
	fmt.Printf("text: %q, err: %v\n", text, err)

	text, err = svc.GetMessage(context.Background(), "", "")
	fmt.Printf("text: %q, err: %v\n", text, err)

	// Output:
	// text: "user[123]: message[00000000-1111-2222-3333-001122334455]", err: <nil>
	// text: "", err: userID: INVALID(has an invalid length), messageID: INVALID(invalid UUID)

}
