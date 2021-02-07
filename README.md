# Roadrunner web dispatcher
The `revo-web-dispatch` command routes and/or redirects web queries to
port 80/443.

# Build
By default the Golang implementation of `revo-web-dispatch` builds for
GNU/Linux running on ARMv7. The requisite toolchain to compile the
source and compress the binary:

  * a `go` compiler, and
  * the `upx` packer.

Run:

```shell
git clone https://github.com/revolution-robotics/revo-web-dispatch.git
make -C revo-web-dispatch
```

To build for another architecture, e.g., GNU/Linux on AMD64, run:

```shell
GOOS=linux GOARCH=amd64 make -C revo-web-dispatch
```

For a full list of supported OSes and architectures, run:

```shell
go tool dist list
```
