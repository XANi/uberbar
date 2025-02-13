# Uberstatus

![byzanz-record -x 3065 -y 0 -h 22 -w 775 uberstatus.gif](doc/uberstatus.gif)

Status line generator for i3wm and (eventually) other WMs

Goals:

* integrate most of core functionality directly (checking cpu, network etc. without extra forking
* Support push and pull (polling) type of plugins
* Async (each plugin on its own timer) and sync mode with optional low power/slow update (WIP)

# Installation

If you have go set up already then *usually* just `go get github.com/XANi/uberstatus` (if any of upstream deps didn't break) if not, get the repo and:

    make # will make binary in application root
    mkdir -p ~/.config/uberstatus

Binaries will be made in current directory. If you dont have config already, copy one:

    cp cfg/uberstatus.default.conf ~/.config/uberstatus/uberstatus.conf # copy default config
    
Then just put it in i3 config in section defining bars

    bar {
        ...
        status_command path/to/uberstatus/uberstatus
        ...
    }


# Operation

Each of plugins operates asynchronously and can send update of its status at any time. On each plugin update state is sent upstream (including previous cached plugin state) so it is possible to have each plugin update separetely.

Each plugin have a name (by default plugin name) and instance (in case of plugins that will be used multiple times, like network interfaces).

# Commandline options

* `-config file/name` - use alternate config file
* `-d` - enable remote debugging on localhost:6060

# Plugins

All plugins have `interval` option defining refresh rate in miliseconds

## CPU

* `prefix`

## CPUFreq


## Clock

Parameters:

* `long_format`/`short_format` - time format, as golang formatting string
* `interval` - update interval in miliseconds. Note that it can be longer than it because of run time so for update every second something like 990 ms is better

## Disk free

```yaml
    - name: disk-root
      instance: df
      plugin: df
      config:
       prefix: "💾"
       mounts:
         - /
         - /var
         - /home
```

`prefix` will be added at beginning of the status bar. name and instance are just to distinguish between different instances

## GPU

Gets stats from GPU. So far only NVIDIA is supported via `nvidia-smi` which needs to be installed prior to running the plugin.

```yaml
    - name: gpu
      plugin: gpu
      config:
        id: 0
        type: nvidia
```

## Memory

Click to get detailed stats

## MQTT

Subscribe to MQTT queue

```yaml
    - name: mqtt
      plugin: mqtt
      config:
        # password included in URL
        address: tcp://mqtt:mqtt@192.168.1.4:1883
        lisp_filter: (aget (split (aget (split x ":" ) 1) ".") 0)
        template: "P: {{ color ( intOr0 .Msg | sub100 | percentToColor ) (intOr0 .Msg | percentToBar) }}  {{ color ( intOr0 .Msg | sub100 | percentToColor ) \"%\"}}"
        template_on_click: "P: {{ color ( intOr0 .Msg | percentToColor )  .Msg }} %"
        subscribe: collectd/some-data/phone/power-battery
```
`lisp_filter` is [zygo](https://github.com/glycerine/zygomys) lisp code ran on the data read from `subscribe` channel.
Main use is for formatting data received from the topic defined in `subscribe`

`template` is Go template string working on struct:

```
Event{
		Msg: p.lastMessage,
		TS:  p.lastMQTTUpdate,
	})
```


## Network

Left click for interface's IP, right click for secondary IP (usually IPv6), middle click to display all addresses.

Parameters:

* `iface` - interface to use

## Ping

* `type` - tcp/http
* `addr` - address of a target, host:port format for tcp, url for http(s)

## Pipe

Accepts data to display in pipe. Data is updated instantly.

* `path` - path to named pipe (will be created if not exist
* `parse_template` - enable template parsing, that allows for using template functions like `{{color #00ff00 "message"}}`. Note that some (like `color`) will only work/make sense with markup enabled 
* `markup` - enable pango markup. Enabled by default
* `interval` - regenerate message every x milliseconds. Only matters if you use templates that will change with time.

This will pass data directly to i3 which means that any pango formatting it supports it works but also that you need to escape any HTMLisms (`<>` and such) on your own.

You can also just disable markup (`markup: false`) or use escaping function in template (`{{escape "<title>hai</title>"}}`)

 


## Pomodoro

Simple pomodoro timer. Click button 1 (left) to start/acknowledge break, button 2 to display stats

## I3blocks

i3blocks-compatible input. It will also pass events in compatible way so plugins like volume can be used

Parameters:

* `command` - command to run
* `prefix` - prefix command with text

Example:

```yaml
    - name: volume
      plugin: i3blocks
      config:
        command: /path/to/i3blocks/volume
```

## Icinga

Icinga2 plugin. Icinga2 must have API enabled and served on given URL

Parameters:
```yaml
    - name: icinga
      plugin: icinga
      config:
        url: https://icinga/plugin/url
        user: uber
        pass: status_pass
        # api_update_interval: 5m # defaults to 30s
        # host_filter: icinga2 filter for hosts
        # service_filter: icinga2 filter for services
```

## Uptime

System uptime. Parameters:

* `prefix`

## Weather

Displays temperature, left click for more detailed weather data

To set it up first, get you token [here](https://openweathermap.org/), then get city name or id and set it in location

Example:

```yaml
    - name: weather
      plugin: weather
      config:
        openweather_api_key: 1111cccccccc11111111
        openweather_location: London
```

weather will automatically update every 10 minutes which is way below their free tier ratelimit
