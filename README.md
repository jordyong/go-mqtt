## Setup WSL port forwarding
Run this cmd on windows terminal with adminstrator
```
$env:WSL_LISTENPORT = "8080"
$env:NET_LISTENPORT = "8080"
$WSL_ADDR = wsl -d debian hostname -I
netsh advfirewall firewall add rule name="Allowing LAN connections" dir=in action=allow protocol=TCP localport=$WSL_LISTENPORT
netsh interface portproxy add v4tov4 listenaddress=localhost listenport=$WSL_LISTENPORT connectaddress=$WSL_ADDR connectport=$NET_LISTENPORT

```

## Connect to WSL Webserver
Use the ip addr of windows machine with the listening port
eg. <http://192.168.96.75:8080>

