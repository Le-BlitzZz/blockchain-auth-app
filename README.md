# Blockchain‑Auth App

A reference implementation of **token‑based authentication on Ethereum** powered by a Golang backend, a Redis cache, and a non‑transferable ERC‑721 (“VIP Pass”).  
Users authenticate by signing a nonce in MetaMask; the backend checks ownership of the NFT on‑chain.

---

## ✨ Quick Start (for Users)

1. **Download the Docker stack**

   ```bash
   mkdir blockchain-auth && cd blockchain-auth
   curl -O https://raw.githubusercontent.com/Le-BlitzZz/blockchain-auth-app/main/app/setup/compose.yaml
   ```

2. **Start the backend**

   ```bash
   docker compose up -d          # backend at http://localhost:8080
   ```

   The compose file pulls a pre‑built image already configured for the Sepolia testnet.

3. **Grab the CLI**

   * Visit the [v1.1.0 release](https://github.com/Le-BlitzZz/blockchain-auth-app/releases/tag/v1.1.0).
   * Download the tarball for your OS / CPU.
   * Extract it; the binary is called `blockchain-auth-client`.

4. **Run an action**

   ```bash
   ./blockchain-auth-client vip        # or mint / burn
   ```

   The CLI opens your browser, MetaMask asks for a signature, and the result is printed in your terminal.

---

## ⚙️ Developer Guide

### 1 · Clone and launch the dev environment

```bash
git clone https://github.com/Le-BlitzZz/blockchain-auth-app.git
cd blockchain-auth-app
docker compose up -d            # app container, Redis, local Hardhat
make terminal                   # shell inside the Go container (backend)
```

### 2 · Build and run the **backend** on the local Hardhat chain

```bash
make build                      # compiles backend (app/bin/app)
make deps                       # inside Hardhat container: npm install + local contract deploy
make start                      # runs backend with configs/local.yml
```

### 3 · Build and test the **CLI client**

Open a second terminal tab:

```bash
cd blockchain-auth-app/client
make build                      # produces client/bin/client
./bin/client vip                # or mint / burn
```

The CLI automatically opens the MetaMask flow and streams status via SSE.

### 4 · Point backend + client to Sepolia

1. **Create secrets**

   * `app/configs/testnet.yml` with your Sepolia contract address, deployer address, private key, RPC URL, and chain ID 11155111.
   * `.env` in `app/` based on `.env.example`:

     ```dotenv
     SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/<yourKey>
     PRIVATE_KEY=0x<privateKey>
     ```

2. **Deploy and start**

   ```bash
   make deps-testnet             # npm install + deploy to Sepolia
   make start-testnet            # ./bin/app -y configs/testnet.yml
   ```

   The CLI doesn’t need changes; it always points to `localhost:8080`.

### 5 · Updating the smart contract

If you modify `contracts/VIPPass.sol`:

```bash
make compile                    # re‑compiles Solidity
make contract                   # regenerates Go bindings
```

Artifacts in `app/artifacts/` and `internal/contract/` are committed so other contributors can build without Solidity installed.

### 6 · Handy Make targets

Backend (`app/`):

* `make build` – compile Go binary  
* `make deps` – install NPM deps & deploy to local chain  
* `make deps-testnet` – deploy to Sepolia  
* `make start` / `make start-testnet` – run backend  

Client (`client/`):

* `make build` – compile CLI  
* `make vip | mint | burn` – run pre‑set actions

Global:

* `docker compose up -d` / `docker compose down` – start / stop full stack  
* `make terminal` / `make terminal-hardhat` – shell into containers  

---

## License

MIT — © Le‑BlitzZz 2025
