package helloworld

import (
	"context"

	"github.com/daveadams/vault-plugin-secrets-helloworld/version"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type backend struct {
	*framework.Backend
}

// Factory returns a configured instance of the backend.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	b.Logger().Info("Plugin " + version.HumanVersion + " successfully initialized")
	return b, nil
}

// Backend returns a configured instance of the backend.
func Backend() *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        "The Hello World secrets engine provides greeting messages.",

		Paths: []*framework.Path{
			// path_greet.go
			// ^greet/<target>
			b.pathGreet(),

			// path_targets.go
			// ^targets (LIST)
			b.pathTargets(),
			// ^targets/<target>
			b.pathTargetsCRUD(),

			// path_config.go
			// ^config
			b.pathConfig(),
		},
	}

	return &b
}
