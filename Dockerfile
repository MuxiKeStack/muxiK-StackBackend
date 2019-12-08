FROM golang:1.12.13 
ENV GO111MODULE "on"
WORKDIR $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
COPY . $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
RUN make
EXPOSE 4096 25 465 587 
CMD ["./main", "-c", "conf/config.yaml"]
