#!/usr/bin/env bash

# Testing on Raspberry Pi Zero with 2017-07-05-raspbian-jessie-lite
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

echo "Installing LIRC with IR Tx at GPIO 22 and IR Rx at GPIO 23..."
apt update && apt install lirc -y

echo "Updating necessary configs..."

cat >> /etc/modules <<EOF
lirc_dev
lirc_rpi gpio_in_pin=23 gpio_out_pin=22
EOF

cat > /etc/lirc/hardware.conf <<EOF
########################################################
# /etc/lirc/hardware.conf
#
# Arguments which will be used when launching lircd
LIRCD_ARGS="--uinput"
# Don't start lircmd even if there seems to be a good config file
# START_LIRCMD=false
# Don't start irexec, even if a good config file seems to exist.
# START_IREXEC=false
# Try to load appropriate kernel modules
LOAD_MODULES=true
# Run "lircd --driver=help" for a list of supported drivers.
DRIVER="default"
# usually /dev/lirc0 is the correct setting for systems using udev
DEVICE="/dev/lirc0"
MODULES="lirc_rpi"
# Default configuration files for your hardware if any
LIRCD_CONF=""
LIRCMD_CONF=""
########################################################
EOF

cat >> /boot/config.txt <<EOF
dtoverlay=lirc-rpi,gpio_in_pin=23,gpio_out_pin=22
dtparam=gpio_in_pull=up
EOF

echo ">>> Installing redis-server..."
apt update && apt install redis-server -y
sed -i "s/daemonize.*/daemonize no/g" /etc/redis/redis.conf

echo ">>> Installing supervisor..."
apt install supervisor
cp supervisor.conf /etc/supervisord/conf.d/

echo ">>> Creating log files..."
touch /var/log/lirc_server.log && touch /var/log/redis.log


echo ">>> Updating supervisor..."
supervisorctl reread

echo ">>> Done installation; rebooting in 10 seconds..."
sleep 10
reboot
