# Dashboard Client

The idea of this repo is to be a redesign of the remotechrome code. Each TV will run Consul which will connect to the admin server as well as this application. This app watches paths in Consul and if they change it will call Chromium dev tools to either open a URL or reload the browser.

Previously, I thought about using confd to watch Consul and then have confd execute a script but I realized that this app could do the same work with less moving parts.

## Configuration

```shell
[chrome]
host = "localhost"
port = 9222

[consul]
address = "localhost:8500"
scheme = "http"
datacenter = "datacenter1"
action = "foo/action"
newurl = "foo/url"
runningurl = "foo/current"

[delay]
interval = 1000
```

### Config Section - Chrome

The **address** and **port** the Chromium dev tools are listening on. localhost and 9222 are the defaults.

### Config Section - Consul

**address** is the address and port of the local Consul server.

**scheme** is the http protocol used to talk to the Consul server. Since it is only talking to the localhost, http is being used.

**datacenter** is the Consul configured datacenter

**action** is the path the application will watch for action commands. If the application were called **tv1** the path would be **tv1/action**. The TV name is where the dashboard-admin will write the action. At the moment the only valid actions are **open** and **reload**

**newurl** is the path the application will watch for what URL to open in the browser. If the application were called **tv1** the path would be **tv1/url**. The application keeps the currently running URL stored in **runningURL** and with compare the value stored in **tv1/url** to **runningURL** if they match the URL will not be loaded. If they do not match the URL will be opened in the browser.

**runningurl** is the path the application will write the currently running URL. It is possible that the new URL and the running URL could be different. A future plan is to setup a temp URL path so that a URL would be shown for some short amount of time and after that time the value of runningurl could be used to reopen the normal path the TV should show.

### Config Section - Delay

**interval** the number of milliseconds the application will wait before looking at the Consul paths to see if they have changed.

## Execution

```shell
% ./dashboard-client -help
Usage of ./dashboard-client:
-conf string
Config file chromium, Consul, delay interval.

% ./dashboard-client -conf config.toml
```

## Starting Chromium

```shell
#!/usr/bin/env bash
CHROME_DATA_DIR=$(mktemp -d)
trap "rm -rf ${CHROME_DATA_DIR}" SIGINT SIGTERM EXIT

# this will take over your screen
/usr/bin/chromium --remote-debugging-port=9222 --user-data-dir="${CHROME_DATA_DIR}" --disable-infobars --kiosk "about:blank"

# to test without it chromium going full screen
#/usr/bin/chromium --remote-debugging-port=9222 --user-data-dir="${CHROME_DATA_DIR}"
```

## References

### Dependencies

<https://github.com/BurntSushi/toml>

<https://github.com/hashicorp/consul/api>

<https://github.com/raff/godet>

<https://consul.io>

### Code being superseded

<https://github.com/mikeplem/remotechrome>

## Testing

### Consul Setup

Here is a simple Consul config

```shell
{
"datacenter": "datacenter1",
"data_dir": "/tmp/consul",
}
```

```shell
mkdir -p /tmp/consul
```

Run Consul like this

```shell
./consul agent -dev -config-file test.hcl
```

Since we are in dev mode nothing exists in Consul so create test values

```shell
./consul kv put foo/url https://osu.edu
./consul kv put foo/action open
```

### Chromium Setup

Use the script above to start Chromium but make sure it does not go full screen

### Dashboard Client Execution

Once you start the dashboard client it will run forever until you hit Ctrl+c

```shell
./dashboard-client -conf config.toml
```

You should see Chromium open up Ohio State's website.

### Updating Chromium

With all the dependent software running we can now update the browser using Consul.

### Open A New URL

```shell
./consul kv put foo/url https://github.com
```

It should not be necessary to have to send the open command but if the page does not open send the open action

```shell
./consul kv put foo/action open
```

### Reload The Browser

```shell
./consul kv put foo/action reload
```

### Client Test Output

```shell
% ./dashboard-client -conf config.toml
2020/01/04 11:04:52 Chrome Connection:  localhost:9222
2020/01/04 11:04:52 Consul Address:  localhost:8500
2020/01/04 11:04:52 Consul Scheme:  http
2020/01/04 11:04:52 Consul Datacenter:  datacenter1
2020/01/04 11:04:52 Consul Action Path:  foo/action
2020/01/04 11:04:52 Consul URL Path:  foo/url
2020/01/04 11:04:52 Loop Delay:  1s
2020/01/04 11:04:52 Open URL:  https://osu.edu
2020/01/04 11:04:52 Requested to open https://osu.edu
2020/01/04 11:05:08 Reload browser
2020/01/04 11:05:08 Reloading browser page
2020/01/04 11:05:20 Open URL:  https://github.com
2020/01/04 11:05:20 Requested to open https://github.com
```
