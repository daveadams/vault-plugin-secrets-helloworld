package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daveadams/vault-plugin-secrets-helloworld/version"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("missing argument")
	}

	switch args[0] {
	case "name":
		fmt.Printf("%s\n", version.Name)
	case "version":
		fmt.Printf("%s\n", version.Version)
	default:
		log.Fatalf("unknown arg %q", args[0])
	}
}
