# mount_sshfs

![travis build status](https://api.travis-ci.com/malikbenkirane/mount_sshfs.svg?branch=master)
![goreportcard](https://goreportcard.com/badge/github.com/malikbenkirane/mount_sshfs)


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
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" github.com/malikbenkirane/mount_sshfs/cmd/mount_sshfs
```
