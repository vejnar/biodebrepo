#!/bin/bash

set -e

export KEYNAME="FC1FE21271936D93DD086699A7C205953C4DCCA5"

if [ ! -d "pkgbuilds" ]; then
    echo "ERROR: pkgbuilds not found"
    exit 1
fi
if [ ! -d "repo" ]; then
    mkdir repo
fi

# Copy packages
cp -pf pkgbuilds/*/*.deb repo/
# Copy repo ID
cp -pf my_repo_id/* repo/

# Move to repo
cd repo

# Packages & Packages.gz
dpkg-scanpackages --multiversion . > Packages
gzip --keep --force -9 Packages

# Release, Release.gpg & InRelease
apt-ftparchive release . > Release
gpg --yes --default-key "${KEYNAME}" -abs --output Release.gpg Release
gpg --yes --default-key "${KEYNAME}" --clearsign --output InRelease Release
