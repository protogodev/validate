package customvalidator

import (
	"regexp"

	v "github.com/RussellLuo/validating/v3"
)

var reUUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func UUID() *v.MessageValidator {
	return v.Is(reUUID.MatchString).Msg("invalid UUID")
}
