# Présentation

Vous trouverez ici toute la documentation nécessaire au développement sur canopsis.

!!! tip "Note"
    Cette page contient le plan de la documentation de développement, qui est en cours d'écriture.

# Process de développement
## Organisation des dépôts
## Process de release
## Nomenclature des messages de commit
<!--  - specification des segments de canopsis (alerts, action, …) -->

# Installation d'un environnement de développement

## Installation des solutions de virtualisation/contenerisation
### LXC

LXC est une solution de contenerisation **système**, à la différence de Docker qui est orienté contenerisation **applicative**.

La différence réside dans la capacité de LXC à démarrer un système quasi complet comme le ferait une machine virtuelle (VM), là où Docker ne permet bien souvent que de démarrer une seule application/programme et n’arrive pas à isoler convenablement les logiciels tels que l’init sans lui octroyer des droits dangereux pour un environnement de production.

L’usage de LXC au sein de Canopsis est dédié à :

 * Permettre le développement rapide sans avoir à reconstruire un paquet complet à chaque fois
 * Être proche d’une VM en terme d’usage sans les problématiques de consommation de ressources
 * Très bonne intégration dans Ubuntu et dossiers partagés sans problèmes de permissions contrairement à VBox

Cependant l’installation et la mise sur pied d’un conteneur LXC est un peu sportive. Vous verrez cependant que vous pouvez vous épargner une grande partie de ces manipulations en conservant une copie préconfigurée d’un conteneur.

#### Installation des outils

Pour la durée de l'installation, restez en superutilisateur grâce à cette commande :  
```bash
sudo -s
```

Une fois cela fait, continuez en superutilisateur jusqu'à la fin de l'installation.

```bash
apt install -y lxc
```

Sur Ubuntu, les services nécessaires pour obtenir des conteneurs LXC fonctionnels et connectés au réseau sont lancés, il n’y a rien à faire de plus.

#### Création d’un conteneur LXC

Pour créer un conteneur LXC nous pouvons :

 * Télécharger via des scripts que nous allons voir une image minimale
 * Décompresser une archive à la main

##### Création assistée

```bash
export DOWNLOAD_KEYSERVER="hkp://p80.pool.sks-keyservers.net:80"

# Debian 8 (Jessie)
lxc-create -t download -n <nom_conteneur> -- --dist debian --release jessie --arch amd64

# Debian 9 (Stretch)
lxc-create -t download -n <nom_conteneur> -- --dist debian --release stretch --arch amd64

# CentOS 7
lxc-create -t download -n <nom_conteneur> -- --dist centos --release 7 --arch amd64
```

Le template `download` peut être lancé sans arguments pour obtenir un mode de sélection interactif :

```bash
export DOWNLOAD_KEYSERVER="hkp://p80.pool.sks-keyservers.net:80"
lxc-create -t download -n <nom_conteneur>
```

##### Sauvegarde et récupération d’un conteneur existant

Les conteneurs sont stockés dans `/var/lib/lxc/<nom_conteneur>`.

Il est utile de récupérer un conteneur lorsque par exemple on a configuré un serveur SSH, déployé sa clef ou préinstallé quelques outils, afin d’éviter de devoir le refaire à chaque fois que l’on veut repartir de zéro.

**Attention :** lorsqu’on souhaite refaire une installation de Canopsis, il n’est pas souhaitable d’embarquer un Canopsis déjà installé dans un conteneur "de récupération".

Tout ce que vous avez à faire sont ces étapes :

**Sauvegarde**

```bash
lxc-stop -n <nom_conteneur>
cd /var/lib/lxc
tar cf <nom_conteneur>.tar <nom_conteneur>
```

**Récupération**

```bash
lxc-stop -k -n <nom_conteneur>
cd /var/lib/lxc
rm -rf <nom_conteneur>
tar xpf <nom_conteneur>.tar
```

#### Configuration et utilisation du conteneur

Pour pouvoir faciliter le travail avec le conteneur nous allons configurer quelques points de montage entre le système hôte et le conteneur.

