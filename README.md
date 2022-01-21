# Manager for YeeLight (wifi) devices

## Instead of preface

Manager is a driver module to serve requests and statuses from different hardware (yeelight in current case). It can receive
status notification and send commands from/to all yeelight lamps and probably some other yeelight products (need test).

All updates and commands send thru mqtt server. I assume different module (named hub) should control mqtt traffic and
serve automation.

## How to YeeLight (and this manager) works (in short)

Operation principle: device discovered in the network by ssdp protocol with little correction from yeelight. In yeelight
implementation used different port for discovery (1982). Discovered devices saves in mqtt server with retain flag.
After device found it store in module internal memory and two connection will be opened. One connection for sending commands
and another one to listen status changes.
All yeelight hardware often close connections. Its normal. Module check it and reopen. In case of send command to device and
device close connection we try to reopen it in one minute. 

## Supported devices

 * Probably any type of ceiling light (test 3 different types)
 * Light strip
 
## Command line paraments

$ ./manager_yeelight -?
flag provided but not defined: -?
Usage of ./manager_yeelight:
  -clientid string
    client id for mqtt server (default "yeelight-1")
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

# Device registration

Module manager_yeelight doen't have its own registration procedure. Need to use manager_xiaomi to do this. In most cases 
manager_xiaomi is compatible and interchangeble with manager_yeeligh. But manager_yeelight use more simple (and sometimes 
more fast communications method).

## Format mqtt strored object

Each device found in the network stored in mqtt as retain message in follow format:

yeelight/78ab4e4 = {"id":78ab4e4,"model":"ceiling24","name":"My ceiling light","ver":5}

Ip address of device doen't stored because address can change by DHCP server in the network. Manager_yeelight doesn't 
requre to use static address for device and address changes serve correctly.

When device found in the network mqtt record will be extended with support commands tag (support) and status variables. 
Status variables set may (and actually should) be different on different types of device. (Lamps do not inform about 
temperature, but thermometers does).

Extended example of device records in mqtt:

yeelight/78ab4e4 = {"id":78ab4e4,"ip":"1.1.1.1","model":"ceiling24","name":"My ceiling light","ver":5,
    "support":"set_scene set_ct_abx adjust_ct set_bright set_name adjust_bright set_default toggle cron_get start_cf set_adjust get_prop set_power cron_add cron_del stop_cf",
    "power":"false","bright":100,"mode":2,"temp":4000,"rgb":0,"hue":0,"sat":0}

## The commands from mqtt server to manage device

Command format: { "cmd":"xxxx", "value1":"aaa", "value2":"bbbb", "value3":"ccccc"}

Commands are dynamically formatted and processed. Different devices has different set of commands. Tag *support* in device 
record define which command are actually supported by this device. Hub or process who format command should know set of 
parameters for each command.

## Known problems

* mqtt client has several (or much more) issues. i.e. reconnect procedure eat too much cpu time
* All commands/parameters is pass to lamp correctly, but not all types of possible response correctly converted to mqtt device update, i.e. music flow and etc.

## License and author

This project licensed under GNU GPLv3.

Author Eugene Chertikhin
