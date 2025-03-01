
    FROM golang:1.22.2-alpine AS build
    
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN go build -o server ./cmd/server
    
    FROM alpine:latest
    WORKDIR /app
    
    COPY --from=build /app/server /app/server
    EXPOSE 8080
    
    CMD ["/app/server"]