# WARNING: it is a very experimental Lachesis

Consensus platform for distributed applications. Don't use for product!

## GitHub

[![Commit Activity](https://img.shields.io/github/commit-activity/w/Fantom-foundation/lachesis-ex.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/lachesis-ex/commits/master)
[![LoC](https://tokei.rs/b1/github/Fantom-foundation/lachesis-ex?category=lines)](https://github.com/Fantom-foundation/lachesis-ex)


[github-issues-image]: https://img.shields.io/github/issues/Fantom-foundation/lachesis-ex.svg?style=flat-square&logo=github
[github-issues-url]: https://github.com/Fantom-foundation/lachesis-ex/issues

## Building the source

Building `lachesis` requires both a Go (version 1.13 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
go build -o ./build/lachesis ./cmd/lachesis
```
The build output is ```build/lachesis``` executable.

Do not clone the project into $GOPATH, due to the Go Modules. Instead, use any other location.

## Running `lachesis`

Going through all the possible command line flags is out of scope here,
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `lachesis` instance.

### Mainnet

Launching `lachesis` for mainnet with default settings:

```shell
$ lachesis
```

### Configuration

As an alternative to passing the numerous flags to the `lachesis` binary, you can also pass a
configuration file via:

```shell
$ lachesis --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to
export your existing configuration:

```shell
$ lachesis --your-favourite-flags dumpconfig
```

#### Docker quick start

One of the quickest ways to get Lachesis up and running on your machine is by using
Docker:

```shell
cd docker/
make
docker run -d --name lachesis-node -v /home/alice/lachesis:/root \
           -p 5050:5050 \
          "lachesis" \
          --port 5050 \
          --nat extip:YOUR_IP
```

This will start `lachesis` with ```--port 5050 --nat extip:YOUR_IP``` arguments, with DB files inside ```/home/alice/lachesis/.lachesis```

Do not forget `--rpcaddr 0.0.0.0`, if you plan to access RPC from other containers
and/or hosts. By default, `lachesis` binds to the local interface and RPC endpoints is not
accessible from the outside.

To find out your enode ID, use:
```shell
docker exec -i lachesis-node /lachesis --exec "admin.nodeInfo.enode" attach
```
To get the logs:
```
docker logs lachesis-node
```

#### Validator

To launch a validator, you have to use `--validator` flag to enable events emitter. Also you have to either use `--unlock` / `--password` flags or unlock
validator account manually. Validator account should be unlocked for signing events.

```shell
$ lachesis --nousb --validator 0xADDRESS --unlock 0xADDRESS --password /path/to/password
```

#### Participation in discovery

Optionally you can specify your public IP to straighten connectivity of the network.
Ensure your TCP/UDP p2p port (5050 by default) isn't blocked by your firewall.

```shell
$ lachesis --nat extip:1.2.3.4
```

## Dev

### Running testnet

To run a testnet node, you have to add `--testnet` flag every time you use `lachesis`:

```shell
$ lachesis --testnet # launch node
$ lachesis --testnet attach # attach to IPC
$ lachesis --testnet account new # create new account
```

### Testing

Lachesis has extensive unit-testing. Use the Go tool to run tests:
```shell
go test ./...
```

If everything goes well, it should output something along these lines:
```
?       github.com/Fantom-foundation/lachesis-ex/event_check/basic_check    [no test files]
?       github.com/Fantom-foundation/lachesis-ex/event_check/epoch_check    [no test files]
?       github.com/Fantom-foundation/lachesis-ex/event_check/heavy_check    [no test files]
?       github.com/Fantom-foundation/lachesis-ex/event_check/parents_check  [no test files]
ok      github.com/Fantom-foundation/lachesis-ex/evm_core   (cached)
ok      github.com/Fantom-foundation/lachesis-ex/gossip (cached)
?       github.com/Fantom-foundation/lachesis-ex/gossip/fetcher [no test files]
?       github.com/Fantom-foundation/lachesis-ex/gossip/occuredtxs [no test files]
ok      github.com/Fantom-foundation/lachesis-ex/gossip/ordering    (cached)
ok      github.com/Fantom-foundation/lachesis-ex/gossip/packsdownloader    (cached)
```

### Operating a private network (fakenet)

Fakenet is a private network optimized for your private testing.
It'll generate a genesis containing N validators with equal stakes.
To launch a validator in this network, all you need to do is specify a validator ID you're willing to launch.

Pay attention that validator's private keys are deterministically generated in this network, so you must use it only for private testing.

Maintaining your own private network is more involved as a lot of configurations taken for
granted in the official networks need to be manually set up.

To run the fakenet with just one validator (which will work practically as a PoA blockchain), use:
```shell
$ lachesis --fakenet 1/1
```

To run the fakenet with 5 validators, run the command for each validator:
```shell
$ lachesis --fakenet 1/5 # first node, use 2/5 for second node
```

If you have to launch a non-validator node in fakenet, use 0 as ID:
```shell
$ lachesis --fakenet 0/5
```

After that, you have to connect your nodes. Either connect them statically or specify a bootnode:
```shell
$ lachesis --fakenet 1/5 --bootnodes "enode://verylonghex@1.2.3.4:5050"
```

### Running the demo

For the testing purposes, the full demo may be launched using:
```shell
cd docker/
make # build docker image
./start.sh # start the containers
./stop.sh # stop the demo
```

The full demo doesn't spin up very fast. To avoid the full docker image building, you may run the integration test instead:
```shell
go test -v ./integration/...
```
Adjust test duration, number of nodes and logs verbosity in the test source code.
