# Manager for YeeLight (wifi) devices

## Instead of preface

Manager is a driver module to serve requests and statuses from different hardware (yeelight in current case). It can receive
status notification and send commands from/to all yeelight lamps and probably some other yeelight products (need test).

All updates and commands send thru mqtt server. I assume different module (named hub) should control mqtt traffic and
serve automation.

## How to YeeLight (and this manager) works (in short)

Operation principle: device discovered in the network by ssdp protocol with little correction from yeelight. In yeelight
implementation used different port for discovery (1982). Discovered devices saves in mqtt server with retain flag.
After device found it stored in module internal memory.  If devices received from mqtt discovery process assign ip addresses
to device. It can take approx. one minutes.
All yeelight hardware often close connections. Its normal. So manager doen't keep connection to device open every time.

## Supported devices

 * Probably any type of ceiling light (test 3 different types)
 * Light strip
 * RGB and basic bulbs and filament as well
 
## Command line paraments
```commandline
$ ./manager_yeelight -?
flag provided but not defined: -?
Usage of ./manager_yeelight:
  -clientid string
    client id for mqtt server (default "yeelight-1")
  -qos int
    qos for messages send to mqtt
  -debug
    print debuging hex dumps
  -keepalive int
    keepalive timeout for mqtt server (default 30)
  -login string
    login string for mqtt server
  -mqtt string
    mqtt server address (default "127.0.0.1:1883")
  -pass string
    password string for mqtt server
```

# Device registration

Module manager_yeelight doen't have its own registration procedure. Need to use manager_xiaomi to do this. In most cases 
manager_xiaomi is compatible and interchangeble with manager_yeeligh. But manager_yeelight use more simple (and sometimes 
more fast communications method).

## Format mqtt strored object

Each device found in the network stored in mqtt as retain message in follow format:

```commandline
yeelight/xxxxxxxx = {"type":1,"model":"ceiling1","id":"xxxxxxxx","name":"New light","support":"get_prop set_default set_power toggle set_bright set_scene cron_add cron_get cron_del start_cf stop_cf set_ct_abx set_name set_adjust adjust_bright adjust_ct","ver":26}
```

Ip address of device doen't stored because address can change by DHCP server in the network. Manager_yeelight doesn't 
requre to use static address for device and address changes serve correctly.

When device found in the network mqtt record will be extended with support commands tag (support) and status variables. 
Status variables set may (and actually should) be different on different types of device. (Lamps do not inform about 
temperature, but thermometers does).

## The commands from mqtt server to manage device

Two different method is possible.

First one - take the last message from mqtt and change some (several) parameters. Very easy and universal.

Next - run methods according yeelight protocol. Method name can be taken from "support" array. "action", "effect" and 
"duration" is method arguments

Different devices has different set of commands. Tag "support" in device record define which command are actually
supported by this device. Hub or process who format command should know set of parameters for each command.

Command format
```commandline
{"type":1,"model":"ceiling1","method":"toggle","effect":"","duration":"","action":""}
```

## Known problems

.... still in early development stage.

## License and author

This project licensed under GNU GPLv3.

Author Eugene Chertikhin
