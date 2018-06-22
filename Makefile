.PHONY: all dist clean

webrtc-build:
	go build -o webrtc-build cmd/main.go

all: webrtc-build

dist:
	./webrtc-build selfdist

clean:
	rm -rf webrtc-build sora-webrtc-build-*

aar:
	docker build --rm -t sora-webrtc-build/docker-aar docker-aar
