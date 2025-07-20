up:
	docker compose up -d

down:
	docker compose down

cleanup:
	docker compose down -v

terminal:
	docker compose exec -u root app bash

terminal-hardhat:
	docker compose exec -u root hardhat bash

terminal-redis:
	docker compose exec -u root redis redis-cli

build:
	go build -o bin/app ./cmd/app/main.go

start:
	./bin/app

deps:
	npm install
	make deploy

deploy:
	npx hardhat run --network localhost scripts/deploy.js

compile:
	npx hardhat compile

contract:
	jq '.abi' artifacts/contracts/VIPpass.sol/VIPPass.json > VIPPass.abi
	jq -r '.bytecode' artifacts/contracts/VIPpass.sol/VIPPass.json > VIPPass.bin
	abigen \
	--abi VIPPass.abi \
	--bin VIPPass.bin \
	--pkg contract \
	--out internal/contract/vippass.go
	rm -rf VIPPass.abi VIPPass.bin
