FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/cep-engine ./cmd/cep-engine
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/cep-engine /cep-engine
EXPOSE 8080
ENTRYPOINT ["/cep-engine"]
