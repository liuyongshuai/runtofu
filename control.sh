#! /bin/bash

set -e

CONF_NAME="service.conf"
SERVICE_NAME="%(APP)"
BINARY_NAME="%(APP)"
LOG_DIR="/home/xiaoju/soda/log/search-broker_service"
APOLLO_DIR="/home/xiaoju/ep/as/store"

CUR_DIR=$(dirname "$0")
MY_PATH=$(cd "$CUR_DIR" && pwd -P)

function setLogPath() {
    mkdir -p "$LOG_DIR"
}

function setApolloPath() {
    mkdir -p "${APOLLO_DIR}"
    mkdir -p "${APOLLO_DIR}/toggles"
    mkdir -p "${APOLLO_DIR}/conf"
}

function setConfigFile() {
    echo "pwd: $MY_PATH"

	CLUSTER_NAME=$(cat $MY_PATH/.deploy/service.cluster.txt)

	if [ -f "$MY_PATH/conf/$CONF_NAME.$CLUSTER_NAME" ]; then
        rm -f "$MY_PATH/conf/$CONF_NAME"
       	ln -s "$MY_PATH/conf/$CONF_NAME.$CLUSTER_NAME" "$MY_PATH/conf/$CONF_NAME"
	else
		echo "Config file is not found for cluster: $CLUSTER_NAME."
		exit 1
	fi
}

function start() {
	setConfigFile
    setLogPath
    setApolloPath
	exec "$MY_PATH/bin/$BINARY_NAME" -config="$MY_PATH/conf/$CONF_NAME"
}

function stop() {
    supervisorctl stop $SERVICE_NAME
}

function restart() {
    stop
    sleep 1
    start
}

function usage() {
    echo "Usage: $0 {start|stop|restart}"
    exit 1
}

if [ $# != 1 ]; then
    usage
fi

case "$1" in
    start|stop|restart)
        $1
        ;;
    *)
        usage
esac
