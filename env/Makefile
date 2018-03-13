build:
	docker build -t oarteam/robin_ci:latest \
                 -t oarteam/robin_ci:$$(date --iso-8601) \
                 .

build-nocache:
	docker build -t oarteam/robin_ci:latest \
                 -t oarteam/robin_ci:$$(date --iso-8601) \
                 --no-cache \
                 .

push:
	docker push oarteam/robin_ci
