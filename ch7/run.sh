#!/bin/bash

./dummy_cpm &
DUMMYCPM_PID=$!

#./rest_server  &
./rest_server -cpuprofile=restserver_cpu.prof &
#./rest_server -memprofile=restserver_mem.prof &
RESTSERVER_PID=$!

sleep .25

oper=""
for i in {1..1000}; do
    if [ "x$oper" = "xstart" ]; then
        oper="stop"
    else
        oper="start"
    fi
    curl -X PUT -d "$oper" http://localhost:8080/wallbox/123 &
    sleep .01
done

kill -SIGTERM $RESTSERVER_PID

sleep .25

kill -SIGTERM $DUMMYCPM_PID

if [ -f /tmp/cpm.sock ]; then
  rm /tmp/cpm.sock
fi
