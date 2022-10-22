# GoZix SQL

[documentation-img]: https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff
[documentation-url]: https://pkg.go.dev/github.com/gozix/sql/v3
[license-img]: https://img.shields.io/github/license/gozix/sql.svg?style=for-the-badge
[license-url]: https://github.com/gozix/sql/blob/master/LICENSE
[release-img]: https://img.shields.io/github/tag/gozix/sql.svg?label=release&color=24B898&logo=github&style=for-the-badge
[release-url]: https://github.com/gozix/sql/releases/latest
[build-status-img]: https://img.shields.io/github/actions/workflow/status/gozix/sql/go.yml?logo=github&style=for-the-badge
[build-status-url]: https://github.com/gozix/sql/actions
[go-report-img]: https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge
[go-report-url]: https://goreportcard.com/report/github.com/gozix/sql
[code-coverage-img]: https://img.shields.io/codecov/c/github/gozix/sql.svg?style=for-the-badge&logo=codecov
[code-coverage-url]: https://codecov.io/gh/gozix/sql

[![License][license-img]][license-url]
[![Documentation][documentation-img]][documentation-url]

[![Release][release-img]][release-url]
[![Build Status][build-status-img]][build-status-url]
[![Go Report Card][go-report-img]][go-report-url]
[![Code Coverage][code-coverage-img]][code-coverage-url]

The bundle provide a SQL integration to GoZix application.

## Installation

```shell
go get github.com/gozix/sql/v3
```

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "sql": {
    "default": {
      "nodes": [
        "postgres://app:password@127.0.0.1:5432/app?sslmode=disable"
      ],
      "driver": "postgres",
      "max_open_conns": 10,
      "max_idle_conns": 10,
      "conn_max_lifetime": "10m"
    }
  }
}
```

## Documentation

You can find documentation on [pkg.go.dev][documentation-url] and read source code if needed.

## Questions

If you have any questions, feel free to create an issue.
