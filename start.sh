#!/bin/sh
#DOCKER STARTUP SCRIPT

/usr/bin/couchdb -a /etc/couchdb/default.ini -a /etc/couchdb/local.ini -b -r 5 -p /var/run/couchdb/couchdb.pid -o /dev/null -e /dev/null -R &

sleep 5s

/go/bin/simple-todo 
