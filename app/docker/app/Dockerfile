FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /src
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/bin/app ./cmd/app

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app

ARG CONTRACT_ADDR
ARG DEPLOYER_ADDR
ARG DEPLOYER_KEY
ARG ETH_CLIENT_URL
ARG CHAIN_ID

# create config.yml consumed by the Go binary
RUN printf 'ContractAddr: "%s"\nDeployerAddr: "%s"\nDeployerKey: "%s"\nEthClientUrl: "%s"\nChainID: "%s"\n' \
    "$CONTRACT_ADDR" "$DEPLOYER_ADDR" "$DEPLOYER_KEY" "$ETH_CLIENT_URL" "$CHAIN_ID" > /app/config.yml

COPY --from=builder /app/bin/app      ./app
COPY --from=builder /src/assets       ./assets
COPY --from=builder /src/frontend     ./frontend

EXPOSE 8080
ENTRYPOINT ["./app", "-y", "/app/config.yml"]
