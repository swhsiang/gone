# Build Stage
FROM wen777/gone:gobuildimage AS build-stage

LABEL app="build-gone"
LABEL REPO="https://github.com/swhsiang/gone"

ENV GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PROJPATH=/gopath/src/github.com/swhsiang/gone

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /gopath/src/github.com/swhsiang/gone
WORKDIR /gopath/src/github.com/swhsiang/gone

RUN make build-alpine

# Final Stage
FROM wen777/gone:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/swhsiang/gone"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/gone/bin

WORKDIR /opt/gone/bin

COPY --from=build-stage /gopath/src/github.com/swhsiang/gone/bin/gone /opt/gone/bin/
RUN chmod +x /opt/gone/bin/gone

CMD /opt/gone/bin/gone