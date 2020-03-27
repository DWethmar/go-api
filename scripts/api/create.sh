#!/bin/bash

mydir="${0%/*}"

cat $(mydir)new-content-item.json | http -f POST localhost:8080
