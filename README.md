[![GitHub
release](https://img.shields.io/github/release/docwhat/docker-gc.svg)](https://github.com/docwhat/docker-gc/releases)
[![Build
Status](https://travis-ci.org/docwhat/docker-gc.svg?branch=master)](https://travis-ci.org/docwhat/docker-gc)
[![GitHub
issues](https://img.shields.io/github/issues/docwhat/docker-gc.svg)](https://github.com/docwhat/docker-gc/issues)
[![Go Report
Card](https://goreportcard.com/badge/github.com/docwhat/docker-gc)](https://goreportcard.com/report/github.com/docwhat/docker-gc)

# Docker GC

> The missing garbage collector for docker

## Installation

### Binaries

I have pre-built binaries for several platform already.  They are available on
the [releases page](https://github.com/docwhat/docker-gc/releases).

### Source

If you have go installed, then you can get the binary `docker-gc`
with the following command:

``` .sh
$ go get -u -v github.com/docwhat/docker-gc
```

Usage
-----

~~~
usage: docker-gc [<flags>]

The missing docker garbage collector.

Flags:
  -h, --help                Show context-sensitive help (also try --help-long
                            and --help-man).
      --version             Show application version.
  -m, --max-image-age=168h  How old to allow images to be before deletion. (Env:
                            DOCKER_GC_MAX_IMAGE_AGE)
  -s, --sweeper-time=15m    How much time between running checks to delete
                            images. (Env: DOCKER_GC_SWEEPER_TIME)
  -v, --verbosity=1         How much logging to stderr. 0 = none. 9 = maximal
                            (Env: DOCKER_GC_VERBOSITY)
~~~

It uses the normal Docker environment variables, so if `docker info` works,
then `docker-gc` should work.

Developers
----------

I use a `Rakefile` to build and test but normal Go commands should work fine.
The `Rakefile` is mainly for convenience, installing linters, and for Travis.

Install [Ruby](https://www.ruby-lang.org/) and you can setup, test and lint the
code.

~~~
$ rake setup test lint
~~~
