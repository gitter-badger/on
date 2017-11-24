# go-cluster
Go cluster membership disposition; automation, discovery, high availability, policy-based management.

[![Go Report Card](https://goreportcard.com/badge/github.com/continuul/go-cluster)](https://goreportcard.com/report/github.com/continuul/go-cluster)

# Building

```
make
```

# Running a Cluster

Node 1:

```
lsr agent --node test --bind localhost:8765
```

Node 2:

```
lsr agent --bind localhost:9876 --join localhost:8765
```

Node 3:

```
lsr agent --node funk --bind localhost:3433 --join localhost:8765
```
