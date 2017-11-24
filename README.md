# on

[![Join the chat at https://gitter.im/continuul/on](https://badges.gitter.im/continuul/on.svg)](https://gitter.im/continuul/on?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Go cluster management: automation, discovery, high availability, membership, policy-based management.

[![Build Status](https://travis-ci.org/continuul/on.svg?branch=master)](https://travis-ci.org/continuul/on)
[![Go Report Card](https://goreportcard.com/badge/github.com/continuul/on)](https://goreportcard.com/report/github.com/continuul/on)
[![Gitter](https://badges.gitter.im/continuul/on.svg)](https://gitter.im/continuul/on?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=body_badge)

# Building

To build the application for both Linux and Mac:

```
make
```

To build the Docker:

```
make docker
```

# Running a Cluster

## Running a Cluster (Docker)

Start a node:

```
docker run continuul/on:0.1.2
```

## Running a Cluster Locally (No Docker)

Start the primordial node:

Node 1:

```
on agent --node test --bind localhost:8765
```

Subsequent nodes you join the first node:

Node 2:

```
on agent --bind localhost:9876 --join localhost:8765
```

Node 3:

```
on agent --node funk --bind localhost:3433 --join localhost:8765
```
