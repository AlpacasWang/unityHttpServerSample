#!/bin/bash

export PRODUCTION=1
PJ=unitySample
FILE=unitySample

start() {
    cd /home/game/git/${PJ}/bin;
    exec nohup /home/game/git/${PJ}/bin/${FILE} > /tmp/${FILE}.out 2>&1&
    echo $! > /home/game/pids/${FILE}.pid
    disown
}

stop() {
    kill `cat /home/game/pids/${FILE}.pid`
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
