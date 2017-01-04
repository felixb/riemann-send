# riemann-send
Send events to riemann. It's that simple.

## Configuration

### Server config

Default riemann server configuration is read from a json file stored in `/etc/riemann-send/server.json`.
Overwrite the location with the help of the environment variable `RIEMANN_SEND_SERVER_CONFIG`.

Please see [the example server config](example/server-config.json).

### Event template

`riemann-send` uses a json stored in `/etc/riemann-send/event.json` as template for default values for each event.
Overwrite the location with the help of the environment variable `RIEMANN_SEND_EVENT_CONFIG`.

Tags and attributes specified on command line are merged/appended to the tags and attributes specified in the template.

Please see [the example event template](example/event.json).

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
