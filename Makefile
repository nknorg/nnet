.PHONY: test
test:
	go test -v ./...

.PHONY: pb
pb:
	protoc -I=. --go_out=paths=source_relative:. protobuf/**/*.proto
