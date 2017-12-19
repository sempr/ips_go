all: golang

golang:
	@echo "--> Generating Go files"
	protoc -I protos/ --go_out=plugins=grpc:protos/ protos/ipsvc.proto
	@echo ""

run:
	go build -o main github.com/sempr/ips_go/cmds
