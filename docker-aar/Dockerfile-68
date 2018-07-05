FROM ubuntu:18.04

MAINTAINER SHINOHARA, Shunichi <shino@shiguredo.jp>

ADD sources.list /etc/apt/sources.list

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update

RUN apt-get install -y \
         tzdata \
         lsb-release \
         net-tools \
         git \
         curl \
         python \
         sudo \
         openjdk-8-jdk-headless

RUN update-java-alternatives -s java-1.8.0-openjdk-amd64

ADD install-build-deps.sh install-build-deps.sh
RUN yes | ./install-build-deps.sh --no-chromeos-fonts --no-prompt

ENV BUILD_URL_PREFIX=https://github.com/shiguredo/sora-webrtc-build/releases/download/
ENV BUILD_TAG=67.28.0
ENV BUILD_TOOL_VERSION=2.1.0
ENV BUILD_TOOL_FILE=sora-webrtc-build-${BUILD_TOOL_VERSION}-linux-amd64

RUN curl -L -O ${BUILD_URL_PREFIX}/${BUILD_TAG}/${BUILD_TOOL_FILE}.tar.gz \
    && tar zxf ${BUILD_TOOL_FILE}.tar.gz \
    && mv ${BUILD_TOOL_FILE} build

WORKDIR /build

RUN ./webrtc-build fetch

# RUN ./webrtc/src/build/install-build-deps.sh

RUN apt-get install -y \
         time

ADD config.json config.json

RUN ./webrtc-build build

RUN ./webrtc-build archive
