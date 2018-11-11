#!/bin/bash

die() { echo "ERROR: $@" >&2; exit 1; }

which vault &>/dev/null || die "Could not find 'vault' in your PATH"

pluginname=vault-plugin-secrets-helloworld
basedir="$( git rev-parse --show-toplevel )"
pluginbin="$basedir/$pluginname"
workdir="$basedir/_workspace"
plugindir="$workdir/plugins"

addr=127.0.0.1:8282
token=hello

[[ -x $pluginbin ]] || die "Plugin binary was not found. Do you need to run 'make'?"
mkdir -p "$workdir" || die "Could not create working directory at $workdir"
mkdir -p "$plugindir" || die "Could not create plugin directory at $plugindir"
cp "$pluginbin" "$plugindir/$pluginname" || die "Could not copy $pluginname to $plugindir"

_config() {
    sleep 2
    export VAULT_ADDR="http://$addr"
    export VAULT_TOKEN="$token"

    vault plugin register -command="$pluginname" -sha256="$( sha256sum "$plugindir/$pluginname" |cut -d' ' -f1 )" "$pluginname"
    vault secrets enable -path=helloworld -plugin-name="$pluginname" -description="Auto-enabled helloworld plugin for testing" plugin
}

( _config; ) &

vault server -dev -dev-listen-address="$addr" -dev-root-token-id="$token" -dev-plugin-dir="$plugindir" -log-level=debug
