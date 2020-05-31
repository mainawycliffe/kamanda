---
title: "Kamanda Installation Guide"
slug: kamanda_login:ci
url: /installation
summary: "How to install kamanda in to your system"
---

## Installation

> This is still an early version, feedback on improvements is much needed.

### Using Shell

#### Bash

```sh
curl https://raw.githubusercontent.com/mainawycliffe/kamanda/master/install.sh | sh
```

#### Powershell

Coming soon.

### Using npm (Node Package Manager)

You can use npm to install Kamanda:

```sh
npm -g install kamanda
```

### Using Yan

You can also use yarn to install Kamanda:

```sh
yarn global add kamanda
```

To check if installation was completed successfully, run the following
command:

```sh
kamanda version
```

You can view all supported commands [here](./docs/kamanda.md) or by running `kamanda help`

```sh
kamanda help
```

### Executable Binaries

You can find the latest binaries for your operating system in the
[releases](https://github.com/mainawycliffe/kamanda/releases).

## Building from the Source

You can build your version of kamanda from the source. Kamanda is written using
Golang and all you will need is Golang to build kamanda from the source.

Follow the following stepts to build kamanda.

- First, makes sure you have Golang installed.
- Clone this repository: `git clone https://github.com/mainawycliffe/kamanda.git`
- Install all go packages: `go get ./...`
- Build kamanda using Golang: `go build -o kamanda`

More installation avenues are coming soon.
