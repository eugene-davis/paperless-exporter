# Paperless-Exporter

A simple Prometheus exporter for [Paperless-ngx](https://docs.paperless-ngx.com/) statistics. At present it supports:

* the total number of documents
* the total number of characters
* the number of documents in the inbox
* the number of documents for each mime type
* the number of file tasks

## Configuring access

You'll need to [create a token](https://docs.paperless-ngx.com/api/#authorization) for a user which has access view on the "PaperlessTask" permission.

## Configuring the Exporter

The exporter is entirely configured with environmental variables.

| Variable | Description | Default |
| -------- | ----------- | ------- |
| VERBOSITY | Sets the logging verbosity can be one of: DEBUG, INFO, WARN, ERROR | INFO |
| URL | URL for the Paperless-ngx instance. | http://localhost:8000 |
| REFRESH_SECS | Refresh frequency for the metrics. | 1800 |
| HOST_HEADER | Allows setting custom headers if needed. | "" |
| PAPERLESS_TOKEN | Paperless token. Either this or PAPERLESS_TOKEN_FILE must be set | "" |
| PAPERLESS_TOKEN_FILE | Path to paperless token file, which is a file with only the paperless token in it. Either this or PAPERLESS_TOKEN must be set. | "" |
| METRICS_PORT | Port for the exporter to listen on. | 8001 |
