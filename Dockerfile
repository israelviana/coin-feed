FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -o /bin/app ./cmd/

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=build /bin/app /app/app

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/app"]