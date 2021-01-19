FROM alpine:3

RUN apk --no-cache add curl
ADD bin/nexthop-static /nexthop

ENTRYPOINT [ "/nexthop" ]
