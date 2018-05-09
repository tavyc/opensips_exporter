# opensips_exporter
A [Prometheus](http://prometheus.io/) exporter for [OpenSIPS](http://www.opensips.org/). It connects to OpenSIPS mi_json interface, fetches metrics and transforms and exposes them for consumption by Prometheus.

## OpenSIPS Configuration
Configure OpenSIPS to push serve JSON formatted stats via mi_json:
```
loadmodule "httpd.so"
loadmodule "mi_json.so"

modparam("httpd", "ip", "127.0.0.1")
modparam("httpd", "port", 8062)
```
