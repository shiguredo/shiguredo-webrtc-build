.PHONY: all dist clean aar copy-aar

IOS_LATEST=m78.8.1

IOS_BUILD_SCRIPT=./scripts/build_all_ios.sh
ANDROID_BUILD_SCRIPT=./scripts/build_all_android.sh

BUILD_IOS_DIR=build/ios
IOS_FRAMEWORK=WebRTC.framework.zip
BUILD_IOS_FRAMEWORK=$(BUILD_IOS_DIR)-$(IOS_LATEST)/$(IOS_FRAMEWORK)

ios-latest: ios-$(IOS_LATEST)
	rm -rf $(BUILD_IOS_DIR)
	mkdir $(BUILD_IOS_DIR)
	cp $(BUILD_IOS_DIR)-$(IOS_LATEST)/$(IOS_FRAMEWORK) $(BUILD_IOS_DIR)

ios-latest-develop: ios-$(IOS_LATEST)-develop
	rm -rf $(BUILD_IOS_DIR)-develop
	mkdir $(BUILD_IOS_DIR)-develop
	cp $(BUILD_IOS_DIR)-$(IOS_LATEST)-develop/$(IOS_FRAMEWORK) $(BUILD_IOS_DIR)-develop

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

AAR_DIR = "android-aar-"$(shell date +%Y%m%d-%H%M%S)

aar-android-%:
	rm -rf docker-aar/config
	mkdir -p docker-aar/config/
	cp -a config/android-$* docker-aar/config/android-aar
	rm -rf docker-aar/scripts
	cp -a scripts docker-aar/scripts
	@echo AAR_VERSION=$(AAR_VERSION)
	rm -f sora-webrtc-$(AAR_VERSION)-android.zip
	docker build --progress=plain --rm -t sora-webrtc-build/docker-aar docker-aar
	$(MAKE) copy-aar

copy-aar:
	(docker rm aar-container > /dev/null 2>&1 ; true)
	docker run --name aar-container sora-webrtc-build/docker-aar /bin/true
	echo "Output dir: " $(AAR_DIR)
	mkdir -p $(AAR_DIR)
	docker cp aar-container:/work/build/android-aar/libwebrtc.aar $(AAR_DIR)/libwebrtc.aar
	docker cp aar-container:/work/build/android-aar/LICENSE.md \
		$(AAR_DIR)/THIRD_PARTY_LICENSES.md
	docker rm aar-container

