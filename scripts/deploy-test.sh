#!/bin/bash
./scripts/build.sh
7z a ./tusd_linux_amd64/tusd.zip ./tusd_linux_amd64/tusd
scp ./tusd_linux_amd64/tusd.zip linode-st:~/goapps/temp/
sh /home/samtech/Develop/scripts/deploy-test.sh 'tusd'
echo "Done"
