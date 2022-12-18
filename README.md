# dilliver

## Development

### Live reload

This project uses [reflex](https://github.com/cespare/reflex) to reload the app
server and recompile assets after changes.

```sh
go install github.com/cespare/reflex@latest
cp reflex.example.conf reflex.conf
reflex -d none -c reflex.conf
```
