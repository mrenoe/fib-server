# Fib-Server

## Fib-Server is a service API that exposes 3 primary endpoints:
`/current`
Returns the current number in the servers fibonacci sequence

`/previous`
Returns the previous number in the fibonacci sequence from current

`/next`
Returns the next number in the fibonacci sequence from current and increments the central counter by 1
*Note:* If the server count reaches the `uint64` max(18446744073709551615), it will reset back to 0


## Deploying

### Binary
This is designed to run on linux machine with 512mb RAM and 1 CPU and the `dist/` folder contains a binary ready to be deployed in that way


### Source Build
Clone the repo to your preferred destination
* Ensure you have go1.13 or higher installed
* Run `go mod download` to download the modules locally
* Run `go build main.go` or if you would prefer a custom named binary `go build -o <appname>`
* For statically linked binaries to ensure portability users can run `env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o <appname>` to produce a linux ready binary

### Running the server
Once you have your binary, navigate to your host machine and within `/etc/systemd/system/` run the command `touch fib.service` in order to create a systemd service entry to manage runtime of the server

A sufficient service file will look as such:
```
[Unit]
Description=Fibonacci Tracker API
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=<binary-directory>/<appname> -addr=:80 -compress=true

[Install]
WantedBy=multi-user.target
```

To start the server from here, run `systemctl start fib.service`
