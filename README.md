# mage

Mage Targets

## Example

Check the [magefiles/](magefiles/) folder for an example.

## Usage

```shell
mkdir magefiles
PKG="$(go list -m)"
( cd magefiles; go mod init "${PKG}/magefiles" )
go work init
go work use . ./magefiles
```

Then create a `magefiles/magefile.go` like the example [magefile.go](magefiles/magefile.go).
