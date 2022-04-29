# cles: CLI client for Elasticsearch

## What is cles?

**cles** is a tool to manage [Elasticsearch](https://www.elastic.co/elasticsearch/) easily.

You can access Elasticsearch without complicated curl options.

### Examples

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
$ cles cat indices
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
   --profile value, -p value  set profile name (default: default)
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
   --profile value, -p value  set profile name (default: default)
   --help, -h                 show help (default: false)
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