En admettant que toutes les sources de Canopsis (moteurs Go, python, webui…) soient dans un seul dossier, par exemple : `${HOME}/canopsis` :

```bash
cd /var/lib/lxc/<nom_conteneur>

echo "lxc.mount.entry = ${HOME}/canopsis mnt/canopsis none rw,bind 0 0" >> config

mkdir rootfs/mnt/canopsis
```

**Démarrer le conteneur :**

```bash
lxc-start <nom_conteneur>
```

Quand nous nous connecterons au conteneur, les sources seront disponibles dans `/mnt/canopsis`.

**Ne séparez pas vos sources** dans différents points de montage : c’est inutile.

**Ne supprimez pas votre conteneur** sans l’avoir arrêté préalablement : le point de montage est actif tant que le conteneur est en cours d’exécution.

Pour savoir quels conteneurs sont actuellement en exécution :

```bash
lxc-ls --active
```

**SSH**

Bien qu’il soit possible d’utiliser le conteneur sans s’y connecter en SSH, **CELA N’EST ABSOLUMENT PAS VIABLE ET VOUS POSERA DES PROBLÈMES.**

Faites nous confiance.

```bash
# Changez le mot de passe root de votre conteneur : vous n’aurez pas besoin d’autre chose qu’un accès root
lxc-attach -n <nom_conteneur> --clear-env -- passwd root
```

Pour Debian :

```bash
lxc-attach -n <nom_conteneur> --clear-env -- apt install -yf openssh-server openssh-client
sed -i /var/lib/lxc/<nom_conteneur>/rootfs/etc/ssh/sshd_config -re 's/^[#]?PermitRootLogin .*/PermitRootLogin yes/g'
lxc-attach -n <nom_conteneur> --clear-env -- systemctl restart ssh
```

Pour CentOS :

```bash
lxc-attach -n <nom_conteneur> --clear-env -- yum install -y openssh
sed -i /var/lib/lxc/<nom_conteneur>/rootfs/etc/ssh/sshd_config -re 's/^[#]?PermitRootLogin .*/PermitRootLogin yes/g'
lxc-attach -n <nom_conteneur> --clear-env -- systemctl enable sshd
lxc-attach -n <nom_conteneur> --clear-env -- systemctl restart sshd
```

#### Se connecter au conteneur

Maintenant que tout est prêt, on peut se connecter au conteneur :

```bash
lxc-attach -n <nom_conteneur> -- ip a
```

On obtient une sortie de ce genre :

```bash
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
6: eth0@if7: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 00:16:3e:46:1a:6d brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.0.3.200/24 brd 10.0.3.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::216:3eff:fe46:1a6d/64 scope link
       valid_lft forever preferred_lft forever
```

L’IP qui nous intéresse est la `10.0.3.200`. Chez vous ce sera nécessairement une autre IP car c’est une réservation DHCP locale qui est établie.

**Vous n’avez plus besoin d’être en root pour vous connecter à votre conteneur, à partir d’ici toutes les commandes seront exécutée avec votre utilisateur "normal".**

```bash
ssh-keygen -t ecdsa
ssh-copy-id ~/.ssh/id_ecdsa.pub root@<ip_conteneur>
ssh root@<ip_conteneur>
```

### Docker
Docker est une solution de contenerisation **applicative**, à la différence de LXC qui est orienté contenerisation **système**.

N’essayez pas de démarrer plusieurs application et certainement pas un *init* dans Docker : vous vous exposez à des problèmes, et nous ne le ferons pas ici car l’usage de Docker par rapport à LXC est différent.

L’usage de Docker au sein de Canopsis permet de :

 * Facilement recetter une branche ou une *release* à venir
 * Fournir des images Docker aux clients et à la communauté
 * Générer les paquets RPM et DEB pour un déploiement chez les clients, pour la communauté mais aussi pour nous dans les conteneurs LXC afin de faciliter le développement

Il n’est **pas** possible de développer avec Docker : chaque processus est exécuté dans une **copie** d’une image : une modification faite sur un conteneur ne sera pas répercutée sur un autre.

#### Installation

Suivre les étapes de la documentation officielle :

