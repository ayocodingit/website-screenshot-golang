# build
FROM --platform=linux/amd64 registry.digitalservice.id/proxyjds/library/golang:1.17-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/services ./src/cmd

FROM --platform=linux/amd64 registry.digitalservice.id/proxyjds/library/alpine:3.16.0
COPY --from=builder /app/bin/services /
COPY --from=builder /app/src /app
EXPOSE 8080
CMD ["/services"]