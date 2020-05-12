# hashi-bin

`hashi-bin` is a `tfenv`-like tool for all HashiCorp products.

It will likely be renamed at some point, this is only a POC.

```
brew install hashi-bin  # this does not work yet
```

```
hashi-bin install packer latest  # install to ~/.hashi-bin/packer/1.5.6
hashi-bin use packer latest      # symlink ^ to /usr/local/bin/packer
packer version                   # Packer v1.5.6
```

Full (current) help output

```
Usage: hashi-bin [--version] [--help] <command> [<args>]

Available commands are:
    download          Download to the current directory.
    install           Install to ~/.hashi-bin/{product}/{version} (or env $HASHI_BIN)
    list              List installed versions of a product.
    list-available    List available versions of a product.
    uninstall         Delete ~/.hashi-bin/{product}/{version} and remove symlink.
    use               Symlink /usr/local/bin/{product} (or env $HASHI_LINKS) -> ~/.hashi-bin/{product}/{version}
```
