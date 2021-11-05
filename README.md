# Probability

[![CI](https://github.com/spenserblack/probability/actions/workflows/ci.yml/badge.svg)](https://github.com/spenserblack/probability/actions/workflows/ci.yml)

Collections where values are returned based on probability

See [examples](./examples_test.go) for example usage.

## Markov

This module comes with a binary for generating random strings of text using Markov chains.
The preferred way to install it is by cloning this repository and running `make install`,
but `go install github.com/spenserblack/probability/cmd/markov` works too.

### Example Usage

```console
$ markov sentence "foo bar" "bar baz" "bar foo"
foo bar foo bar baz
```

The above is just example output. No guarantees that it will happen
(although it should after enough tries).
