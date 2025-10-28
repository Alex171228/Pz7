# ---------- Build stage ----------
FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---------- Runtime stage (distroless) ----------
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /app/server /app/server
EXPOSE 8080
USER nonroot:nonroot
ENV PORT=8080
CMD ["/app/server"]
