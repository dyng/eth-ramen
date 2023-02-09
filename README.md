## Ramen - A Terminal Interface for Ethereum 🍜

Ramen is a good-old terminal UI to interact with [Ethereum Network](https://ethereum.org/en/). It allows you to observe latest chain status, check account's balance and transaction history, navigate blocks and transactions, view smart contract's source code or call its functions, and many things more!

Here are some demos:

1. View accounts and transactions

    <img src="https://user-images.githubusercontent.com/1492050/216913618-9d9d708a-d350-442d-83c2-3a6c7a77614f.gif"/>

2. Call contract's function

    <img src="https://user-images.githubusercontent.com/1492050/216915763-f6c429a5-f980-4dcc-87ac-77b6cd65b62a.gif"/>

3. Sign in and transfer ethers

    <img src="https://user-images.githubusercontent.com/1492050/216917516-8cf8b02b-ba95-466a-ac6c-450f53b7f6e4.gif"/>

Additionally, **Ramen is also well designed for smart contract development**. Ramen can connect to a local chain (such as the one provided by Hardhat) to view transaction history of smart contract in development, call functions for testing, or verify its storage. Just works like Etherscan, but for your own chain!

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
