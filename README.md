# mgomgo [![Build Status](https://travis-ci.org/furushchev/mgomgo.svg)](https://travis-ci.org/furushchev/mgomgo)

database migration tool for mongodb

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
