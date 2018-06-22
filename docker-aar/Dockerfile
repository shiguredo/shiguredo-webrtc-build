FROM ubuntu:16.04

MAINTAINER SHINOHARA, Shunichi <shino@shiguredo.jp>

ADD sources.list /etc/apt/sources.list

RUN apt-get update \
  && apt-get install -y \
         net-tools \
         git \
         curl \
         python \
  && apt-get clean \
  && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

RUN python --version

ENV BUILD_URL_PREFIX=https://github.com/shiguredo/sora-webrtc-build/releases/download/
ENV BUILD_TAG=66.8.2
ENV BUILD_TOOL_VERSION=2.1.0
ENV BUILD_TOOL_FILE=sora-webrtc-build-${BUILD_TOOL_VERSION}-linux-amd64

RUN curl -L -O ${BUILD_URL_PREFIX}/${BUILD_TAG}/${BUILD_TOOL_FILE}.tar.gz \
    && tar zxf ${BUILD_TOOL_FILE}.tar.gz \
    && mv ${BUILD_TOOL_FILE} build

WORKDIR /build

RUN ./webrtc-build fetch

RUN ./webrtc/src/build/install-build-deps.sh

RUN ./webrtc-build build

RUN ./webrtc-build archive
