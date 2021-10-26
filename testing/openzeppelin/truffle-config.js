module.exports = {
  test_directory: "./openzeppelin-contracts/test",
  migrations: "./openzeppelin-contracts/migrations",
  contracts_directory: "./openzeppelin-contracts/contracts",
  contracts_build_directory: "./openzeppelin-contracts/build/output",
  mocha: {
    reporter: "mocha-spec-json-output-reporter",
    reporterOptions: {
      fileName: "output.json",
    },
  },
  networks: {
    development: {
      host: "127.0.0.1",
      port: 23889, //Switch to 23888 for local HTTP Server, look at Makefile run-janus
      network_id: "*",
      gas: 3000000,
      gasPrice: "0x64",
    },
    testing: {
      host: "127.0.0.1",
      port: 23888,
      network_id: "*",
      gas: 3000000,
      gasPrice: "0x64",
    },
    docker: {
      host: "janus",
      port: 23889,
      network_id: "*",
      gas: 3000000,
      gasPrice: "0x64",
    },
    ganache: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
    },
    testnet: {
      host: "hk1.s.qtum.org",
      port: 23889,
      network_id: "*",
      from: "0x7926223070547d2d15b2ef5e7383e541c338ffe9",
      gasPrice: "0x64",
    },
  },
  compilers: {
    solc: {
      version: "^0.8.0",
      settings: {
        optimizer: {
          enabled: true,
          runs: 1,
        },
      },
    },
  },
};