https://docs.docker.com/install/linux/docker-ce/ubuntu/

#### Configuration

Aucune configuration n’est nécessaire.

## Installation des technologies

### Venv Python

Un environnement virtuel Python est utilisé dans tous les environnement Canopsis : RPM, DEB, Docker… tous embarquent en venv au complet.

Il est plutôt logique d’utiliser un venv pour vous afin de développer.

Nous allons admettre que les sources de canopsis se trouvent dans `${HOME}/canopsis`

#### virtualenv

Sur Ubuntu 18.xx, la version par défaut de Python est encore Python 2, qui est pour le moment la base du code de Canopsis.

```bash
sudo apt install -y python python-virtualenv virtualenv
virtualenv ~/doc/pyvenv/canopsis
```


### Installation de Go

Téléchargez la [dernière version de Go](https://golang.org/dl/) et suivez [les instructions](https://golang.org/doc/install) pour l'installer.  

Il est préférable de lire les instructions complètes, néanmoins voici un résumé très court des instructions.  

```bash
wget https://dl.google.com/go/go<version>.linux-amd64.tar.gz
rm -rf /usr/local/go && tar xf go<version>.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin
```

Définir l'environnement go :

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

mkdir -p $GOPATH/{bin,src}
mkdir -p $GOPATH/src/git.canopsis.net/canopsis/go-revolution
```


Une fois cela fait, installez [Glide](https://glide.sh/), disponible dans les dépôts Ubuntu, Debian et CentOS, sous le nom `golang-glide`.  

Glide, à l'instar de Pip et npm, est un _packet manager_ pour le langage Go et ses modules.  


### Utilitaires

#### Robot3t

Anciennement Robomongo, [Robo3t](https://robomongo.org/) est un logiciel d'administration pour mongoDB.  

Pour l'installer, voici quelques commandes qui devraient vous simplifier la vie :  

```bash
wget https://download.robomongo.org/1.2.1/linux/robo3t-1.2.1-linux-x86_64-3e50a65.tar.gz
sudo tar -xvzf robo3t-1.2.1-linux-x86_64-3e50a65.tar.gz -C /opt
rm robo3t-1.2.1-linux-x86_64-3e50a65.tar.gz
sudo chmod +x robo3t ./robo3t
sudo ln -s /opt/robo3t-1.2.1-linux-x86_64-3e50a65/bin/robo3t /usr/local/bin/robo3t
```


#### Postman

Pour tester les API, Postman est un excellent logiciel utilitaire permettant de construire puis d'envoyer des requêtes. C'est un utilitaire bien plus simple et pratique que `curl`, et permet aussi de partager facilement avec votre équipe les requêtes que vous utilisez.


Pour l'installer, voici quelques commandes qui devraient vous simplifier la vie :  

```bash
wget https://dl.pstmn.io/download/latest/linux64 -O postman.tar.gz
sudo tar -xzf postman.tar.gz -C /opt
rm postman.tar.gz
sudo ln -s /opt/Postman/Postman /usr/local/bin/postman
```

# Backend
## Python
### Installation de nouvelle source python
### Structure du projet
<!--
  - organisation des packages
  - architecture à mettre en place : modele, adapter, api
-->
### Création d'engines
### Création d'API

# Golang
## Compilation
## Architecture du projet
## Création de moteurs
## Implémentation de source de données externes (pour l'event-filter)

# Base de données
<!--
## default_entities
### Présentation générale
### Présentation de la structure d'un document.
## periodical_alarms
### Présentation générale
### Présentation de la structure d'un document.
-->

# Front-end
## Mise en place de l'environnement de développement
## Technologies utilisées
## Structure du projet
## Règles de style
## Les mixins, helpers et filters
## Le store Vuex
## Guides de création nouvelle fonctionnalité
### Modal
### Vue
### Widget (+ Paramètres du widget)

# API

[Présentation de toutes les routes disponibles](API.md)

  * Pbehavior
  * Event-filter
  * [Healthcheck](./healthcheck/api_v2_healthcheck.md)

# Gestion de la documentation
