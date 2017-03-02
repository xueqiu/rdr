RDR: redis data reveal
=================================================

RDR(redis data reveal) is a tool to parse redis rdbfile. Comparing to [redis-rdb-tools](https://github.com/sripathikrishnan/redis-rdb-tools) ,rdr implemented by golang, much faster (especially for big rdbfile).

## Usage

```
NAME:
   rdr - a tool to parse redis rdbfile

USAGE:
   rdr [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

COMMANDS:
     show     show statistical information of rdbfile by webpage
     keys     get all keys from rdbfile
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
[Download](https://baidu.com) here.

## Exapmle
```
$ ./rdr show -p 8080 example1.rdb example2.rdb
```

```
$ ./rdr keys example.rdb
```

## 

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.