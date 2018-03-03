Build
=====

``` bash
docker build \
    -t oarteam/robin_ci:latest \
    -t oarteam/robin_ci:$(date --iso-8601) \
    .
```

Push to Docker Hub
==================

``` bash
docker push oarteam/robin_ci
```
