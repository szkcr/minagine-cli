FROM golang:1.19 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o minagine-cli ./src


FROM alpine:3.17

ENV TZ=Asia/Tokyo
RUN apk add --no-cache tzdata

RUN apk add --no-cache chromium chromium-chromedriver

COPY --from=builder /build/minagine-cli /app/minagine-cli

CMD /app/minagine-cli
