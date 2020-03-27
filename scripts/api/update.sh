#!/bin/bash

# cat ./updated-content-item.json | http -f POST localhost:8080/6

LASTID=$(http -f GET localhost:8080 | jq '.[-1].id')

cat ./updated-content-item.json | http -f POST localhost:8080/$LASTID