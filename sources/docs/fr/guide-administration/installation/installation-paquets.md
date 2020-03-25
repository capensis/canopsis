# Installation de Canopsis avec un paquet Debian ou CentOS

Canopsis peut être installé à l'aide de paquets sur Debian 9 (« *stretch* ») et CentOS 7. Les binaires sont compilés pour l'architecture x86-64.

Cette procédure décrit l'installation de l'édition open-source de Canopsis en mono-instance.

L'ensemble des procédures décrites doivent être réalisées avec l'utilisateur `root`.

## Prérequis

Assurez vous d'avoir suivi les [prérequis réseau et de sécurité](pre-requis-parefeu-et-selinux.md).

!!! attention
    Notez que Canopsis installe ses propres versions d'InfluxDB, MongoDB, Redis et RabbitMQ, et que seules ces versions sont validées pour fonctionner avec Canopsis. Veillez à ne pas remplacer ces versions par vos propres versions, que ce soit de façon intentionnelle, ou par l'ajout de dépôt tiers qui pourraient écraser les versions installées avec Canopsis (ex : installation des dépôts officiels InfluxDB pour l'ajout d'un Telegraf).

## Étape 1 : installation des paquets

### Sur Debian 9

Il peut être nécessaire d'activer les dépôts `non-free` de Debian, dans le cas d'une installation CAT.

Application des dernières mises à jour du système :
```sh
apt update
apt upgrade
```

Ajout du dépôt Canopsis (qui permettra aussi de récupérer les mises à jour) :
```sh
apt install apt-transport-https
echo "deb [trusted=yes] https://repositories.canopsis.net/pulp/deb/debian9-canopsis/ stable main" > /etc/apt/sources.list.d/canopsis.list
```

Installation de l'édition open-source de Canopsis :
```sh
apt update
apt install canopsis-core
```

### Sur CentOS 7

Activation d'EPEL et application des dernières mises à jour du système :
```sh
yum install yum-utils epel-release
yum update
```

Installation d'une version plus récente de Python 2.7 (pré-requis pour le `webserver` Canopsis, depuis Canopsis 3.13.0), [depuis SCL](https://www.softwarecollections.org/en/scls/rhscl/python27/) :
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

## Étape 2 : initialisation de l'environnement

Une fois le paquet installé, vous devez lancer une commande d'initialisation de votre environnement. Cette commande s'occupe notamment d'installer les briques et réglages nécessaires à Canopsis, de finaliser les fichiers de configuration et d'activer les bons moteurs.

Cette commande n'est destinée à être lancée qu'**une seule fois après une installation** de Canopsis. Pour les mises à jour, suivez la procédure décrite dans le [Guide de mise à jour](../mise-a-jour/index.md).

### Versions de Canopsis 3.30.0 et inférieures

Dans les versions de Canopsis 3.30.0 et inférieures, la commande à utiliser est la suivante. Elle déploie une installation mono-nœud avec les [moteurs Python « historiques »](../moteurs/index.md#moteurs-python) :

```sh
canoctl deploy
```

Vous pouvez alors procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).

### Versions de Canopsis 3.31.0 et supérieures

À partir de [Canopsis 3.31.0](../../notes-de-version/3.31.0.md), vous pouvez choisir de déployer facilement votre environnement Canopsis mono-nœud avec les [moteurs Python « historiques »](../moteurs/index.md#moteurs-python) ou avec les [moteurs Go « nouvelle génération »](../moteurs/index.md#moteurs-go).

L'environnement Python reste disponible pour l'instant, mais les efforts de développement se concentrent dorénavant sur les moteurs Go, qui apportent notamment de meilleures performances et de nouvelles fonctionnalités.

Si vous souhaitez réaliser une installation de Canopsis en environnement Python, exécutez la commande suivante :

```sh
canoctl deploy-python
```

Si vous souhaitez plutôt bénéficier des moteurs Go nouvelle génération (recommandé), exécutez les commandes suivantes :

```sh
# Sur Debian :
apt install canopsis-engines-go
# Sur CentOS :
yum install canopsis-engines-go

# Puis, dans tous les cas :
canoctl deploy-go
```

Vous pouvez alors procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).
