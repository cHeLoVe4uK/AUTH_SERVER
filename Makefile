run:
	go run cmd/api/main.go

lint:
	golangci-lint run --disable-all -E unused -E gofumpt -E govet -E errcheck

fix:
	golangci-lint run --disable-all -E gofumpt --fix
