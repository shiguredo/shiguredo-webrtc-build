.PHONY: all dist clean aar

AAR_VERSION = $(shell ./webrtc-build version | awk '{print $$NF}')

webrtc-build:
	go build -o webrtc-build cmd/main.go

all: webrtc-build

dist:
	./webrtc-build selfdist

clean:
	rm -rf webrtc-build sora-webrtc-build-*

docker-aar/WebrtcBuildVersion.java:
	@echo "package org.webrtc;" > docker-aar/WebrtcBuildVersion.java
	@echo "public interface WebrtcBuildVersion {" > docker-aar/WebrtcBuildVersion.java
	@grep '"webrtc_' docker-aar/config.json | sed \
		-e 's/^ *"/    public static final String /' \
		-e 's/": *"/ = "/' \
        -e 's/",/";/ ' \
		>> docker-aar/WebrtcBuildVersion.java
	@echo "}" >> docker-aar/WebrtcBuildVersion.java

aar: docker-aar/WebrtcBuildVersion.java
	@echo AAR_VERSION=$(AAR_VERSION)
	rm -f sora-webrtc-$(AAR_VERSION)-android.zip
	docker build --rm -t sora-webrtc-build/docker-aar docker-aar
	(docker rm aar-container > /dev/null 2>&1 ; true)
	docker run --name aar-container sora-webrtc-build/docker-aar /bin/true
	docker cp aar-container:/build/sora-webrtc-$(AAR_VERSION)-android.zip .
	docker rm aar-container
