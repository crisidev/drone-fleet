# Docker image for the Drone Fleet Pluging
#
#     cd $GOPATH/src/github.com/crisidev/drone-fleet
#     make deps build docker

FROM alpine:3.3

RUN apk update && \
  apk add -U ca-certificates openssh-client && \
  rm -rf /var/cache/apk/*

ADD .build/* /bin/

ENTRYPOINT ["/bin/drone-fleet"]
