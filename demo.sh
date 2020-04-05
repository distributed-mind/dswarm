#!/usr/bin/env bash

killall dswarm &>/dev/null
# ./build.sh

for i in {1..10} ; do
    ./dswarm & 2>&1
done

sleep 600

killall dswarm
