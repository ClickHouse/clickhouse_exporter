# Contributing notes

## Local setup

The easiest way to make a local development setup is to use Docker Compose.

```
docker-compose up
make init
make
```


## Vendoring

We use [dep](https://github.com/golang/dep) to vendor dependencies. Please use released version, i.e. do not `go get`
from `master` branch.
