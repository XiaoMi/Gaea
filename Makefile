all: build test

build: parser gaea gaea-cc
gaea:
	@bash genver.sh
	go build -o ./bin/gaea ./cmd/gaea
gaea-cc:
	go build -o ./bin/gaea-cc ./cmd/gaea-cc
parser:
	cd parser && make && cd ..
clean:
	@rm -rf bin
	@rm -f .coverage.out .coverage.html
test:
	go test -coverprofile=.coverage.out ./...
	go tool cover -func=.coverage.out -o .coverage.func
	tail -1 .coverage.func
	go tool cover -html=.coverage.out -o .coverage.html
build_with_coverage:
	go test -c cmd/gaea/main.go cmd/gaea/main_test.go -coverpkg ./... -covermode=count -o bin/gaea
