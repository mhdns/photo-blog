FROM golang:alpine AS builder
ENV GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM alpine
WORKDIR /app
EXPOSE 5000
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/photo_blog .
COPY --from=builder /app/public /app/public
ENTRYPOINT [ "/app/photo_blog" ]