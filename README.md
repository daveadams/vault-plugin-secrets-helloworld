# vault-plugin-secrets-helloworld

This is a plugin backend for [HashiCorp Vault][vault] that provides a minimal
and trivial but relatively complete example of a secrets engine plugin to serve
as a basic model for development of actual functional plugins.

This code was patterned heavily after the official [gcpkms secrets plugin][gcpkms].

For more details about how Vault plugins work, see the [plugin architecture documentation][pluginarch].

## Building

Running `make` or `go build ./cmd/vault-plugin-secrets-helloworld` in the root
directory of this project should generate a `vault-plugin-secrets-helloworld`
binary.

## Registration

To register this plugin with Vault, first copy the binary to the plugin directory
configured for your running instance of Vault, then register the plugin with a
command similar to this:

    $ vault plugin register \
          -command=vault-plugin-secrets-helloworld \
          -sha256=$( sha256sum vault-plugin-secrets-helloworld |cut -d' ' -f1 ) \
          helloworld-plugin

See the [plugin registration docs][plugindocs] for more details.

## Usage

Once the plugin is registered as above, you can enable it on a given path:

    $ vault secrets enable \
          -path=helloworld \
          -plugin-name=helloworld-plugin \
          -description="Example of the helloworld plugin" \
          plugin

Then you can configure a target:

    $ vault write helloworld/targets/everyone display_name="all y'all"

And configure a salutation:

    $ vault write helloworld/config salutation=Howdy

And finally read a greeting:

    $ vault read -field message helloworld/greet/everyone
    Howdy, all y'all!

Use `vault path-help helloworld` to see full documentation on the options
available on each endpoint.

## Testing

You can interactively test this plugin in development mode using the scripts
in the `test` directory:

Run `make test` in one shell window. This will start Vault in dev mode and
register and enable the plugin at `helloworld/`.

Then in another shell window, run `source test/env.sh`. This will set `VAULT_ADDR`
and `VAULT_TOKEN` appropriately to interact with the test server. Then, you can
try such commands as:

    $ vault path-help helloworld

    $ vault write helloworld/targets/france display_name="tout le monde"

    $ vault write helloworld/config salutation=Bonjour

    $ vault read helloworld/greet/france

## License

This code is licensed under the Mozilla Public License (MPL) 2.0. See the
LICENSE file for more information.

[vault]: https://www.vaultproject.io
[gcpkms]: https://github.com/hashicorp/vault-plugin-secrets-gcpkms
[plugindocs]: https://www.vaultproject.io/docs/plugin/index.html
[pluginarch]: https://www.vaultproject.io/docs/internals/plugins.html
