#!/bin/bash
DATA="{\"pid\":$1,\"path\":\"/Users/terryhaowu/Downloads\"}"
echo $DATA
curl -d $DATA http://127.0.0.1:18081/file/pull