package helloworld

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
)

// errImmutable is a logical coded error that is returned when the user tries to
// modify an immutable field.
func errImmutable(s string) error {
	return logical.CodedError(400, fmt.Sprintf("cannot change %s after key creation", s))
}
