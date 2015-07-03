FROM golang
ADD src /go/src
RUN go install github.com/jbitor/cli/jbitor-web
ENTRYPOINT /go/bin/jbitor-web
EXPOSE 8080

