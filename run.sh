#!/bin/bash

export PRODUCTION=1
PJ=unitySample
FILE=unitySample

DIR=`dirname $0`

start() {
    cd ${DIR}/bin;
    nohup ./${FILE} > /tmp/${FILE}.out 2>&1&
    touch ${DIR}/${FILE}.pid
    echo $! > ${DIR}/${FILE}.pid
}

stop() {
    cd ${DIR}/bin;
    kill `cat ${DIR}/${FILE}.pid`
}


case $1 in
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
  *)
  echo "usage: run.sh {start|stop|restart}" ;;
esac
exit 0
