FROM golang:1.23-alpine AS build

LABEL maintainer="CodeVault <codevaultllc@gmail.com>" \
      description="A production-ready Dockerfile for the Minerva Golang project"

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/minerva ./cmd/main.go

FROM scratch

LABEL maintainer="CodeVault <codevaultllc@gmail.com>" \
      description="A lightweight production image for the Minerva Golang project" \
      version="1.0"

USER 1001

COPY --from=build /app/minerva /minerva

EXPOSE 3000

CMD ["/minerva"]
