FROM golang:1.16
RUN mkdir /app
ADD . /app/
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY="https://goproxy.cn,direct"
WORKDIR /app/
RUN go mod tidy
RUN make
EXPOSE 4096 25 465 587
CMD ["./main"]

#ENV GOPROXY "https://goproxy.cn,direct"
#ENV GO111MODULE "on"
#WORKDIR $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
#COPY . $GOPATH/src/github.com/MuxiKeStack/muxiK-StackBackend
#RUN make
#EXPOSE 4096 25 465 587
#CMD ["./main", "-c", "conf/config.yaml"]