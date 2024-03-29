FROM golang:1.19

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential  libc-dev  pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends

ENV IMAGEMAGICK_VERSION=6.9.10-11

RUN cd && \
	wget https://github.com/ImageMagick/ImageMagick6/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib
RUN ls /usr/local/bin/ccp
RUN ls /go
RUN ls /go/src
RUN ls /go/src/gitlab.com/tokend/course-certificates/
WORKDIR /go/src/gitlab.com/tokend/course-certificates/ccp
COPY . .

#
#RUN export CGO_ENABLED=1
#RUN   export GO111MODULE=off
#RUN    export GOOS=linux
RUN  go install
RUN    go build

ENTRYPOINT  ["ccp"]