require("@nomiclabs/hardhat-ethers");
require("dotenv").config();

const { SEPOLIA_RPC_URL, PRIVATE_KEY } = process.env;

module.exports = {
  solidity: {
    compilers: [{ version: "0.8.28" }],
    settings: { viaIR: true },
  },
  networks: {
    localhost: { url: "http://127.0.0.1:8545" },

    sepolia: {
      url: SEPOLIA_RPC_URL || "",
      chainId: 11155111,
      accounts: PRIVATE_KEY ? [PRIVATE_KEY] : [],
    },
  },
};
