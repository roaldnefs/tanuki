<p align="left"><img src="logo/tanuki.png" alt="tanuki" height="160px"></p>

[![Travis CI](https://img.shields.io/travis/roaldnefs/tanuki.svg?style=for-the-badge)](https://travis-ci.org/roaldnefs/tanuki)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/roaldnefs/tanuki)
[![Github All Releases](https://img.shields.io/github/downloads/roaldnefs/tanuki/total.svg?style=for-the-badge)](https://github.com/roaldnefs/tanuki/releases)

Named after the raccoon dog logo of GitLab. A tool for performing actions on GitLab repos or a single repo.

* [Installation](README.md#installation)
     * [Binaries](README.md#binaries)
     * [Via Go](README.md#via-go)
* [Usage](README.md#usage)

## Installation

### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/roaldnefs/tanuki/releases).

### Via Go

```console
$ go get github.com/roaldnefs/tanuki
```

## Usage

```console
$ tanuki -h
A tool for performing actions on GitLab repos or a single repo.

Usage:
  tanuki [command]

Available Commands:
  help        Help about any command
  version     Print the version number of Tanuki

Flags:
      --config string   config file (default is $HOME/.tanuki.yaml)
  -h, --help            help for tanuki
  -t, --toggle          Help message for toggle

Use "tanuki [command] --help" for more information about a command.
```
