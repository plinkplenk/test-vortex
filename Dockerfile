FROM golang:1.22-alpine as build

COPY . .

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux

RUN go build -o /api -ldflags='-s -w' ./cmd/api/

FROM scratch as final
COPY --from=build /api /api
CMD ["/api"]

