#!/bin/bash

http -f GET :8080/ | jq '.[-1]'
