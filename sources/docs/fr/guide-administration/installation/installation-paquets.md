# Installation de Canopsis avec un paquet Debian ou CentOS

Cette procédure décrit l'installation de l'édition open-source de Canopsis en mono-instance.

L'ensemble des procédures décrites doivent être réalisées avec l'utilisateur `root`.

## Pré-requis

Canopsis peut être installé à l'aide de paquets sur Debian 9 (« *stretch* ») ou sur CentOS 7.

!!! attention
    Les prochaines versions de Canopsis ne seront bientôt plus compatibles avec **Debian 8** (« *jessie* »). Pensez à mettre à jour votre environnement vers Debian 9 pour continuer de profiter des mises à jour de l'outil.

## Installation des paquets

### Sur Debian 9

**Attention :** Sur Debian 9, le paquet Canopsis Cat à partir de la version 3.24.0 nécessite d'activer les dépôts `non-free` de Debian.

Application des dernières mises à jour de votre système :
```sh
apt update
apt upgrade
```

Ajout du dépôt Canopsis (qui permettra aussi de récupérer les mises à jour) :
```sh
apt install apt-transport-https lsb-release
echo "deb [trusted=yes] https://repositories.canopsis.net/pulp/deb/debian$(cat /etc/debian_version | cut -d'.' -f 1)-canopsis/ stable main" \
> /etc/apt/sources.list.d/canopsis.list
```

Installation de l'édition open-source de Canopsis :
```sh
apt update
apt install canopsis-core
```

### Sur CentOS 7

!!! note
    Les versions de CentOS inférieures à CentOS 7 ne sont **pas** prises en charge.

Activation d'EPEL et application des dernières mises à jour de votre système :
```sh
yum install yum-utils epel-release
yum update
```

Installation d'une version plus récente de Python 2.7 (pré-requis pour le `webserver` Canopsis en environnement CentOS 7, depuis Canopsis 3.13.0), [depuis SCL](https://www.softwarecollections.org/en/scls/rhscl/python27/) :
```sh
yum install centos-release-scl
yum install python27
```

Ajout du dépôt Canopsis (qui permettra aussi de récupérer les mises à jour) :
```sh
echo "[canopsis]
name = canopsis
baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis/
gpgcheck=0
enabled=1" > /etc/yum.repos.d/canopsis.repo
```

Installation de l'édition open-source de Canopsis :
```sh
yum install canopsis-core
```

## Mise en service

!!! attention
    La commande suivante ne doit être exécutée que lors de votre **première** installation de Canopsis, sans quoi certains éléments de configuration seront totalement réinitialisés.

    Pour procéder à une mise jour de Canopsis, voir [la documentation de mise à jour](../../mise-a-jour/).

Une fois le paquet installé, vous pouvez déployer une configuration **mono-instance** à l'aide de la commande suivante :
```sh
canoctl deploy
```
