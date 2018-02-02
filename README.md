# ElastSec

Connects to Elasticsearch, parses heartbeat writes, creates human readable alerts, focusing on what machine the event came from.

## Motivation

[ElastAlert](https://github.com/Yelp/elastalert) was too heavyweight, carrying too many alerting features. Also, ElastAlert's enhancement modules did not play well
with query_keys.

Furthermore it's more feasible to create machine-oriented event data by redoing ElastAlert's necessary work from the ground up.

## Usage

1. Set `ES_ADDR` to your ElasticSearch address.
2. Add `-w /etc/ -p wa` to your auditbeat.yml
3. `./elastsec`

## Scope

Beats:

1. [Filebeat](https://www.elastic.co/products/beats/filebeat)
2. [Auditbeat](https://www.elastic.co/products/beats/auditbeat)
