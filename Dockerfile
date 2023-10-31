#
# NOTE: THIS DOCKERFILE IS GENERATED VIA "apply-templates.sh"
#
# PLEASE DO NOT EDIT IT DIRECTLY.
#
FROM golang:1.21.3
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"
COPY src ./app
WORKDIR ./app
RUN go mod tidy ; go build -o /go/bin/enum
RUN rm -rf *
ENTRYPOINT ["/go/bin/enum"]