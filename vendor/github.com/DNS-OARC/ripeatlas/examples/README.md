# Examples using ripeatlas Go bindings

To use any of these examples, do the following from the top directory of
`ripeatlas`.

```shell
mkdir -p "$GOPATH/src/github.com/DNS-OARC"
ln -s "$PWD" "$GOPATH/src/github.com/DNS-OARC/"
make dep
make
cd examples/
make
```

After building all examples programs, each one can be run with `-help` to see
what options are available.

## measurements

This will fetch and print metadata about measurements.

## probes

This will fetch and print metadata about probes.

## reader

This will read and display results from measurements.

## streamer

This will open a stream to a specific type of measurement or all and display
results as they come in.
