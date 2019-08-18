#!/bin/bash

while true;
do
        server=`ps aux | grep Hole| grep -v grep`
        if [ ! "$server" ]; then
           ./Hole &
        fi
        sleep 300
done
