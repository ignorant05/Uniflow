.PHONY: install build test cov-test fmt lint clean 

args = main.go
out = uniflow 
test_args = ./...

install: 
	@echo "Installing dependencies..."
	go mod tidy 
	@echo "Dependencies Installed Successfully"
	@echo "You can proceed now..."

build: 
	@echo "Building Uniflow..."
	go build -o $(out) $(args)  
	@echo "Uniflow Built Successfully"

test: 
	@echo "Running tests..."
	go test $(test_args) -v
	@echo "Tests ran"

cov-test: 
	@echo "Running coverage tests..."
	go test -cover $(test_args) -v
	@echo "Tests ran"

fmt:
	@echo "Running go fmt..."
	go fmt $(test_args) 
	@echo "Formatting completed"

lint:
	@echo "Running go vet..."
	go vet $(test_args)
	@echo "Running staticcheck"
	@which staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed, skipping"
	@echo "Linting completed"

clean: 
	@echo "Cleaning Up..."
	rm $(out) 
	@echo "Old binary cleaned up"

