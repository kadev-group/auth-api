FROM golang:1.21.9-alpine as builder

ENV GITHUB_ACCESS_TOKEN=0

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache git
RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./main.go

FROM alpine:3.18 as production

COPY --from=builder /build/auth ./auth

CMD ["/auth"]
