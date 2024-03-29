# devkit


### certificate
```
mkdir -p var/cert
keyfile=var/cert/server.key
certfile=var/cert/server.crt

openssl req -newkey rsa:2048 -x509 -nodes -keyout "$keyfile" -new -out "$certfile" -subj /CN=localhost
```