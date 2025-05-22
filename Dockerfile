FROM docker.io/golang:1.23.8 as build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./
RUN go build -o /out/myapp

FROM scratch

COPY --from=build /out/myapp /myapp
COPY migrations /migrations

ENTRYPOINT ["/myapp"]
