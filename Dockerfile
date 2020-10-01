# Stage 1: Build backend
FROM golang:1.15.2-buster AS backend

WORKDIR /go/src/app

RUN apt update
COPY go.mod ./
# COPY go.sum ./
RUN go mod download
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app

# Stage 2: Execute
FROM alpine
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=backend /go/bin/app .
COPY templates ./templates
COPY assets ./assets
COPY docs ./docs
EXPOSE 2983

ENTRYPOINT ["./app"]
