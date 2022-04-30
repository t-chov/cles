# Example index

It uses [Motor Vehicle Collisions by NYC](https://data.cityofnewyork.us/Public-Safety/Motor-Vehicle-Collisions-Crashes/h9gi-nx95)

## Requirement

You have to make connection to Elasticsearch.
I recommend [single-node cluster](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-dev-mode) mode.

## How to use

**Create an index**
```sh
$ cles indices create -b /path/to/crashes.mapping.json nyc-crashes
health  status  index   uuid    pri     rep     docs.count      docs.deleted    store.size      pri.store.size
green   open    .geoip_databases        tbTGZ7mWTNyr1hfw0oMUeQ  1       0       40      0       37.7mb  37.7mb
green   open    nyc-crashes     6jQtdiq1Qfaqdz4QaaglXA  1       0       0       0       226b    226b
# or use stdin
$ cles indices create nyc-crashes < /path/to/crashes.mapping.json
```

**Create an alias**
```sh
$ cles indices alias nyc-crashes crashes
alias   index   routing.index   routing.search  is_write_index
crashes nyc-crashes     -       -       -
```

**Download crash data from NYC Open data**
Link: https://data.cityofnewyork.us/Public-Safety/Motor-Vehicle-Collisions-Crashes/h9gi-nx95

**Indexing**

The convert script is written in Python.

```sh
$ python3 converter.py /path/to/crashes.csv > crashes.ndjson
$ cles bulk index --source crashes.ndjson -i collision_id nyc-crashes
# or use stdin
$ python3 converter.py /path/to/crashes.csv | cles bulk index -i collision_id nyc-crashes
```

**Create a search template**
```sh
$ cles search-template -b feature/add_example create collision-point-search
# or use stdin
$ cles search-template create collision-point-search < /path/to/collision-point-search.template.json
```

**Search with a template**
```sh
# search around Yankee Stadium
$ cles search-template search -i nyc-crashes --params='{"latitude": "40.8296", "longitude":"-73.9262", "distance": "500m"}' collision-point-search
# or use stdin
$ cles search-template search -i nyc-crashes collision-point-search < /path/to/stadium.params.json
```