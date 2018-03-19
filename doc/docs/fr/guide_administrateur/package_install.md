## Installation

### CentOS / RedHat 7

```bash
yum localinstall canopsis-cat-<version>-1.el7.centos.x86_64.rpm
```

### Debian 8 / 9

```bash
apt-get update
dpkg -i canopsis-cat-1-<version>.amd64.<platform>.deb
apt install -f
```

## Configuration

Les fichiers suivants doivent être modifiés :

 * `/opt/canopsis/etc/common/mongo_store.conf`
 * `/opt/canopsis/etc/infux/storage.conf`
 * `/opt/canopsis/etc/amqp.conf`

Ou a défaut vous pouvez insérer dans votre `/etc/hosts` si c’est une installation locale :

```
127.0.0.1 localhost mongodb influxdb rabbitmq
```

## Intégrations externes

 * MongoDB
 * RabbitMQ
 * InfluxDB

Cette partie de la documentation a été réalisée pour CentOS 7. Adapter pour les autres distributions.

### MongoDB

Se référer à la documentation d’installation des paquets `mongodb-org` en version `3.4` : https://docs.mongodb.com/v3.4/administration/install-on-linux/

Vour trouverez les fichiers `create_admin.js` et `create_user.js` dans `doc/docs/files/mongo`.

Désactiver l’authentification MongoDB puis :

```bash
systemctl start mongod
mongo load create_admin.js
```

Activer l’authentification MongoDB `security.authorization: enabled` puis :

```
systemctl restart mongod
mongo load create_user.js
```

### RabbitMQ

Installer RabbitMQ puis :

```bash
systemctl start rabbitmq-server

rabbitmqctl add_user admin admin
rabbitmqctl set_user_tags admin administrator
rabbitmqctl add_vhost canopsis
rabbitmqctl set_permissions -p canopsis admin ".*" ".*" ".*"
rabbitmqctl add_user cpsrabbit canopsis
rabbitmqctl set_user_tags cpsrabbit canopsis
rabbitmqctl set_permissions -p canopsis cpsrabbit ".*" ".*" ".*"
rabbitmq-plugins enable rabbitmq_management

systemctl restart rabbitmq-server
```

## InfluxDB

Installer InfluxDB 0.10.x

```bash
systemctl start influxdb

influx --execute "CREATE USER admin WITH PASSWORD 'admin' WITH ALL PRIVILEGES"
influx --execute "CREATE USER cpsinflux WITH PASSWORD 'canopsis'"
influx --execute "GRANT ALL ON canopsis TO cpsinflux"

sed -i "s/auth-enabled = false/auth-enabled = true/" /etc/influxdb/influxdb.conf

systemctl restart influxdb
```

## Init

Des unités `systemd` sont disponibles :

 * `canopsis-engine@<module>-<name>.service`
 * `canopsis-webserver.service`

Voici tous les engines qui vous pouvez activer dans `core` et `cat` :

```bash
systemctl enable canopsis-engine@dynamic-alerts.service
systemctl enable canopsis-engine@cleaner-cleaner_alerts.service
systemctl enable canopsis-engine@cleaner-cleaner_events.service
systemctl enable canopsis-engine@dynamic-context-graph.service
systemctl enable canopsis-engine@event_filter-event_filter.service
systemctl enable canopsis-engine@eventstore-eventstore.service
systemctl enable canopsis-engine@linklist-linklist.service
systemctl enable canopsis-engine@dynamic-pbehavior.service
systemctl enable canopsis-engine@dynamic-perfdata.service
systemctl enable canopsis-engine@scheduler-scheduler.service
systemctl enable canopsis-engine@selector-selector.service
systemctl enable canopsis-engine@task_dataclean-task_dataclean.service
systemctl enable canopsis-engine@task_importctx-task_importctx.service
systemctl enable canopsis-engine@task_linklist-task_linklist.service
systemctl enable canopsis-engine@task_mail-task_mail.service
systemctl enable canopsis-engine@ticket-ticket.service
systemctl enable canopsis-engine@dynamic-watcher.service

systemctl enable canopsis-webserver.service
```

Quelques exemples de gestion des services avec systemd :

```bash
# Démarrer tout canopsis
/opt/canopsis/bin/canopsis-systemd start
# Récupérer les status détaillés
/opt/canopsis/bin/canopsis-systemd status

# Désactiver un engine cassé et le supprimé des failures
systemctl disable canopsis-engine@badinstance-badinstance
systemctl reset-failed

# Redémarrer tout canopsis
/opt/canopsis/bin/canopsis-systemd restart

# Arrêter tout canopsis
/opt/canopsis/bin/canopsis-systemd stop

# Lister toutes les unités canopsis avec un affichage compacte
systemctl list-units -a "canopsis*"

# Redémarrer le webserver
systemctl restart canopsis-webserver

# Redémarrer un engine
systemctl restart canopsis-engine@cooleng-cooleng
```

Le fichier `/opt/canopsis/etc/amqp2engines.conf` est toujours en vigeur.

### Nombre de process

Pour le moment le nombre de processus lancés via `engine-launcher` est fixé dans les unités.

Pour changer le nombre d’instances :

```bash
mkdir -p /etc/systemd/system/canopsis-engine@<module>-<name>
cat > /etc/systemd/system/canopsis-engine@<module>-<name>/workers.conf << EOF
[Service]
Environment=WORKERS=X
EOF
```

Remplacer `X` par le nombre de workers désiré. Par défaut `1`.

```bash
systemctl daemon-reload
systemctl restart canopsis-engine@<module>-<name>.service
```
