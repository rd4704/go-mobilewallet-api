FROM golang:1.11-alpine AS build-env

RUN apk --update add git

ENV GOOS=linux
ENV GOARCH=386
ENV CGO_ENABLED=0

WORKDIR /projects/src
COPY ./src .

RUN go build -o mobilewallet ./cmd/mobilewallet/main.go

FROM scratch
COPY --from=build-env /projects/src/mobilewallet .

ENTRYPOINT ["./mobilewallet"]
