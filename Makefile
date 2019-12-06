all: test build 

build: 
	go build  -o bin/conget -v  
test:
	go test -v ./...
clean:
	rm bin/conget 
	rm bin/windows_amd64/conget
	rm bin/darwin_amd64/conget		

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64/conget -v 

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/conget -v 	
    	