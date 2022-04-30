# cles: CLI client for Elasticsearch

[![unittest status](https://github.com/t-chov/cles/workflows/unittest/badge.svg)](https://github.com/t-chov/cles/workflows/unittest/badge.svg)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/t-chov/cles/main/LICENSE)

## What is cles?

**cles** is a tool to manage [Elasticsearch](https://www.elastic.co/elasticsearch/) easily.

You can access Elasticsearch without complicated curl options.

### Examples

An example project is available in the [`example directory`](./example/)

```sh
# The setting file is available.
# windows  : $APPDATA\cles\config.toml
# unix-like: $HOME/.config/cles/config.toml
$ cat ${HOME}/.config/cles/config.toml
[[profile]]
name = "default"
address = ["http://localhost:9200"]
username = ""
password = ""
sniff = false
$ cles -p default cat indices
health  status  index   uuid    pri     rep     docs.count      docs.deleted    store.size      pri.store.size
green   open    .geoip_databases        m0EVcoSZSAuGe7Mj5fB9tg  1       0       40      0       37.7mb  37.7mb
green   open    foo     IDKzsLa9Q2KCqN7soMoyLw  1       0       0       0       226b    226b
```

## Installation

### go get

```
go install github.com/t-chov/cles@latest
```

## Commands

```
NAME:
   cles - Command line client for Elasticsearch

USAGE:
   cles [global options] command [command options] [arguments...]

VERSION:
   0.X.Y

COMMANDS:
   indices, i, index    operate indices
   cat, c               exec cat API
   search-template, st  operate search templates
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile value, -p value  set profile name (default: "default")
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)
```


**indices** (alias: _i_ , _index_ )

```
NAME:
   cles indices - operate indices

USAGE:
   cles indices command [command options] [arguments...]

COMMANDS:
   alias, a        manage alias
   create, c, new  create index
   delete, rm      delete index
   mapping, m      get mapping
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h                 show help (default: false)
```

**cat** (alias: _c_)

```
NAME:
   cles cat - exec cat API

USAGE:
   cles cat command [command options] [arguments...]

COMMANDS:
   aliases, a  cat aliases
   indices, i  cat indices
   help, h     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

**search-template** (alias: _st_)

```
NAME:
   cles search-template - operate search templates

USAGE:
   cles search-template command [command options] [arguments...]

COMMANDS:
   list, ls        list search template
   create, c, new  create search template
   delete, rm      delete search template
   render          render search template
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h                 show help (default: false)
```

**bulk** (alias: _b_)

```
NAME:
   cles bulk - operate bulk API

USAGE:
   cles bulk command [command options] [arguments...]

COMMANDS:
   index, i  exec bulk index from ndjson
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

### Bulk index

`cles bulk index` executes [Bulk API](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html)

The file format of `--source` is different from the original format.
You can use ndjson format, but you don't insert a command line.

**wrong**

```json
{ "index" : { "_index" : "test", "_id" : "1" } } // don't add 
{ "field1" : "value1" }
{ "index" : { "_index" : "test", "_id" : "2" } } // don't add
{ "field1" : "value2" }
```

**right**

```json
{ "field1" : "value1" }
{ "field1" : "value2" }
```

The example command is below.

```sh
$ cles bulk index --source /path/to/source.ndjson <INDEX_NAME>
```


## Environment variables

**ES_ADDRESS**
This value is used as the address of Elasticsearch. You can use `,` as a separator.

**ES_USERNAME**
This value is used as the user name for authentication of Elasticsearch.

**ES_PASSWORD**
This value is used as the password for authentication of Elasticsearch.

**ES_SNIFF**
If you set this value, the client will use [sniffing](https://www.elastic.co/jp/blog/elasticsearch-sniffing-best-practices-what-when-why-how).

## FAQ

### I have no active conneciton with single-node mode.

If you have this error message with single-node mode, you have to set `ES_SNIFF` to false.

```
initClient failure! cles: no active connection found: no Elasticsearch node available
```