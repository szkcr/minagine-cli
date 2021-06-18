FROM golang:1.16 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o automate-minagine


FROM alpine:3.14.0

ENV TZ=Asia/Tokyo
RUN apk add --no-cache tzdata

RUN apk add --no-cache chromium chromium-chromedriver

COPY --from=builder /build/automate-minagine /app/automate-minagine

CMD /app/automate-minagine
