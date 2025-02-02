# Constructor

Constructor is a small utility CLI tool that generates a constructor for a new struct in Go. The constructor that is
created is heavily opinionated and designed to work with the [functional options pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
to allow for a flexible and expressive API.

As an added bonus, constructor is also able to generate tests for the constructor it creates, with complete coverage.

## Installation

To install Constructor, run the following command:

```bash
go get -u github.com/MadsRC/constructor
```

## Usage

To generate a constructor for a new struct, run the following command:

```bash
constructor --name MyStruct --package mypackage
```

This will print the constructor to the console. If you want to write the constructor to a file, you can use the
`--output` flag:

```bash
constructor --name MyStruct --package mypackage --output mystruct_constructor.go
```

To generate tests for the constructor, use the `--test` flag:

```bash
constructor --name MyStruct --package mypackage --test
```

## Build

This project uses [mise](https://mise.jdx.dev) to manage dev tools, environments and tasks. To build the project, run
the following command:

```bash
mise run build
```

This will result in a binary being created in the `dist/` directory.

Alternatively, you can create a development/debug build by running:

```bash
mise run build --dev
```

This will create a binary with debug information.

## Testing

### Acceptance tests

The acceptance test suite is written as a separate GoLang module in the [tests/acceptance](tests/acceptance) directory.
To run the acceptance tests, you need to have a constructor binary built and available.

To run the acceptance tests, run the following command:

```bash
mise run test:acceptance dist/constructor
```

The arguments passed here is the path to the constructor binary.
