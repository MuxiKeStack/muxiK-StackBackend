FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make
CMD ["/app/main", "-c", "conf/config.yaml"]
