# gobar

An i3bar/swaybar protocol-compatible program. Includes some basic modules and
click support.

## Building

```bash
go build ./cmd/gobar
```

## Usage

Place this in your i3 or sway config:

```i3config
bar {
  status_command path/to/gobar
}
```
