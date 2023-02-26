## Ramen - A Terminal Interface for Ethereum 🍜

[![Go Report Card](https://goreportcard.com/badge/github.com/dyng/ramen)](https://goreportcard.com/report/github.com/dyng/ramen)
[![Release](https://img.shields.io/github/v/release/dyng/ramen.svg)](https://github.com/derailed/k9s/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/mum4k/termdash/blob/master/LICENSE)


Ramen is a good-old terminal UI to interact with [Ethereum Network](https://ethereum.org/en/). It allows you to observe latest chain status, check account's balance and transaction history, navigate blocks and transactions, view smart contract's source code or call its functions, and many things more!

## Features

- [x] View an account's type, balance and transaction history.
- [x] View transaction details, including sender/receiver address, value, input data, gas usage and timestamp.
- [x] Decode transaction input data and display it in a human-readable format.
- [x] Call contract functions.
- [x] Import private key for transfer and calling of [non-constant](https://docs.ethers.org/v4/api-contract.html) functions.
- [ ] View contract's [ABI](https://docs.soliditylang.org/en/v0.8.13/abi-spec.html), source code, and storage.
- [x] Keep syncing with network to retrieve latest blocks and transactions.
- [ ] Show account's assets, including [ERC20](https://ethereum.org/en/developers/docs/standards/tokens/erc-20/) tokens and [ERC721](https://ethereum.org/en/developers/docs/standards/tokens/erc-721/) NFTs.
- [ ] Windows support.
- [ ] [ENS](https://ens.domains/) support.
- [ ] Navigate back and forth between pages.
- [ ] Customize key bindings and color scheme.
- [ ] Support more Ethereum JSON-RPC providers.
- [ ] Support Polygon, Binance Smart Chain, and other EVM-compatible chains.

<img src="https://user-images.githubusercontent.com/1492050/221394602-d7aaba0e-b9f8-4d73-8ddb-e81d45f289ed.gif"/>

Additionally, Ramen is also well designed for smart contract development. **Ramen can connect to a local chain (such as the one provided by [Hardhat](https://hardhat.org/))** to view transaction history of smart contract in development, call functions for testing, or verify its storage. Just works like [Etherscan](https://etherscan.io/), but for your own chain!

## Installation

### Using Package Manager

#### Homebrew

```shell
brew tap dyng/ramen && brew install ramen
```

More package managers are coming soon!

### Using Prebuilt Binaries

You can choose and download the prebuilt binary for your platform from [release page](https://github.com/dyng/ramen/releases).

### Building From Source

If you want to experience the latest features, and don't mind the risk of running an unstable version, you can build Ramen from source.

1. Clone repository

    ```shell
    git clone https://github.com/dyng/ramen.git
    ```

2. Run `go build` command

    ```shell
    go build -o ramen
    ```

## Quick Start

Ramen requires an Ethereum [JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/) provider to communicate with Ethereum network. Currently only Alchemy and local node is supported by Ramen. More providers will be added soon.

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

Then you can start Ramen by running the following command:

```shell
# connect to Mainnet
./ramen --network mainnet

# connect to Goerli Testnet
./ramen --network goerli
```

#### Key Bindings

Ramen inherits key bindings from underlying UI framework [tview](https://github.com/rivo/tview), the most frequently used keys are the following:

| Key | Action |
|---|---|
|`j`, `k`|Move cursor up and down|
|`enter`|Select an element|
|`tab`|Switch focus among elements|

#### Connect Local Network

[Hardhat](https://hardhat.org/) / [Ganache](https://trufflesuite.com/ganache/) provides a local Ethereum network for development purpose. Ramen can be used as an user interface for these local networks.

```shell
./ramen --provider local
```

## Troubleshoting

If you come across some problems when using Ramen, please check the log file `/tmp/ramen.log` to see if there are any error messages. You can also run Ramen in debug mode with command:

```shell
ramen --debug
```

If you still can't figure out the problem, feel free to open an issue on [GitHub](https://github.com/dyng/ramen/issues/new)

## Special Thanks

Ramen is built on top of many great open source projects, special thanks to [k9s](https://github.com/derailed/k9s) and [podman-tui](https://github.com/containers/podman-tui) for inspiration.

## License

Ramen is released under the Apache 2.0 license. See LICENSE for details.
