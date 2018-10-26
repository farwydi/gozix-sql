# GoZix SQL

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
