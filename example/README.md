# Example index

It uses [Motor Vehicle Collisions by NYC](https://data.cityofnewyork.us/Public-Safety/Motor-Vehicle-Collisions-Crashes/h9gi-nx95)

## How to use

**Create new index**
```sh
$ cles indices create -b /path/to/crashes.mapping.json nyc-crashes
health  status  index   uuid    pri     rep     docs.count      docs.deleted    store.size      pri.store.size
green   open    .geoip_databases        tbTGZ7mWTNyr1hfw0oMUeQ  1       0       40      0       37.7mb  37.7mb
green   open    nyc-crashes     6jQtdiq1Qfaqdz4QaaglXA  1       0       0       0       226b    226b
# or use stdin
$ cat /path/to/crashes.mapping.json | cles indices create nyc-crashes
```

**Create new alias**
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