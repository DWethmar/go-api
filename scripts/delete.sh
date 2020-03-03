#!/bin/bash

LASTID=$(http -f GET localhost:8080 | jq '.[-1].id')

http -f DELETE localhost:8080/$LASTID