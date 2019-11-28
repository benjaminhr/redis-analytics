# Redis Analytics

NOTE: Accepts all WS connections from all origins

Usage: `./wss -h 127.0.0.1 -p 6379 -db 0 -chan browsers`

All flags are optional:
- `-h` host
  - Default value: 127.0.0.1
- `-p` port
  - Default value: 6379
- `-pass` Redis password
  - Default value: <empty_string>
- `-db` Redis DB index
  - Default value: 0
- `-chan` Channel to listen to
  - Default value: browsers

To recompile `go build -o wss .`