# Manager for YeeLight (wifi) devices

This is part of Maja suite project.

Module connects supported YeeLight WiFi devices and exchange requests with *mqtt* server.

I assume devices was previously initialized by mobile app or *xiaomi_reg* command. When registration procedure finish device_id and device_token will be shared and stored in mqtt database. After application starts and connect to mqtt server all registered devices datas will be send to application.

## Compile

Run *go build* command. As result, you will have one binary file to run.

## Notice

Author: Eugene Chertikhin <e.chertikhin@crestwavetech.com>

Licensed under GNU GPL.
