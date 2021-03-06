[![GitHub release](https://img.shields.io/github/release/docwhat/docker-gc.svg)](https://github.com/docwhat/docker-gc/releases)
[![Docker Image Layers](https://images.microbadger.com/badges/image/docwhat/docker-gc.svg)](https://microbadger.com/images/docwhat/docker-gc)
[![GoDoc](https://godoc.org/github.com/docwhat/docker-gc?status.svg)](https://godoc.org/github.com/docwhat/docker-gc)

[![Build Status](https://travis-ci.org/docwhat/docker-gc.svg?branch=master)](https://travis-ci.org/docwhat/docker-gc)
[![Go Report Card](https://goreportcard.com/badge/github.com/docwhat/docker-gc)](https://goreportcard.com/report/github.com/docwhat/docker-gc)
[![Code Coverage](https://codecov.io/gh/docwhat/docker-gc/branch/master/graph/badge.svg)](https://codecov.io/gh/docwhat/docker-gc)
[![GitHub issues](https://img.shields.io/github/issues/docwhat/docker-gc.svg)](https://github.com/docwhat/docker-gc/issues)

Docker GC
=========

> The missing garbage collector for docker

Installation
------------

### Binaries

I have pre-built binaries for several platform already. They are available on the [releases page](https://github.com/docwhat/docker-gc/releases).

### Container

I also have pre-built containers available on [Docker Hub](https://hub.docker.com/r/docwhat/docker-gc/).

You can use this via `docker run`:

``` .sh
$ docker run -d -v /var/run/docker:/var/run/docker --name=gc docwhat/docker-gc:latest
```

### Source

If you have go installed, then you can get the binary `docker-gc` with the following command:

``` .sh
$ go get -u -v github.com/docwhat/docker-gc
```

Usage
-----

    usage: docker-gc_darwin_amd64 [<flags>]

    The missing docker garbage collector.

    Flags:
      -h, --help                Show context-sensitive help (also try --help-long
                                and --help-man).
          --version             Show application version.
      -m, --max-image-age=168h  How old to allow images to be before deletion. (Env:
                                DOCKER_GC_MAX_IMAGE_AGE)
      -s, --sweeper-time=15m    How much time between running checks to delete
                                images. (Env: DOCKER_GC_SWEEPER_TIME)
      -d, --dangle-safe-duration=30m
                                How old should a dangle image be before deletion.
                                (Env: DOCKER_GC_DANGLE_SAFE_DURATION)
      -q, --quiet               Don't show any output. (Env: DOCKER_GC_QUIET)

It uses the normal Docker environment variables, so if `docker info` works, then `docker-gc` should work.

Developers
----------

I love contributions! Read [CONTRIBUTING.md](CONTRIBUTING.md) for details.
