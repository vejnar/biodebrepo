#!/bin/bash

set -e

if [ "$(whoami)" != "root" ]; then
    echo "ERROR: run as root"
    exit 1
fi

if [ ! -f "PKGBUILD" ]; then
    echo "ERROR: file PKGBUILD not found"
    exit 1
fi

CWD=$(pwd)

# Source the build recipe
source PKGBUILD

# Install dependencies
pacman -S ${depends[@]}
pacman -S ${makedepends[@]}

# Prepare build directory
lfiles=$(ls -1)
mkdir "$CWD/build"
chown build: "$CWD/build"
cp -r $lfiles "$CWD/build"
cd "$CWD/build"

# Build package
sudo -u build makepkg

# Create DEBIAN directory
cd "$CWD/build/pkg"
mkdir -p "${pkgname}/DEBIAN"

cat >${pkgname}/DEBIAN/control << EOF
Package: ${pkgname}
Description: ${pkgdesc}
Homepage: ${url}
Version: ${pkgver}-${pkgrel}
Architecture: amd64
License: ${license}
Maintainer: None
Provides: ${pkgname}
Conflicts: ${pkgname}
EOF

# Create Debian dependency list
deb_depends=$(pacman -C ${depends[@]})
if [ "$deb_depends" != "" ]; then
    echo "Depends: ${deb_depends}" >> ${pkgname}/DEBIAN/control
fi

cat ${pkgname}/DEBIAN/control

# Remove Arch metadata
rm -f ${pkgname}/.BUILDINFO ${pkgname}/.MTREE ${pkgname}/.PKGINFO

# Create Debian package
dpkg-deb --root-owner-group --build ${pkgname}
dpkg-name ${pkgname}.deb

# Cleaning
mv *deb ../..
cd $CWD
rm -fr "$CWD/build"
