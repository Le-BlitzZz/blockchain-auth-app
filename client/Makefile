.PHONY: build

build:
	go build -o bin/client ./cmd/client

start:
	./bin/client

hello:
	./bin/client hello

mint:
	./bin/client mint

vip:
	./bin/client vip

burn:
	./bin/client burn

private_ca:
	openssl genrsa -out certs/ca.key 4096
	openssl req -x509 -new -nodes -key certs/ca.key \
	-subj "/CN=LeBlitzZz-Root-CA" \
	-days 7300 -out certs/ca.pem

client_key_csr:
	openssl genrsa -out certs/client.key 2048
	openssl req -new -key certs/client.key \
	-subj "/CN=blockchain-auth-cli" \
	-out certs/client.csr

sign_csr_with_ca: private_ca client_key_csr
	openssl x509 -req \
	-in certs/client.csr \
	-CA certs/ca.pem -CAkey certs/ca.key \
  	-CAcreateserial -days 3650 \
	-out certs/client.crt

clean:
	rm -f certs/client.csr certs/ca.srl
