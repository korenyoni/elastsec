# ElastSec

Connects to Elasticsearch, parses heartbeat writes, creates human readable alerts with focus on what machine the event came from.

## Motivation

[ElastAlert](https://github.com/Yelp/elastalert) was too heavyweight, carrying too many alerting features. Also, ElastAlert's enhancement modules did not play well
with query_keys.

Furthermore it's to create machine-oriented event data by redoing ElastAlert's necessary work from the ground up.

## Usage

Set `ES_ADDR` to your ElasticSearch address.

## Scope

Beats:

1. [Filebeat](https://www.elastic.co/products/beats/filebeat)
2. [Auditbeat](https://www.elastic.co/products/beats/auditbeat)
