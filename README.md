# mount_sshfs

![travis badge](https://api.travis-ci.com/malikbenkirane/mount_sshfs.svg?branch=master&status=passed)


## Get

```
go get github.com/malikbenkirane/mount_sshfs
```

## Usage

Either, look at

```sh
go run github.com/malikbenkirane/mount_sshfs -h
```

or use

```sh
go run github.com/malikbenkirane/mount_sshfs -config remote.yml 2>/dev/null | sh
```

also make sure that the line with `user_allow_other` in `/etc/fuse.conf` is uncommented.

## Build

### Linux

```sh
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags='-w -s' -o mount_sshfs
```
