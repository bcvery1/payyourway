run: build
	cp dist/payyourway .
	./payyourway

build: *.go
	go build -o dist/payyourway .

package: *.go assets/* payyourway-windows-4.0-amd64.exe payyourway-darwin-10.6-amd64
	go build -o payyourway_linux

payyourway-windows-4.0-amd64.exe: *.go assets/*
	xgo --targets=windows/amd64 -ldflags='-H=windowsgui' github.com/bcvery1/payyourway

payyourway-darwin-10.6-amd64: *.go assets/*
	xgo --targets=darwin/amd64 github.com/bcvery1/payyourway

zip: package
	mv payyourway-windows-4.0-amd64.exe payyourway.exe
	mv payyourway-darwin-10.6-amd64 payyourway_mac

	mkdir dist

	zip dist/windows.zip payyourway.exe -r assets
	zip dist/mac.zip payyourway_mac -r assets
	zip dist/linux.zip payyourway_linux -r assets

clean:
	rm -f payyourway-windows-4.0-amd64.exe payyourway.exe payyourway-darwin-10.6-amd64 payyourway_mac payyourway_linux
	rm -rf dist
