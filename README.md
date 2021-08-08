# LVM metrics for Prometheus

[![Latest release](https://img.shields.io/github/v/release/hansmi/prometheus-lvm-exporter)][releases]
[![Release workflow](https://github.com/hansmi/prometheus-lvm-exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/hansmi/prometheus-lvm-exporter/actions/workflows/release.yaml)
[![CI workflow](https://github.com/hansmi/prometheus-lvm-exporter/actions/workflows/ci.yaml/badge.svg)](https://github.com/hansmi/prometheus-lvm-exporter/actions/workflows/ci.yaml)
[![Go reference](https://pkg.go.dev/badge/github.com/hansmi/prometheus-lvm-exporter.svg)](https://pkg.go.dev/github.com/hansmi/prometheus-lvm-exporter)

Prometheus exporter for the [Logical Volume Manager][lvm2] (LVM,
[Wikipedia][wikipedia]). It is only compatible with Linux and has been tested
with LVM 2.03. All fields related to physical volumes, volume groups and
logical volumes are reported, either as a standalone metric for numeric values
or as a label on a per-entity info metric.

## Usage

`prometheus-lvm-exporter` listens on TCP port 8081 by default. To listen on
another address use the `-web.listen-address` flag (e.g.
`-web.listen-address=127.0.0.1:3000`).

TLS and HTTP basic authentication is supported through the [Prometheus exporter
toolkit][toolkit]. A configuration file can be passed to the `-web.config` flag
([documentation][toolkitconfig]).

See the `--help` output for more flags.

## Installation

Pre-built binaries are provided for all [releases][releases]:

* Binary archives (`.tar.gz`)
* Debian/Ubuntu (`.deb`)
* RHEL/Fedora (`.rpm`)

With the source being available it's also possible to produce custom builds
directly using [Go][golang] or [GoReleaser][goreleaser].

## Example metrics

```
lvm_pv_info{pv_fmt="lvm2",pv_name="/dev/sda1",pv_uuid="yc1zVe-…"} 1
lvm_pv_info{pv_fmt="lvm2",pv_name="/dev/sdb1",pv_uuid="WVIH97-…"} 1

lvm_pv_free_bytes{pv_uuid="WVIH97-…"} 9.14358272e+08
lvm_pv_free_bytes{pv_uuid="yc1zVe-…"} 1.040187392e+09

lvm_lv_size_bytes{lv_uuid="BUfEXc-…"} 1.30023424e+08
lvm_lv_size_bytes{lv_uuid="ijb7Yx-…"} 4.194304e+06
```

More examples can be found in the [testdata files](./testdata/).

[lvm2]: https://sourceware.org/lvm2/
[wikipedia]: https://en.wikipedia.org/wiki/Logical_Volume_Manager_(Linux)
[toolkit]: https://github.com/prometheus/exporter-toolkit
[toolkitconfig]: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
[releases]: https://github.com/hansmi/prometheus-lvm-exporter/releases/latest
[golang]: https://golang.org/
[goreleaser]: https://goreleaser.com/

<!-- vim: set sw=2 sts=2 et : -->
