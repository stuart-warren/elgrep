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
Queries server at env ELASTICSEARCH_HOST or http://localhost:9200 by default

Currently only queries the index with todays date.

See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax for details on query string syntax

Examples:
--------
```
$ export ELASTICSEARCH_URL=http://myes.box:9200
# Get latest logs *
$ elgrep 
# Get latest logs where myfield has the value stuff
$ elgrep myfield:stuff
# Get logs where myfield has a value starting with st
$ elgrep myfield:st*
# Boolean operators (OR if not specified)
$ elgrep (myfield:stuff AND thatfield:other)
# Regex (see docs above)
$ elgrep myfield:/st.+/
# Show more fields
$ elgrep -f @timestamp,message,myfield,thatfield
```
