## Ramen - A Terminal Interface for Ethereum 🍜

Ramen is a powerful terminal UI to interact with [Ethereum Network](https://ethereum.org/en/). It allows you to observe latest chain status, check account's balance and transaction history, navigate blocks and transactions, view smart contract's source code or call its functions, and many things more!

(screenshots here)

Additionally, Ramen is also well designed for smart contract development. Ramen can connect to a local chain (such as the one provided by Hardhat) to view transaction history of smart contract in development, call functions for testing, or verify its storage. Just works like Etherscan, but for your own chain!

## Installation

Currently Ramen is under active development and its interface, key-bindings, configurations are subject to change. As a result, you can install Ramen only by building from source at this time.

### Building From Source

1. Clone repository

    ```shell
    git clone https://github.com/dyng/ramen.git
    ```

2. Run `go build` command

    ```shell
    go build -o ramen
    ```

## Quick Start

Ramen requires an Ethereum [JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/) provider to communicate with Ethereum network. Currently only Alchemy and local node is supported by Ramen. Other providers will be added soon.

In addition to the Ethereum JSON-RPC provider, Ramen also relies on the Etherscan API to access certain information that is not easily obtainable through the JSON-RPC alone, such as transaction histories and ETH prices.

To access Alchemy and Etherscan's service, you need an Api Key respectively. Please refer to their guides to obtain your own Api Key.

- [Alchemy Quickstart Guide](https://docs.alchemy.com/lang-zh/docs/alchemy-quickstart-guide)
- [Etherscan: Getting an API key](https://docs.etherscan.io/getting-started/viewing-api-usage-statistics)

When the API keys are ready, you can create a configuration file `.ramen.json` in your home directory and place keys there.

```json
{
    "apikey": "your_json_rpc_provider_api_key",
    "etherscanApikey": "your_etherscan_api_key"
}
```

Then you can launch Ramen now.

```shell
# connect to Mainnet
./ramen --network "mainnet"

# connect to Goerli Testnet
./ramen --network "goerli"
```
### Connect Local Network

[Hardhat](https://hardhat.org/) / [Ganache](https://trufflesuite.com/ganache/) provides a local Ethereum network for development purpose. Ramen can be used as an user interface for these local networks.

```shell
./ramen --provider "local"
```

## License

Ramen is released under the Apache 2.0 license. See LICENSE for details.
