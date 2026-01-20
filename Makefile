.PHONY: install, build

AAR_FILE=mobile.aar

install:
	go mod tidy && go get golang.org/x/mobile/bind && gomobile init

build: 
	rm -rf ./build && mkdir ./build && gomobile bind -target=android -androidapi=21 -o ./build/$(AAR_FILE)
