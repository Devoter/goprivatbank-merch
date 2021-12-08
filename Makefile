NAME=app
GIT_COMMIT=$(shell git describe --tags --long --always --dirty)
BUILD_TIMESTAMP=$(shell date +%s)
LDFLAGS=-ldflags "-X main.Version=${GIT_COMMIT} -X main.BuildTimestamp=${BUILD_TIMESTAMP}"
ALLBIN=${NAME}
CGO_ENABLED=0
GOOS=
GOARCH=
GOARM=
INSTALL_PREFIX=/usr/local/bin

all: ${ALLBIN}

clean:
	rm --force ${ALLBIN}

clean-${NAME}:
	rm --force app

${NAME}:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${NAME} ./cmd/app/...

#generate:
#

test:
	go test ./...

coverage:
	go test ./... --coverprofile=coverage.out 
	go tool cover --html=coverage.out

install:
	mkdir -p ${INSTALL_PREFIX}
	cp ${NAME} ${INSTALL_PREFIX}

uninstall:
	rm --force ${NAME} ${INSTALL_PREFIX}
	rmdir --ignore-fail-on-non-empty ${INSTALL_PREFIX}
