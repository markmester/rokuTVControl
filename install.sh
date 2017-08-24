#!/usr/bin/env bash

# Testing on Raspberry Pi Zero with 2017-07-05-raspbian-jessie-lite
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

echo ">>> Installing redis-server..."
apt update && apt install redis-server -y
sed -i "s/daemonize.*/daemonize no/g" /etc/redis/redis.conf

echo ">>> Installing supervisor..."
apt install supervisor
cp supervisor.conf /etc/supervisord/conf.d/

echo ">>> Creating log files..."
touch /var/log/redis.log


echo ">>> Updating supervisor..."
supervisorctl reread && supervisorctl reload
