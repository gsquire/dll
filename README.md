# dll

**D**efer **L**oop **L**inter

[![Build Status](https://travis-ci.org/gsquire/dll.svg?branch=master)](https://travis-ci.org/gsquire/dll)

A simple linter to find `defer` statements inside of for loops in Go source.

## Why?
It's often erroneous to use `defer` inside of a loop as it can lead to memory leaks or other
unintended behavior. It can also be easy to miss this in a code review as using `defer` to
close sockets or files is a common Go idiom. This tool aims to point these out by simply printing
the line of a `defer` statement when it is found inside of a loop.

## Install

```sh
go get github.com/gsquire/dll
```

## Usage

```sh
dll source.go

dll *.go
```

## Contributing
Found a bug? Found a case this didn't catch? Great! Feel free to open an issue or add a test case!

## License
MIT
