all: ufm

init:
	mkdir -p _output

ufm: init
	go build -o _output/ufm cmd/ufm/main.go

fmt:
	gofmt -w `find cmd pkg -name *.go -type f | xargs`
