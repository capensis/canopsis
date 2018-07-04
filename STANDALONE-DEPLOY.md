# Deploy a Standalone Canopsis from Packages

## CentOS/RedHat 7

### CORE

```
yum install -y epel-release
yum localinstall -y canopsis-core-2.7.0-1.el7.x86_64.rpm
canoctl deploy
```

### CAT

```
yum install -y epel-release
yum localinstall -y canopsis-cat-2.7.0-1.el7.x86_64.rpm
canoctl deploy
```

## Debian 8 / Jessie

### CORE

```
dpkg -i canopsis-core-2.7.0-1.amd64.jessie.deb
apt install -y -f
canoctl deploy
```

### CAT

```
dpkg -i canopsis-cat-2.7.0-1.amd64.jessie.deb
apt install -y -f
canoctl deploy
```

## Debian 9 / Stretch

### CORE

```
dpkg -i canopsis-core-2.7.0-1.amd64.stretch.deb
apt install -y -f
canoctl deploy
```

### CAT

```
dpkg -i canopsis-cat-2.7.0-1.amd64.stretch.deb
apt install -y -f
canoctl deploy
```
