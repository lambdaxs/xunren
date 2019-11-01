all:build

build:main.go
	go build -ldflags "-X main.VERSION=2.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X 'main.GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`'" -o xunren_server main.go

move:
	\cp -f xunren_server /root/servers/bin/

start:
	supervisorctl restart xunren_server
