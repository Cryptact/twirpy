# Twirpy

> Twirp is a framework for service-to-service communication emphasizing simplicity and minimalism.
> It generates routing and serialization from API definition files and lets you focus on your application's logic
> instead of thinking about folderol like HTTP methods and paths and JSON.
>
> -- <cite>[Twirp's README](https://github.com/twitchtv/twirp/blob/main/README.md)</cite>

Twirpy is a Python implementation of the Twirp framework.
It currently supports [Twirp Wire Protocol v7](https://twitchtv.github.io/twirp/docs/spec_v7.html).

This repository contains:
* a protoc (aka the Protocol Compiler) plugin that generates sever and client code;
* a Python package with common implementation details.

## Installation

### Runtime Library

The runtime library package contains common types like `TwirpServer` and `TwirpClient` that are used by the generated code.

Add the Twirp package to your Python project with:
```
pip install twirpy
```

### Code Generator

You need to install `go` and the `protoc` compiler in your system.
Then, install the `protoc-gen-twirpy` protoc plugin to generate code.

First, install the following prerequisites:
- [Go](https://golang.org/): For installation instructions, see [Go’s documentation](https://golang.org/doc/install).
- [Protocol Buffers](https://developers.google.com/protocol-buffers) compiler: For installation instructions, see [Protocol Buffer Compiler Installation documentation](https://github.com/protocolbuffers/protobuf#protobuf-compiler-installation). You can also use your package manager (e.g. `brew install protobuf` on macOS).

The installed plugin need to be accessible by the protoc compiler.
Set GOBIN (see [go help environment](https://golang.org/cmd/go/#hdr-Environment_variables)) to define where the tool dependencies will be installed.
You might need to add GOBIN to your PATH:
```sh
export GOBIN=$HOME/go/bin
export PATH=$GOBIN:$PATH
```

Then, install the plugin with:
```sh
go install github.com/cryptact/twirpy/protoc-gen-twirpy@latest
```

## Generate and run

Use the protoc plugin to generate twirp server and client code.
```sh
protoc --python_out=. --pyi_out=. --twirpy_out=. example/rpc/haberdasher/service.proto
```

For more information on how to generate code, see the [example](example/README.md).

## Development

We use [`hatch`](https://hatch.pypa.io/latest/) to manage the development process.

To open a shell with the development environment, run: `hatch shell`.
To run the linter, run: `hatch fmt --check` or `hatch fmt` to fix the issues.

## Standing on the shoulders of giants

- The initial version of twirpy was made from an internal copy of https://github.com/daroot/protoc-gen-twirp_python_srv
- The work done by [Verloop](https://verloop.io/) on [the initial versions of Twirpy](https://github.com/verloop/twirpy).
- The `run_in_threadpool` method comes from https://github.com/encode/starlette
