#!/bin/bash

# cat ./updated-content-item.json | http -f POST localhost:8080/6

LASTID=$(http -f GET localhost:8080 | jq '.[-1].id')
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

cat ${SCRIPTPATH}/update-content-item.json | http -f POST localhost:8080/$LASTID
