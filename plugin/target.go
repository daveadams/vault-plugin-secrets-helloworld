package helloworld

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
)

var (
	ErrTargetNotFound = errors.New("target not found")
)

// Target represents a target from the storage backend.
type Target struct {
	// ID is the identifier of the key in Vault.
	ID string `json:"id"`

	// DisplayName is the name to use when constructing a greeting to this target.
	DisplayName string `json:"display_name"`
}

// GetTarget retrieves the target from the storage backend based on the given ID, or an error
// if it does not exist.
func (b *backend) GetTarget(ctx context.Context, s logical.Storage, id string) (*Target, error) {
	entry, err := s.Get(ctx, "targets/"+id)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to retrieve target %q: {{err}}", id), err)
	}
	if entry == nil {
		return nil, ErrTargetNotFound
	}

	var rv Target
	if err := entry.DecodeJSON(&rv); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to decode entry for %q: {{err}}", id), err)
	}
	return &rv, nil
}

// ListTargets returns the list of target IDs.
func (b *backend) ListTargets(ctx context.Context, s logical.Storage) ([]string, error) {
	entries, err := s.List(ctx, "targets/")
	if err != nil {
		return nil, errwrap.Wrapf("failed to list targets: {{err}}", err)
	}
	return entries, nil
}
