#!/bin/bash

GOOS=linux GOARCH=arm GOARM=6 go build
scp HomeAuto pi@192.168.0.100:~
