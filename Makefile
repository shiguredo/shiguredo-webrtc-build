.PHONY: all dist clean aar copy-aar

webrtc-build:
	go build -o webrtc-build cmd/main.go

all: webrtc-build

dist:
	./webrtc-build selfdist

clean:
	rm -rf webrtc-build sora-webrtc-build-*

docker-aar/WebrtcBuildVersion.java:
	@echo "package org.webrtc;" > docker-aar/WebrtcBuildVersion.java
	@echo "public interface WebrtcBuildVersion {" >> docker-aar/WebrtcBuildVersion.java
	@grep -E '"(webrtc_|maint)' docker-aar/config.json | sed \
		-e 's/^ *"/    public static final String /' \
		-e 's/": *"/ = "/' \
        -e 's/",/";/ ' \
		>> docker-aar/WebrtcBuildVersion.java
	@echo "}" >> docker-aar/WebrtcBuildVersion.java

AAR_VERSION = $(shell ./webrtc-build -config docker-aar/config.json version | awk '{print $$NF}')

aar:
	@rm -f docker-aar/WebrtcBuildVersion.java
	$(MAKE) docker-aar/WebrtcBuildVersion.java
	@echo AAR_VERSION=$(AAR_VERSION)
	rm -f sora-webrtc-$(AAR_VERSION)-android.zip
	docker build --rm -t sora-webrtc-build/docker-aar docker-aar
	$(MAKE) copy-aar

copy-aar:
	(docker rm aar-container > /dev/null 2>&1 ; true)
	docker run --name aar-container sora-webrtc-build/docker-aar /bin/true
	docker cp aar-container:/build/sora-webrtc-$(AAR_VERSION)-android.zip .
	docker rm aar-container
