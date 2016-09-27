# mgomgo [![Build Status](https://travis-ci.org/furushchev/mgomgo.svg)](https://travis-ci.org/furushchev/mgomgo)

database migration tool for mongodb

### Installation

- Binary Install

  Download binary from [Release Page](https://github.com/furushchev/mgomgo/releases), then install.

  ```bash
# for linux 64bit
sudo dpkg -i mgomgo_X.X.X_amd64.deb
```

- Install from source

  1. Install `go`, setup `GOPATH`
  2. Build from source

    ```bash
cd $GOPATH
go get github.com/furushchev/mgomgo
```

### Usage

```bash
mgomgo <src> <dst>
```

e.g.:

  ```bash
mgomgo -c 10 -t 60 mongodb://user:pass@abc.com:27017/db_1/col_1 mongodb://user2:pass@def.co.jp:27018/db_2/col_2
```

  - `-c`: Concurrent session to insert data to destination
  - `-t`: Timeout second to connect database server

### Author

Yuki Furuta <<furushchev@jsk.imi.i.u-tokyo.ac.jp>>
