# Installation de Canopsis par les paquets système

## Pré-requis

Canopsis doit être installé sur l'un des environnements Linux suivants :
*  CentOS 7 ;
*  Debian 8 (« *jessie* ») ;
*  ou Debian 9 (« *stretch* »).

## Installation

### Sur Debian

```sh
apt install software-properties-common lsb-release
add-apt-repository \
  "deb [trusted=yes] https://repositories.canopsis.net/pulp/deb/debian$(cat /etc/debian_version | cut -d'.' -f 1)-canopsis/ stable main"
apt update
apt install canopsis-core
```

Pour déployer une configuration **mono-instance :**
```sh
canoctl deploy
```

### Sur CentOS 7

```sh
yum install yum-utils epel-release
echo "[canopsis]
name = canopsis
baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis/
gpgcheck=0
enabled=1" > /etc/yum.repos.d/canopsis.repo
yum install canopsis-core
```

Pour déployer une configuration **mono-instance :**
```sh
canoctl deploy
```
