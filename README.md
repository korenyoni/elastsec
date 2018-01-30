# ElastSec

Connects to Elasticsearch, parses heartbeat writes, creates human readable alerts.

## Motivation

[ElastAlert](https://github.com/Yelp/elastalert) was too heavyweight, carrying too many alerting features. Also, ElastAlert's enhancement modules did not play well
with query_keys.

Furthermore it's to create machine-oriented event data by redoing ElastAlert's necessary work from the ground up.

## Usage

Set `ES_ADDR` to your ElasticSearch address.
