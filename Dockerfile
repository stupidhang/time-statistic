FROM golang:latest

ARG PROJECT_NAME=time-statistic
ENV GOPROXY=https://goproxy.cn
WORKDIR /usr/src/time-statistic
COPY . /usr/src/time-statistic
RUN go build -o bin/ cmd/server.go

EXPOSE 8001
ENTRYPOINT ["bin/./server"]