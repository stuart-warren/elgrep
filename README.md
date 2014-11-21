elgrep
======

CLI tool for Elasticsearch Logstash indexes - in progress

Usage:
-----
```
$ elgrep -h
Usage of elgrep:
  -f=[@timestamp message]: fields to return (comma separated)
  -j=false: results in json
  -m=500: max number of results
  -prefix="logstash-": index name prefix (before date)
  args: query string query
```
Queries server at env ELASTICSEARCH_URL or http://localhost:9200 by default

Currently only queries the index with todays date.

See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax for details on query string syntax

