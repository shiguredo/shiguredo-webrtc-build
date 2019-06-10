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
         openjdk-8-jdk-headless \
         time

RUN update-java-alternatives -s java-1.8.0-openjdk-amd64

ADD install-build-deps.sh install-build-deps.sh
RUN yes | ./install-build-deps.sh --no-chromeos-fonts --no-prompt

WORKDIR /work

ADD scripts scripts

ADD config config

RUN ./scripts/build_all_android.sh config/android-aar
