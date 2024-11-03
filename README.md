# Bio-packages repo

A public repo of packages for biology complementing [Debian](https://www.debian.org), [Ubuntu](https://ubuntu.com) and [Debian Med](https://www.debian.org/devel/debian-med/).

## Rationale

* Packages complementing [Debian Med](https://www.debian.org/devel/debian-med/) for Debian or Ubuntu
* Building packages on Debian/Ubuntu using the excellent and simple [makepkg](https://gitlab.archlinux.org/pacman/pacman/blob/master/scripts/makepkg.sh.in) from [Archlinux](https://archlinux.org) and the `PKGBUILD`s from the [Arch User Repository (AUR)](https://aur.archlinux.org). In addition:
    1. Wrapper of [Archlinux](https://archlinux.org) [pacman](https://pacman.archlinux.page)
    2. Script `makedeb` wrapping `makepkg` to build *.deb* packages
* Building packages in isolation using container ([Apptainer](https://apptainer.org) is currently supported)
* Public repo

* Sources
    * [Assaf Morami](https://assafmo.github.io/2019/05/02/ppa-repo-hosted-on-github.html)
    * [Arch User Repository (AUR)](https://aur.archlinux.org)

## Packages

List of available [packages](pkgbuilds).

## Usage

### Debian 12 (*bookworm*)

```bash
sudo wget -O - https://github.com/vejnar/biodebrepo/releases/latest/download/key.gpg | gpg --dearmor -o /etc/apt/trusted.gpg.d/biodebrepo.gpg
sudo wget -O /etc/apt/sources.list.d/biodebrepo.list https://github.com/vejnar/biodebrepo/releases/latest/download/biodebrepo.list
sudo apt update
sudo apt install bowtie2 star-cshl=2.7.11b-1
```

## Build package

### Install

1. Requirements
    * For Archlinux
        ```bash
        pacman -S apptainer apt debootstrap dpkg gh
        ```
2. Build Apptainer container
    ```bash
    cd container
    apptainer build debian-makepkg.sif debian-makepkg.def
    ```
3. Configure [gh](https://cli.github.com) for your repo
4. Set the repo URL in `my_repo_id/biodebrepo.list`
5. Export GPG key in `my_repo_id/key.gpg`

### Manual

To manually compile a package, for example [Bowtie2](https://github.com/BenLangmead/bowtie2)
    1. Navigate to direcotry with `PKGBUILD`
        ```bash
        cd pkgbuilds/bowtie2_2.5.4/
        ```
    2. Execute `makedeb` in container
        ```bash
        apptainer exec --fakeroot --writable-tmpfs debian-makepkg.sif makedeb
        ```

### Create repo

1. Install requirements and build container
2. Build all packages
    ```bash
    cd pkgbuilds/
    ../utils/build_repo.sh
    ```
3. Build repo
    ```bash
    cd ..
    ./utils/update_repo.sh
    ```
3. Upload repo
    ```bash
    ./utils/upload_repo.sh
    ```
