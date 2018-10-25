# Guide Administrateur

**TODO (DWU) :** regénéraliser cette doc (qui était prévue dans le cas particulier d'un fulldisk).

# Section : Administration avancée / Backup & Restore

# Gestion d'un cas de disque plein sur Mongo 2

## Description

On prend l'exemple d'une infrastructure Mongo 2 où `/var` est totalement plein (ou bientôt totalement plein) au point de bloquer tous les services.

L'idée est que Mongo est généralement responsable de la situation : on peut alors tenter de faire une sauvegarde puis une restauration de sa base, afin qu'il nettoie l'espace disque inutilement alloué par d'anciennes données. Mongo a cependant besoin de quelques Go libres sur la partition afin de pouvoir effectuer cette opération : on va donc effectuer quelques manipulations pour retrouver un peu de place.

**Important :** les manipulations suivantes sont à moduler en fonction de l'espace disque nécessaire à la remise en place du service. Il n'est pas forcément nécessaire d'exécuter la totalité des procédures décrites dans ce document.

**Pré-requis :** on part du principe que l'on dispose d'un espace de stockage de dimension suffisante pour pouvoir y déplacer temporairement des fichiers. On prend ici l'exemple d'un point de montage `/data` accessible en NFS.

## Libération d'espace disque

### Trouver les services ayant le plus rempli la partition

Valider que `/var` est bien la partition pleine :
```
# df -h
...
/dev/mapper/vg_data-lv_var                   116G   116G   0G  100% /var
...
```

Y chercher les fichiers ou dossiers occupant le plus d'espace :
```
# du -hs /var/log/* /var/lib/* | grep G
 3G     /var/log/syslog
93G     /var/lib/mongodb
16G     /var/lib/rabbitmq
```

Il est possible que RabbitMQ ait empiré la situation en flushant continuellement sur le disque les données qu'il n'arrive plus à envoyer à MongoDB, comme on peut le voir dans le résultat de la commande précédente.

### Arrêt des services

Les services ont besoin d'être arrêtés afin d'arrêter les écritures en cours et de pouvoir libérer de l'espace pour Mongo.

```
# systemctl stop influxdb.service
# systemctl stop rabbitmq-server.service
```

Il se peut que certains services soient totalement bloqués à cause de la partition pleine : le service n'arrivera alors même plus à s'arrêter. Juger au cas par cas s'il est nécessaire d'employer la manière forte (`kill -9`...) pour terminer le processus.

On arrête aussi les services internes à Canopsis :
```
# su - canopsis
(canopsis) service amqp2engines* mstop
(canopsis) service webserver stop
(canopsis) exit
```

Puis, vérifier avec un `top` ou un `ps` qu'il n'y a pas d'autre service pouvant remplir la partition ou bloquer le reste de la manipulation :
```
# top
```

### Pré-requis : libération de quelques Go avant la restauration MongoDB

Quelques Go doivent être disponibles dans `/var` afin de pouvoir lancer la restauration de MongoDB. À titre d'exemple, 3 Go libres *peuvent* être suffisants.

Les étapes suivantes ne sont donc nécessaires que jusqu'au moment où l'on a libéré suffisamment d'espace. Il faut donc  relancer les commandes `df`/`du` précédentes jusqu'à avoir libéré les quelques Go dont MongoDB a besoin pour lancer une restauration.

*Si nécessaire*, déplacer temporairement les fichiers InfluxDB :
```
# systemctl status influxdb.service
[S'assurer que le service est bien arrêté avant de lancer la prochaine commande]
# mv /var/lib/influx /data/
```

*Si nécessaire*, récupérer quelques Go en vidant quelques fichiers de logs :
```
# :> /var/log/daemon.log
# :> /var/log/syslog
```

*Si nécessaire*, désactiver temporairement le journal Mongo :
```
# vi /etc/mongod.conf
[Passer la valeur "journal = true" à "journal = false"]
# systemctl restart mongod.service
```

### Sauvegarde de la base MongoDB courante

La commande suivante fait une sauvegarde de la base de données courante, et peut prendre quelques heures (penser à un `tmux` pour reprendre la main plus tard si nécessaire) :
```
[Vers un point de montage disposant de suffisamment d'espace pour une sauvegarde de la base]
# MONGO_DATE=$(date '+%Y%m%d') ; mkdir -p /data/backup-mongo$MONGO_DATE
# mongodump -u cpsmongo -p canopsis -d canopsis -o /data/backup-mongo$MONGO_DATE/
```

**Temps d'exécution :** à titre d'exemple, il a fallu 4h30 pour faire la sauvegarde d'une base de données Mongo de ce type vers le point de montage NFS :
```js
> db.stats()
{
        "db" : "canopsis",
        "collections" : 35,
        "objects" : 20 827 775,
        "avgObjSize" : 1903.0303928288067,
        "dataSize" : 39635888840,
        "storageSize" : 47567834864,
        "numExtents" : 153,
        "indexes" : 125,
        "indexSize" : 26848643136,
        "fileSize" : 92230647808,
        "nsSizeMB" : 16,
        "dataFileVersion" : {
                "major" : 4,
                "minor" : 5
        },
        "extentFreeList" : {
                "num" : 0,
                "totalSize" : 0
        },
        "ok" : 1
}
```

### Restauration de la sauvegarde MongoDB

Après avoir vérifié que la commande `mongodump` précédente a fonctionné, on peut effacer puis réimporter la base, afin que MongoDB puisse purger l'espace de stockage qui était encore réservé par les données supprimées.

Suppression de la base courante :
```
# mongo admin -u admin -p admin
> use canopsis
[s'assurer de bien avoir une sauvegarde OK avant d'exécuter la commande suivante !]
> db.dropDatabase()
> exit
```

Restauration de la sauvegarde. Cette étape peut, à nouveau, prendre quelques heures :
```
# mongorestore -u admin -p admin /data/backup-mongo$MONGO_DATE/
```

**Temps d'exécution :** à titre d'exemple, il a fallu 2h30 pour restaurer la sauvegarde décrite plus haut (la très grande majorité du temps étant sur les index).

## Remise en place des services

### Éventuelle remise en place du journal MongoDB

*Si on a dû désactiver le journal MongoDB précédemment*, on doit maintenant le remettre en place :
```
# vi /etc/mongod.conf
[Passer la valeur "journal = false" à "journal = true"]
# systemctl restart mongod.service
```

### Éventuelle restauration d'InfluxDB

*Si l'on a dû déplacer* `/var/lib/influxdb` précédemment pour libérer de l'espace, il faut maintenant le restaurer :
```
# cp -pR /data/influx /var/lib/
# systemctl start influxdb.service
```

### Éventuel nettoyage de RabbitMQ

**Si et seulement si** RabbitMQ a été totalement planté par la partition pleine au point de ne plus pouvoir redémarrer (c'est-à-dire au point que `systemctl restart rabbitmq-server.service` n'aboutit jamais), on le réinstalle afin de repartir sur une base propre.

**Attention**, on perdra alors les données RabbitMQ existantes :
```
# apt-get remove --purge rabbitmq-server
# rm -rf /var/lib/rabbitmq
# apt-get update
# apt-get install rabbitmq-server
```

Mise en place des accès Canopsis par défaut :
```
# rabbitmqctl add_user cpsrabbit canopsis
# rabbitmqctl set_user_tags cpsrabbit administrator
# rabbitmqctl add_vhost canopsis
# rabbitmqctl set_permissions -p canopsis cpsrabbit ".*" ".*" ".*"
```

Activation de la console web de gestion de RabbitMQ :
```
# rabbitmq-plugins enable rabbitmq_management
# rabbitmq-plugins list
[vérifier que le module "rabbitmq_management" est bien activé et que l'interface web fonctionne]
```

### Redémarrage des moteurs Canopsis

Redémarrer rapidement les moteurs internes à Canopsis :
```
# su - canopsis
(canopsis) service amqp2engines* mstart
(canopsis) service webserver start
(canopsis) exit
```

### Tests

Il faut maintenant s'assurer que l'ensemble des services et l'interface Canopsis fonctionnent correctement.
