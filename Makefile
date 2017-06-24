.PHONY: clean prepare api api/client build docker/push docker/build run deps

ARCH ?= amd64
GOARCH ?= ${ARCH}
GOARM ?= 7

PKG=github.com/muka/dyndns
NAME=`basename ${PKG}`
GOPKGS=${GOPATH}/src
GOPKGSRC=${GOPKGS}/${PKG}
IMAGE="opny/${NAME}-${ARCH}"
CGO=0

clean:
	rm -rf ./build

prepare:
	mkdir -p ./build

deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

api:
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=google/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. api/api.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. api/api.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. api/api.proto

api/client: api
	cd api && swagger generate client -f api.swagger.json -m api/models -c api/client

build: prepare api
	CGO_ENABLED=${CGO} ARCH=${ARCH} GOARCH=${GOARCH} GOARM=${GOARM} go build -o ./build/${NAME} main.go

run: api
	go run cli/cli.go

docker/build:
	docker build . -t ${IMAGE}

docker/push: docker/build
	docker push ${IMAGE}

all: deps build api/client docker/build
