# devkit
## devkit-client
### 本地接口
* 推送文件：
``` http
POST /file/push?address=<target.server-address>&source=<source.file-path>&target=<target.file-path>&override={target.file-override} HTTP/1.1

<optional.source.file-content>
```
## devkit-server
配合 `devkit` 命令可实现：
* pull - 从 devkit-client 拉取用户选定的文件
* push - 向 devkit-client 推送指定的文件内容

## devkit-relay

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
keyfile=var/cert/server.key
certfile=var/cert/server.crt
openssl req -newkey rsa:2048 -x509 -nodes -keyout "$keyfile" -new -out "$certfile" -subj /CN=$name

keyfile=var/cert/client.key
certfile=var/cert/client.crt
openssl req -newkey rsa:2048 -x509 -nodes -keyout "$keyfile" -new -out "$certfile" -subj /CN=$name


```
