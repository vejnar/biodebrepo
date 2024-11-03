#!/bin/bash

set -e

for fpack in $(ls -1 "."); do
    echo $fpack
    cd $fpack
    if [ -f "$(ls -1 ${fpack}-?_amd64.deb)" ]; then
        echo "$fpack already built"
    else
        echo "$fpack building"
        apptainer exec --fakeroot --writable-tmpfs ../../container/debian-makepkg.sif makedeb
    fi
    cd ..
done
