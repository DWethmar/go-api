#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
cat ${SCRIPTPATH}/new-content-item.json | http -f POST localhost:8080
