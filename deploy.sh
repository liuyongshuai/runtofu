#!/bin/bash

git pull
make
cd output
cp -rf ~/service.conf ./conf/
./control.sh restart

