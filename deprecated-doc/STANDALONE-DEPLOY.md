# Deploy a Standalone Canopsis from Packages

The package can be obtained from a full build, see [documentation](/doc/docs/fr/guide_administrateur/packaging.md).

After building a release and packages, to ease you the job of copy pasting the following commands:

```bash
pkgver=<your_version_built>
pkgrel=<your_release_1_by_default>
```

## CentOS/RedHat 7

### CORE

```
yum install -y epel-release
yum localinstall -y canopsis-core-${pkgver}-${pkgrel}.el7.x86_64.rpm
canoctl deploy
```

### CAT

```
yum install -y epel-release
yum localinstall -y canopsis-cat-${pkgver}-${pkgrel}.el7.x86_64.rpm
canoctl deploy
```

## Debian 8 / Jessie

### CORE

```
dpkg -i canopsis-core-${pkgver}-${pkgrel}.amd64.jessie.deb
apt install -y -f
canoctl deploy
```

### CAT

```
dpkg -i canopsis-cat-${pkgver}-${pkgrel}.amd64.jessie.deb
apt install -y -f
canoctl deploy
```

## Debian 9 / Stretch

### CORE

```
dpkg -i canopsis-core-${pkgver}-${pkgrel}.amd64.stretch.deb
apt install -y -f
canoctl deploy
```

### CAT

```
dpkg -i canopsis-cat-${pkgver}-${pkgrel}.amd64.stretch.deb
apt install -y -f
canoctl deploy
```

# Develop in a standalone installation

```
chown -R canopsis:canopsis /opt/canopsis/lib/python2.7
su - canopsis
pip install -U . /path/to/canopsis/sources/canopsis
```
