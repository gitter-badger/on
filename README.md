# on
Go cluster membership disposition; automation, discovery, high availability, policy-based management.

[![Build Status](https://travis-ci.org/continuul/on.svg?branch=master)](https://travis-ci.org/continuul/on)
[![Go Report Card](https://goreportcard.com/badge/github.com/continuul/on)](https://goreportcard.com/report/github.com/continuul/on)

# Building

```
make
```

# Running a Cluster

Node 1:

```
on agent --node test --bind localhost:8765
```

Node 2:

```
on agent --bind localhost:9876 --join localhost:8765
```

Node 3:

```
on agent --node funk --bind localhost:3433 --join localhost:8765
```
