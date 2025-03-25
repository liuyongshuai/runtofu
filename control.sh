#!/bin/bash

workspace=$(cd $(dirname $0) && pwd -P)

action="$1"
conf="./conf/service.conf"
mkdir log

function usage() {
        echo "manage runtofu process"
        echo
        echo "usage: $0 start|stop|restart|status"
        echo
        echo "    conf   optional, default ./conf/service.conf"
        echo
}

function start() {
		nohup bin/runtofu -config=./conf/service.conf > ./abc.txt 2>&1 &
		status
}

function stop() {
	pkill -f 'bin/runtofu -config=./conf/service.conf'
	kill -9 $(ps aux|grep runtofu|grep -v grep|awk '{print $2}')
    sleep 0.5
}

function status() {
	pid=$(pgrep -f 'bin/runtofu -config=./conf/service.conf')
        if [ "x$pid" == "x" ]; then
            echo "runtofu is not running"
        else
            echo "runtofu is running with pid [$pid]"
        fi
}

case $action in
    start)
		start
		;;
    stop)
		stop
		;;
    restart)
		stop
		start
		;;
    status)
		status
		;;
    *)
		usage
		;;
esac