FROM golang:1.12.13 
WORKDIR $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
COPY . $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
RUN make
EXPOSE 4096 
CMD ["./main", "-c", "conf/config.yaml"]
