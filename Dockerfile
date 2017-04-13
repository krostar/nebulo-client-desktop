FROM krostar/go-gtk:1.8-3.18

RUN apt-get update
RUN apt-get -qq install -y libcanberra-gtk3-module

RUN go get -u github.com/twitchtv/retool
RUN mkdir -p /go/src/github.com/krostar/nebulo-client-desktop
WORKDIR /go/src/github.com/krostar/nebulo-client-desktop

CMD make
