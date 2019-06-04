.PHONY: all dist clean aar copy-aar

IOS_BUILD_SCRIPT=./scripts/build_all_ios.sh
ANDROID_BUILD_SCRIPT=./scripts/build_all_android.sh

ios-%-nofetch:
	 $(IOS_BUILD_SCRIPT) --nofetch config/ios-$*

ios-%:
	$(IOS_BUILD_SCRIPT) config/ios-$*

android-%-nofetch:
	 $(ANDROID_BUILD_SCRIPT) --nofetch config/android-$*

android-%:
	$(ANDROID_BUILD_SCRIPT) config/android-$*

webrtc-build:
	go build -o webrtc-build cmd/main.go

all: webrtc-build

dist:
	./webrtc-build selfdist

clean:
	rm -rf webrtc-build sora-webrtc-build-*

AAR_VERSION = $(shell ./webrtc-build -config docker-aar/config.json version | awk '{print $$NF}')

aar-%:
	rm -rf docker-aar/config
	cp -a config/android-$* docker-aar/config
	rm -rf docker-aar/scripts
	cp -a scripts docker-aar/scripts
	@echo AAR_VERSION=$(AAR_VERSION)
	rm -f sora-webrtc-$(AAR_VERSION)-android.zip
	docker build --rm -t sora-webrtc-build/docker-aar docker-aar
	$(MAKE) copy-aar

copy-aar:
	(docker rm aar-container > /dev/null 2>&1 ; true)
	docker run --name aar-container sora-webrtc-build/docker-aar /bin/true
	docker cp aar-container:/work/build/android/libwebrtc.aar .
	docker cp aar-container:/work/build/android/android-release/LICENSE.md \
		THIRD_PARTY_LICENSES.md
	# docker rm aar-container

