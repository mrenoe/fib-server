# Fib-Server

## Assumptions and Caveats
* This server uses a global counter to track where the server is at in the fibonacci sequence - this has numeric limits. i.e. for fib(n) = x, n is between 0 and 18446744073709551615 and x is between 0 and 2^63 (machine dependent)
* Future versions could use a memcache/redis store to maintain sessions and users so that each user can track `fib(n)` separately

## Benchmarks
Current benchmarks using [`wrk`](https://github.com/wg/wrk):
Test: `wrk -t12 -c500 -d20s http://<hosted-fib-server>/next`
```
Running 20s test @ http://<hosted-fib-server>/next
  12 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   103.92ms  163.41ms   1.86s    96.19%
    Req/Sec   312.98    302.56     1.68k    80.59%
  56910 requests in 20.10s, 331.02MB read
  Socket errors: connect 251, read 0, write 0, timeout 0
Requests/sec:   2831.04
Transfer/sec:     16.47MB
```
Given that there are errors, there's certainly places to improve the server infrastructure and likely configuration

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
This is designed to run on linux machine with 512mb RAM, 1 CPU, and port 80 open to HTTP traffic. The `dist/` folder contains a binary ready to be deployed in that way.


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
