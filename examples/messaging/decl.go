package messaging

import (
	"github.com/protogodev/validate/examples/messaging/customvalidator"
)

var _ = []any{
	// type=string args=0
	customvalidator.UUID,
}
