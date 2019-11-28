# rallytics

Usage: `./rallytics [options, see below]`

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
- `-wsport` Websocket server port number
  - Default value: 8080

## testing

    docker-compose up --build --force-recreate --scale rallytics=3
    # (or be awesome and use https://github.com/matti/doc)

    open http://rallytics.localtest.me
