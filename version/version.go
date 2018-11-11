package version

import "fmt"

const (
	// Name is the name of the plugin.
	Name = "vault-plugin-secrets-helloworld"

	// Version is the version of the release.
	Version = "1.0.0"
)

var (
	// GitCommit is the specific git commit of the plugin. This is completed by
	// the compiler.
	GitCommit string

	// HumanVersion is the human-formatted version of the plugin.
	HumanVersion = fmt.Sprintf("%s v%s (%s)", Name, Version, GitCommit)
)
