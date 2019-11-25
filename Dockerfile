FROM golang:latest 
WORKDIR $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
COPY . $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
RUN make
CMD ["./main", "-c", "conf/config.yaml"]
