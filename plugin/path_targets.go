package helloworld

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathTargets defines the endpoint for listing defined greeting targets on the backend.
func (b *backend) pathTargets() *framework.Path {
	return &framework.Path{
		Pattern:         "targets/?",
		HelpSynopsis:    "List defined greeting targets.",
		HelpDescription: "This endpoint lists all defined greeting targets.",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathTargetsList,
		},
	}
}

// pathTargetsList responds to LIST ./targets requests.
func (b *backend) pathTargetsList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	targets, err := b.ListTargets(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(targets), nil
}

// pathTargetsCRUD defines the endpoint for managing the list of defined targets on the backend.
func (b *backend) pathTargetsCRUD() *framework.Path {
	return &framework.Path{
		Pattern:         "targets/" + framework.GenericNameRegex("id"),
		HelpSynopsis:    "Provide a secret greeting message to the given target.",
		HelpDescription: "Request this path to get a very friendly message of salutation to the named target.",

		Fields: map[string]*framework.FieldSchema{
			"id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Identifier of the pre-defined target of your greeting.",
			},
			"display_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The formal name to use when constructing the greeting message for this target.",
			},
		},

		ExistenceCheck: b.pathTargetsExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathTargetsWrite,  // C
			logical.ReadOperation:   b.pathTargetsRead,   // R
			logical.UpdateOperation: b.pathTargetsWrite,  // U
			logical.DeleteOperation: b.pathTargetsDelete, // D
		},
	}
}

// pathTargetsExistenceCheck is used to check if a given key exists.
func (b *backend) pathTargetsExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	id := d.Get("id").(string)
	if t, err := b.GetTarget(ctx, req.Storage, id); err != nil || t == nil {
		return false, nil
	}
	return true, nil
}

// pathTargetsRead corresponds to GET hello/targets/:id and shows information
// about the target.
func (b *backend) pathTargetsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id := d.Get("id").(string)

	t, err := b.GetTarget(ctx, req.Storage, id)
	if err != nil {
		if err == ErrTargetNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	data := map[string]interface{}{
		"id":           t.ID,
		"display_name": t.DisplayName,
	}

	return &logical.Response{
		Data: data,
	}, nil
}

// pathTargetsWrite corresponds to PUT/POST hello/targets/:id and creates a new target.
func (b *backend) pathTargetsWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id := d.Get("id").(string)
	displayName := d.Get("display_name").(string)

	entry, err := logical.StorageEntryJSON("targets/"+id, &Target{
		ID:          id,
		DisplayName: displayName,
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to create storage entry: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, errwrap.Wrapf("failed to write to storage: {{err}}", err)
	}

	return nil, nil
}

// pathTargetsDelete corresponds to DELETE hello/targets/:id and deletes an existing
// target from the backend storage.
func (b *backend) pathTargetsDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id := d.Get("id").(string)

	_, err := b.GetTarget(ctx, req.Storage, id)
	if err != nil {
		if err == ErrTargetNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	if err := req.Storage.Delete(ctx, "targets/"+id); err != nil {
		return nil, errwrap.Wrapf("failed to delete from storage: {{err}}", err)
	}

	return nil, nil
}
