run: build
	cp dist/payyourway .
	./payyourway

build: *.go
	go build -o dist/payyourway .
