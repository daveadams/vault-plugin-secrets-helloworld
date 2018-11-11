package helloworld

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathConfig defines the config endpoint on the backend.
func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern:         "config",
		HelpSynopsis:    "Configure the greeting message.",
		HelpDescription: "Use this endpoint to configure the salutation word or phrase used to construct a greeting message.",

		Fields: map[string]*framework.FieldSchema{
			"salutation": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The salutation used to build the greeting message for a target.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigUpdate,
		},
	}
}

// pathConfigRead corresponds to READ hello/config and is used to
// read the current configuration.
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	c, err := b.GetConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"salutation": c.Salutation,
		},
	}, nil
}

// pathConfigUpdate corresponds to UPDATE hello/config and is used to update the config.
func (b *backend) pathConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the current configuration
	c, err := b.GetConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Update the configuration
	changed, err := c.Update(d)
	if err != nil {
		return nil, logical.CodedError(400, err.Error())
	}

	// Only do the following if the config is different
	if changed {
		// Generate a new storage entry
		entry, err := logical.StorageEntryJSON("config", c)
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate JSON configuration: {{err}}", err)
		}

		// Save the storage entry
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, errwrap.Wrapf("failed to persist configuration to storage: {{err}}", err)
		}
	}

	return nil, nil
}
