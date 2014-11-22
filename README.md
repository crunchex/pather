# pather

A command-line tool to make working with Unix paths easier.

## Features

- duplicates "echo $PATH" functionality
- provides a simple list of path elements
- provides a detailed list of path elements with a best-guess at where they were set

## Installation

Install [go](https://golang.org).
Then run:
```
$ go get github.com/crunchex/pather
```

## Usage

```
Usage of pather:
    -d, --detailed-list=false: use a (detailed) long listing format
    -h, --help: displays pather help
    -l, --list=false: use a long listing format
```

Currently, pather supports Ubuntu and OS X for detailed lists.

## Contact

Github: issues and pull requests welcome!

## License

The MIT License (MIT)

Copyright (c) 2014 Mike Lewis
