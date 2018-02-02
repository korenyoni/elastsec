# ElastSec

Connects to Elasticsearch, parses heartbeat writes, creates human readable alerts, focusing on what machine the event came from.

## Motivation

[ElastAlert](https://github.com/Yelp/elastalert) was too heavyweight, carrying too many alerting features. Also, ElastAlert's enhancement modules did not play well
with query_keys.

Furthermore it's more feasible to create machine-oriented event data by redoing ElastAlert's necessary work from the ground up.

## Usage

1. Set `ES_ADDR` to your ElasticSearch address.
2. Add `-w /etc/ -p wa` to your auditbeat.yml
3. Use the following auditbeat configuration:
```
- module: audit
  metricsets: [kernel]
  kernel.audit_rules: |

    # Identity changes.
    -w /etc/group -p wa -k identity
    -w /etc/passwd -p wa -k identity
    -w /etc/gshadow -p wa -k identity
    -w /etc/ -p wa

    # Unauthorized access attempts.
    -a always,exit -F arch=b64 -S open,creat,truncate,ftruncate,openat,open_by_handle_at -F exit=-EACCES -k access
    -a always,exit -F arch=b64 -S open,creat,truncate,ftruncate,openat,open_by_handle_at -F exit=-EPERM -k access

- module: audit
  metricsets: [file]
  file.paths:
  - /bin
  - /usr/bin
  - /sbin
  - /usr/sbin
  - /etc

```
4. In filebeat.yml, under `filebeat.prospectors`, add: `scan_frequency: 1s`
5. `make && ./elastsec`

## Scope

Beats:

1. [Filebeat](https://www.elastic.co/products/beats/filebeat)
2. [Auditbeat](https://www.elastic.co/products/beats/auditbeat)
