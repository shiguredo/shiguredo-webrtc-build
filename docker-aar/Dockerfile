FROM ubuntu:16.04

MAINTAINER SHINOHARA, Shunichi <shino@shiguredo.jp>

ADD sources.list /etc/apt/sources.list

RUN apt-get update

RUN apt-get install -y \
         lsb-release \
         net-tools \
         git \
         curl \
         python \
         sudo

ADD install-build-deps.sh install-build-deps.sh
RUN yes | ./install-build-deps.sh --no-chromeos-fonts

ENV BUILD_URL_PREFIX=https://github.com/shiguredo/sora-webrtc-build/releases/download/
ENV BUILD_TAG=66.8.2
ENV BUILD_TOOL_VERSION=2.1.0
ENV BUILD_TOOL_FILE=sora-webrtc-build-${BUILD_TOOL_VERSION}-linux-amd64

RUN curl -L -O ${BUILD_URL_PREFIX}/${BUILD_TAG}/${BUILD_TOOL_FILE}.tar.gz \
    && tar zxf ${BUILD_TOOL_FILE}.tar.gz \
    && mv ${BUILD_TOOL_FILE} build

WORKDIR /build

RUN ./webrtc-build fetch

# RUN ./webrtc/src/build/install-build-deps.sh

RUN ./webrtc-build build

RUN ./webrtc-build archive
