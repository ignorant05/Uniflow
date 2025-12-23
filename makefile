args = main.go
out = uniflow 
test_args = ./tests/unit/github/

install: 
	go mod tidy 

build: 
	go build -o $(out) $(args)  

test: 
	go test $(test_args) -v

clean: 
	rm $(out) 

