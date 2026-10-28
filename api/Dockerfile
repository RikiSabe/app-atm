FROM golang:1.25-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/inventory-api .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/inventory-api .

EXPOSE 5665
CMD ["/app/inventory-api"]
