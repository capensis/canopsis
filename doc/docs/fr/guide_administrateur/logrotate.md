# Rotation des logs Canopsis

## RabbitMQ, MongoDB, InfluxDB, Redis

Ces dépendances sont généralement installées par le biais de composants systèmes qui disposent déjà d'une rotation des logs, réalisée par `logrotate(8)`.

## gunicorn (serveur web intégré de Canopsis)

Il s'agit des logs `/opt/canopsis/var/log/webserver.log` et `/opt/canopsis/var/log/webserver-access.log`, peuplés lors d'un accès à l'interface web de Canopsis.

**Pré-requis :** Gunicorn doit être suffisamment récent pour disposer d'une rotation des logs. Les versions de Gunicorn ≥ 19.7.1, installées depuis Canopsis 2.5, sont capables de réaliser cette opération.

La rotation des logs de Gunicorn peut être réalisée à l'aide d'un fichier de ce type :

```
$ cat /etc/logrotate.d/gunicorn
/opt/canopsis/var/log/webserver*.log {
	daily
	rotate 7
	missingok
	notifempty
	compress
	delaycompress
	sharedscripts
	postrotate
		pkill -USR1 gunicorn
	endscript
	create 644 canopsis canopsis
}
```

Il s'agit ici d'une rétention quotidienne des logs du serveur web de Canopsis, sur 7 jours glissants. Voir `logrotate.conf(5)` pour d'autres possibilités de configuration.

Le lancement de cette rotation peut être forcé à l'aide de la commande suivante :
```bash
logrotate -vf /etc/logrotate.d/gunicorn
```

## Moteurs Canopsis

Il s'agit des logs contenus dans `/opt/canopsis/var/log/engines/`.

Les moteurs de Canopsis (qu'ils soient en Python ou en Go) ne disposent, pour l'instant, pas de rotation de logs. Leur taille reste cependant faible, de manière générale. 

Voir le [ticket #787](https://git.canopsis.net/canopsis/canopsis/issues/787) à ce sujet.
