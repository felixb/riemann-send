# riemann-send
Send events to riemann. It's that simple.

## Configuration

Default riemann server configuration is read from a json file stored in `/etc/riemann-send/server.json`.
Overwrite the location with the help of the evnironment variable `RIEMANN_SEND_SERVER_CONFIG`.

Please see [the example server config](example/server-config.json).

## Usage

### Send an event

```bash
riemann-send \
  -host srv.example.org \
  -service "example service" \
  -state ok \
  -metric 13.36 \
  -description "everything is fine"
```

### Send an json formatted event

```bash
riemann-send -json example/event.json # from file
riemann-send -json - < example/event.json # or stdin
```
