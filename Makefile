MODULE=github.com/grpcmd/grpcmd

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=module=${MODULE} \
	--go-grpc_out=. --go-grpc_opt=module=${MODULE} \
	proto/grpcmd_service.proto

.PHONY: test
test:
	go test

.PHONY: update-test
update-test:
	go test --update
