FROM golang:1.18.10 as build

ENV GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/cm-api-server
WORKDIR /cm-api-server
# 将代码复制到容器中
COPY . .
# 编译成二进制可执行文件app
RUN rm -rf ./bin && go build -o ./bin/cm-api-server
# 移动到用于存放生成的二进制文件的 /build 目录
WORKDIR /build
RUN cp -r /cm-api-server/bin .

FROM ubuntu:20.04
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN sed -i "s@http://.*archive.ubuntu.com@http://mirrors.tuna.tsinghua.edu.cn@g" /etc/apt/sources.list && \
    sed -i "s@http://.*security.ubuntu.com@http://mirrors.tuna.tsinghua.edu.cn@g" /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y vim net-tools tree gcc g++ p7zip-full

ENV TZ "Asia/Shanghai"
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y tzdata && \
    echo $TZ > /etc/timezone && \
    ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure tzdata -f noninteractive

WORKDIR /cm-api-server

COPY --from=build /build/ .

WORKDIR /cm-api-server/bin