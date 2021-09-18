FROM node:15.11.0 as builder
ARG VUE=/usr/src/vue
COPY ./dashboard $VUE
WORKDIR $VUE
RUN yarn config set registry https://registry.npm.taobao.org
RUN yarn install && npm run build:prod

FROM golang:alpine AS development
WORKDIR $GOPATH/src
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
ENV CGO_ENABLED=1
COPY . .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add libpcap-dev
RUN apk add build-base
RUN go mod tidy & go mod vendor
RUN go build -a -ldflags '-extldflags "-static"' -o ./bin/sqlaudit  ./cmd/sqlaudit 


FROM alpine:latest AS production
WORKDIR /opt/sqlaudit
COPY --from=development /go/src/bin/sqlaudit .
RUN mkdir etc
COPY --from=development /go/src/etc/config.yaml etc/
RUN mkdir -p dashboard/dist
COPY --from=builder /usr/src/vue/dist dashboard/dist

EXPOSE 6969
EXPOSE 9797
EXPOSE 9898
ENTRYPOINT ["./sqlaudit", "-config", "etc/config.yaml"]

