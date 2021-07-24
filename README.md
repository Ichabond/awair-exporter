# awair-exporter
[![Go Report Card](https://goreportcard.com/badge/github.com/Ichabond/awair-exporter?style=flat-square)](https://goreportcard.com/report/github.com/Ichabond/awair-exporter)
## Overview

Awaire Exporter is a basic Prometheus exporter for the [Awair Local API](https://support.getawair.com/hc/en-us/articles/360049221014-Awair-Element-Local-API-Feature). All data exposed through the API is exported as a metric

## Use
`awair-exporter $ENDPOINT`

A sample systemd unit file is also provided in [awair-exporter.service](awair-exporter.service)