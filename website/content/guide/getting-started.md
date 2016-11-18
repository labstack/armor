+++
title = "Getting Started"
description = "Getting started with Armor"
[menu.side]
  name = "Getting Started"
  parent = "guide"
  weight = 1
+++

## Installation

### Binary

- Download the latest armor release for your platform from https://github.com/labstack/armor/releases
- Copy the armor binary to somewhere on the `PATH` so that it can be executed e.g. `/usr/local/bin` or `%PATH%`

### Homebrew

`brew install armor`

### Go

`go get -u github.com/labstack/armor/cmd/armor`

### Docker

`docker pull labstack/armor`

## Usage

Type `armor` in your terminal

```sh
❯ armor

 _______  ______    __   __  _______  ______
|   _   ||    _ |  |  |_|  ||       ||    _ |
|  |_|  ||   | ||  |       ||   _   ||   | ||
|       ||   |_||_ |       ||  | |  ||   |_||_
|       ||    __  ||       ||  |_|  ||    __  |
|   _   ||   |  | || ||_|| ||       ||   |  | |
|__| |__||___|  |_||_|   |_||_______||___|  |_|

https://armor.labstack.com               v0.2.4

Uncomplicated HTTP server, supports HTTP/2 and
auto TLS
_______________O/______________________________
               O\

 ⇛ http server started on :8080
```

This starts armor on address `:8080`, serving the current directory listing using
the default config. Go to http://localhost:8080 to browse the directory.

Armor can also be run as in a Docker:

```sh
docker run \
  -p 8080:80 \
  -v <config_file>:/etc/armor/config.json \
  -v <volume_to_mount>:/var/www \
    labstack/armor -c /etc/armor/config.json
```

### Command-line Flags

- `-p` http listen port
- `-c` config file
- `-v` print the version

### [Configuration]({{< ref "guide/configuration.md">}})
