FROM golang:1.24.4 AS builder

WORKDIR /app
COPY . /app

RUN go install

# Statically compile our app for use in a distroless container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static

RUN adduser --system --group runner-user # Crear un usuario

COPY --from=builder /app/app /app

USER runner-user
ENTRYPOINT ["/app"]
