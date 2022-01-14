# Manager for YeeLight (wifi) devices

This is part of Maja suite project.

Module connects supported YeeLight WiFi devices and exchange requests with *mqtt* server.

I assume devices was previously initialized by mobile app or *xiaomi_reg* command. When registration procedure finish 
device_id and device_token will be shared and stored in mqtt database. After application starts and connect to mqtt 
server all registered devices will be send to application.

## Compile

Run *go build* command. As result, you will have one binary file to run.

## Requests from/to mqtt

founded devices stored:

yeelight/78ab4e4 = {"id":78ab4e4,"ip":"1.1.1.1","type":"Celing light","model":"ceiling24","name":"","ver":5,
    "support":"set_scene set_ct_abx adjust_ct set_bright set_name adjust_bright set_default toggle cron_get start_cf set_adjust get_prop set_power cron_add cron_del stop_cf",
    "power":"false","bright":100,"mode":2,"temp":4000,"rgb":0,"hue":0,"sat":0}

commands send from hub :

yeelight/78ab4e4 = {"cmd":"aaaaaa", "value1":"x", "value2":"y", "value3":"z"}

## Notice

Author: Eugene Chertikhin <e.chertikhin@crestwavetech.com>

Licensed under GNU GPL.
