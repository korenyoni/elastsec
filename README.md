# ElastSec

Connects to Elasticsearch, parses heartbeat writes, and creates the following host-oriented alerts:

1. New Priviledge Escalation
2. Failed file change attempt
3. Failed file access attempt
4. File permissions change
5. New SSH connection
6. Failed SSH connection attempt (password, invalid user)

## Motivation

[ElastAlert](https://github.com/Yelp/elastalert) was too heavyweight, carrying too many alerting features. Also, ElastAlert's enhancement modules did not play well
with query_keys.

Furthermore it's more feasible to create machine-oriented event data by redoing ElastAlert's necessary work from the ground up.

## Usage

1. Set `ES_ADDR` to your ElasticSearch address, `ESEC_SLACK_WEBHOOK` to your slack webhook, and `STMP_SEND_ADDR` to the email you would like to notify.
2. `ESEC_AGG_DURATION` and `ESEC_EMAIL_DURATION` can be optionally set (e.g. `2hr`,`24h`). It is recommended to add a couple extra more seconds for email as it will capture the aggregation events.
3. Add `-w /etc/ -p wa` to your auditbeat.yml
4. Use the following auditbeat configuration:
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
5. In filebeat.yml, under `filebeat.prospectors`, add: `scan_frequency: 1s`
6. `make && ./elastsec`

## Requirements

1. [Elasticsearch](https://www.elastic.co/products/elasticsearch)
2. [Filebeat](https://www.elastic.co/products/beats/filebeat)
3. [Auditbeat](https://www.elastic.co/products/beats/auditbeat)
4. `sendmail` configured via `ssmtp` (including revaliases) or another SMTP utility.

You will need a version of Go relatively recent to `1.9.3` to build the binary yourself. A glide configuration and lock-file is included.
