#!/bin/sh
cd /srv

# Extract raspios.img into base/
[ -d "base" ] && rm -rf base
mkdir -p base
mkdir -p bootpart rootpart
loopback=$(losetup -fP --show raspios.img)
mount ${loopback}p1 bootpart
mount ${loopback}p2 rootpart
rsync -a rootpart/ base
rsync -a bootpart/ base/boot/firmware
umount bootpart
umount rootpart
losetup -d ${loopback}
