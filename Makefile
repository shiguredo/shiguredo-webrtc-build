.PHONY: all dist clean

webrtc-build:
	go build webrtc-build.go

all: webrtc-build

dist:
	./webrtc-build selfdist

clean:
	rm -rf webrtc-build sora-webrtc-build-*
