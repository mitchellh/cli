# Go CLI Library [![GoDoc](https://godoc.org/github.com/mitchellh/cli?status.png)](https://godoc.org/github.com/mitchellh/cli)
> A Go library for implementing powerful command-line interfaces

[![GitHub tag](https://img.shields.io/github/tag/mitchellh/cli?include_prereleases=&sort=semver)](https://github.com/mitchellh/cli/releases/)
[![License](https://img.shields.io/badge/License-Mozilla_Public_License_2.0-blue)](#license)

_cli_ is the library that powers the CLI for:
[Packer](https://github.com/mitchellh/packer),
[Serf](https://github.com/hashicorp/serf),
[Consul](https://github.com/hashicorp/consul),
[Vault](https://github.com/hashicorp/vault),
[Terraform](https://github.com/hashicorp/terraform), and
[Nomad](https://github.com/hashicorp/nomad).

## Features

* Easy sub-command based CLIs: `cli foo`, `cli bar`, etc.
* Support for nested subcommands such as `cli foo bar`.
* Optional support for default subcommands so `cli` does something
  other than error.
* Support for shell autocompletion of subcommands, flags, and arguments
  with callbacks in Go. You don't need to write any shell code.
* Automatic help generation for listing subcommands.
* Automatic help flag recognition of `-h`, `--help`, etc.
* Automatic version flag recognition of `-v`, `--version`.
* Helpers for interacting with the terminal, such as outputting information,
  asking for input, etc. These are optional, you can always interact with the
  terminal however you choose.
* Use of Go interfaces/types makes augmenting various parts of the library a
  piece of cake.

## Example

A simple example of creating and running a CLI app:

```go
package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"foo": fooCommandFactory,
		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
```

## License

Released under [Mozilla Public License 2.0](/LICENSE) by [@mitchellh](https://github.com/mitchellh).
