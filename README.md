## Ramen - A Terminal Interface for Ethereum 🍜

Ramen is a good-old terminal UI to interact with [Ethereum Network](https://ethereum.org/en/). It allows you to observe latest chain status, check account's balance and transaction history, navigate blocks and transactions, view smart contract's source code or call its functions, and many things more!

1. View account balance, transactions etc.

    <img src="https://user-images.githubusercontent.com/1492050/215658150-93b09da7-52b2-4366-ba24-4a56668cf2a8.png"/>

2. Call contract's function

    <img src="https://user-images.githubusercontent.com/1492050/215658167-b38bcf0b-8dd4-4b95-8198-5411fe3fd7e0.png"/>

Additionally, Ramen is also well designed for smart contract development. Ramen can connect to a local chain (such as the one provided by Hardhat) to view transaction history of smart contract in development, call functions for testing, or verify its storage. Just works like Etherscan, but for your own chain!

## Installation

**Currently Ramen is under active development and its interface, key-bindings, configurations are subject to change.**

As a result, you can install Ramen only by building from source at this moment.

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

When the API keys are ready, you can create a configuration file `.ramen.json` in your home directory (e.g. `~/.ramen.json`) and place keys there.

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

### Key Bindings

Ramen inherits key bindings from underlying UI framework [tview](https://github.com/rivo/tview), the most frequently used keys are the following:

| Key | Action |
|---|---|
|`j`, `k`|Move cursor up and down|
|`enter`|Select an element|
|`tab`|Switch focus among elements|

### Connect Local Network

[Hardhat](https://hardhat.org/) / [Ganache](https://trufflesuite.com/ganache/) provides a local Ethereum network for development purpose. Ramen can be used as an user interface for these local networks.

```shell
./ramen --provider "local"
```

## Contribution

Ramen is an open source project, any kind of contribution is welcome! Feel free to open issues for feature request, bug report or discussions.

## License

Ramen is released under the Apache 2.0 license. See LICENSE for details.
