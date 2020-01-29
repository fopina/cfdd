# cfdd
[![Go Report Card](https://goreportcard.com/badge/github.com/fopina/cfdd)](https://goreportcard.com/report/github.com/fopina/cfdd)
![release](https://github.com/fopina/cfdd/workflows/release/badge.svg)
[![Pulls](https://img.shields.io/docker/pulls/fopina/cfdd.svg)](https://hub.docker.com/r/fopina/cfdd)
[![Layers](https://images.microbadger.com/badges/image/fopina/cfdd.svg)](https://hub.docker.com/r/fopina/cfdd)

Dynamic DNS updater for Cloudflare

cfdd monitors your external IP and updates an A record in your cloudflare account when it changes.

## Instalation

* Use `go get`:

  ```
  go get github.com/fopina/cfdd
  ```

* Download a pre-built binary from [releases](https://github.com/fopina/cfdd/releases).

* Use the [fopina/cfdd](https://hub.docker.com/r/fopina/cfdd) docker image

# Usage

```bash
$ cfdd -h
Usage of cfdd:
  -d, --domain string   Domain (or Zone) that the record is part of
  -e, --email string    Email for authentication
  -h, --help            this
  -p, --polling int     Number of seconds between each check for external IP (use 0 to run only once)
  -r, --record string   Record to be updated
  -t, --token string    API token for authentication - use @/path/to/file to read it from a file instead
  -v, --version         display version
```

You need to get your Global API Token as described [here](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345682)

One time use:

```
cfdd --token <GLOBAL API TOKEN> \
     --record dynsubdomain.yourdomain.tld \
     --domain yourdomain.tld \
     --email 'your@cloudflare.account'
```

If you want to run it continuously use the `--polling` flag specifying the polling interval (number of seconds between IP checks).

Token can also be read from a file instead of the command line, to make it possible to use with docker secrets (and others), as in the follow compose example:

```
version: '3.2'

services:
  web:
    image: fopina/cfdd
    command:
      - --email
      - "your@cloudflare.account"
      - --token
      - "@/run/secrets/cf_api_token"
      - --record
      - dynsubdomain.yourdomain.tld
      - --domain
      - yourdomain.tld
      - --polling
      - "60"
    secrets:
      - cf_api_token

secrets:
  cf_api_token:
    external: true
```
