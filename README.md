# elastic-copy
Copy data from / to elasticsearch very fast

## Motivation / use case

``elasticdump`` is a very nice tool for copying a single index to another server. 

But sometimes, you need to copy more than a single index, and as fast as possible!

This why I build ``elasticcopy``.

## Multi threading

``elasticcopy`` is fast, because it use multi threads to copy data in parallel. 

The parallel unit is a ``shard``, if you want to copy 12 indices of 3 shards, this will create 36 tasks.

## Usage

```
elasticcopy --source=http://localhost:9200 --target=http://prod:9200 --indices=events-1,events-2 --threads=12
```

