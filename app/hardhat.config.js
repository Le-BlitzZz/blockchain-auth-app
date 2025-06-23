require("@nomiclabs/hardhat-ethers");

module.exports = {
    solidity: {
        compilers: [
            {
                version: "0.8.28",
            },
        ],
        settings: {
            viaIR: true
        }
    },
    networks: {
        localhost: { url: "http://127.0.0.1:8545" }
    }
};
