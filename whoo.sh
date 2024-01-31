#!/bin/bash
# Copy files from WSL (working area) to rpi (test server)

rsync -vvru ./* tmiku@rpi:~/go/src/mikuserv