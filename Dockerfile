FROM golang:1.14.1-buster AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o fapr main.go

FROM gcr.io/distroless/base-debian10
EXPOSE 4000
COPY --from=build /app/fapr /
CMD ["/fapr"]