# build stage
FROM golang:1.21.5-alpine3.19  as builder
RUN mkdir /app
COPY ${PWD}/ /app/
WORKDIR /app
RUN go build -o server ./cmd/server/main.go

FROM alpine
COPY --from=builder /app/server /app/
WORKDIR /app
CMD ["./server"]
