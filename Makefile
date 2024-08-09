test:
	go test -v ./...
lint:
	gocritic check ./...
	revive ./...
	golint ./...
	goconst ./...
	golangci-lint run
	go vet ./...
	staticcheck ./...