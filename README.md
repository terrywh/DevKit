# devkit

## dependencies
### bun
> https://bun.sh
```
powershell -c "irm bun.sh/install.ps1 | iex"
```

### npm
``` bash
bun install
```

### authorize
``` bash
mkdir -p etc
touch etc/devkit.yaml
```

### certificate
``` bash
mkdir -p var/cert
name=`hostname`

keyfile=var/cert/client.key
certfile=var/cert/client.crt
openssl req -newkey rsa:2048 -x509 -nodes -keyout "$keyfile" -new -out "$certfile" -subj /CN=$name

keyfile=var/cert/server.key
certfile=var/cert/server.crt
openssl req -newkey rsa:2048 -x509 -nodes -keyout "$keyfile" -new -out "$certfile" -subj /CN=$name
```
