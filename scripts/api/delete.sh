#!/bin/bash

LASTID=$(http --verify=no -f GET https://go-api.192.168.88.100.nip.io/content | jq '.[-1].id')

http --verify=no -f DELETE https://go-api.192.168.88.100.nip.io/content/$LASTID