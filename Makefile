install ::
	asdf plugin add swag https://github.com/behoof4mind/asdf-swag.git
	asdf plugin add golang 
	asdf install

tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

test-i ::
	@go test -v -count=1 -tags=integration ./tests/integration

# if you are using docker rancher 
test-i-rancher ::
	export CGO_CFLAGS="-Wno-gnu-folding-constant" \
	export DOCKER_HOST=unix://$${HOME}/.rd/docker.sock && \
    export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock && \
    export TESTCONTAINERS_HOST_OVERRIDE=$$(rdctl shell ip a show vznat | awk '/inet / {sub("/.*",""); print $$2}')  && \
    echo $${DOCKER_HOST} - $${TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE} - $${TESTCONTAINERS_HOST_OVERRIDE} && \
	go test -v -count=1 -tags=integration ./tests/integration -coverprofile=coverage.out

test-a :: test test-i

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down

swagger ::
	swag fmt && swag init -g cmd/server/main.go --output docs/swagger

coverage:: test
	go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'
	go tool cover -html=coverage.out -o coverage.html
	# open coverage.html