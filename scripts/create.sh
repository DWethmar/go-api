#!/bin/bash

cat $(dirname $(realpath $0))/new-content-item.json | http -f POST localhost:8080
