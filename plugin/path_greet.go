package helloworld

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathGreet defines the greetings path on the backend.
func (b *backend) pathGreet() *framework.Path {
	return &framework.Path{
		Pattern:         "greet/" + framework.GenericNameRegex("id"),
		HelpSynopsis:    "Provide a secret greeting message to the target identified by 'id'.",
		HelpDescription: "Request this path to get a very friendly message of salutation to the named target.",

		Fields: map[string]*framework.FieldSchema{
			"id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Identifier of the pre-defined target of your greeting.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathGreetRead,
		},
	}
}

// pathGreetRead returns a greeting for the named target.
func (b *backend) pathGreetRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id := d.Get("id").(string)

	t, err := b.GetTarget(ctx, req.Storage, id)
	if err != nil {
		if err == ErrTargetNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	c, err := b.GetConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"message": c.Salutation + ", " + t.DisplayName + "!",
		},
	}, nil
}
