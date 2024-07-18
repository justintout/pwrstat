# pwrstat-serve

> Restful JSON API serving `pwrstat` status

## Installation

### ...using `go install`

```
$ go install github.com/justintout/pwrstat/cmd/pwrstat-serve@latest
$ sudo setcap cap_setuid=ep ~/.go/bin/pwrstat-serve
```

### ...using a release tarball

> Not yet supported.

### ...by building manually

```
$ git clone github.com/justintout/pwrstat.git
$ go -C pwrstat/cmd/pwrstat-serve/ build -o ../../pwrstat-serve
$ sudo setcap cap_setuid=ep pwrstat/pwrstat-serve
```

## Running

```
$ pwrstat-serve
```

`pwrstat`'s default installation requires root to run.
This binary, by default, uses `setuid` to elevate to root to execute `pwrstat` and will need capabilities to do so.
Thus, you must run `sudo setcap cap_setuid=ep pwrstar-serve` to allow the binary to setuid.
If you have a user who can run `pwrstat` without root privileges, you can pass the `-noroot` flag and the binary won't use setuid.

`pwrstat`'s default installation location is `/usr/sbin/pwrstat`.
If you have the binary installed in a nonstandard location, you can use the `-path` flag to provide the alternate path to `pwrstat`.

## Usage

```
$ ./pwrstat-serve -h
Usage of ./pwrstat-serve:
  -host string
        host for server to listen on, default: 0.0.0.0 (default "0.0.0.0")
  -noroot
        execute pwrstat without elevating to root, default: false
  -path string
        path to the pwrstat executable, default: /usr/sbin/pwrstat (default "/usr/sbin/pwrstat")
  -port int
        port for server to listen on, default: 7977 (default 7977)
```
