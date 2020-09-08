#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
cat ${SCRIPTPATH}/new-content-item.json | http --verify=no -f POST https://go-api.192.168.88.100.nip.io/content
