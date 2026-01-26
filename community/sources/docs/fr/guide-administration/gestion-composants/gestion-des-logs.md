# Gestion des fichiers journaux

## Docker Compose

La mise en place d'une politique de rétention des logs nécessite la présence du logiciel `logrotate`.

Une fois que `logrotate` est installé sur votre machine, créer le fichier `/etc/logrotate.d/docker-container` suivant :

```
/var/lib/docker/containers/*/*.log {
  rotate 7
  daily
  compress
  minsize 100M
  notifempty
  missingok
  delaycompress
  copytruncate
}
```

Pour vérifier la validité de la configuration logrotate ajoutée, lancez la commande :

```sh
logrotate -dv /etc/logrotate.d/docker-container
```

Si vous souhaitez forcer une exécution manuelle de cette rotation sur-le-champ, vous pouvez éventuellement lancer la commande :

```sh
logrotate -fv /etc/logrotate.d/docker-container
```

## RPM (el8/el9)

Les logs de Canopsis sont gérées par `journald`, la rotation des logs est donc gérée par journald et non par logrotate.

Pour voir les logs de canopsis, il faut désormais passer par 

```sh
journalctl -u 'canopsis*' -f
```

!!! important

       Si vous souhaitez avoir les logs d'un seul service, vous devez remplacer `canopsis*` par le nom d'un des services de Canopsis.  
       Pour voir la liste des services de canopsis :
       ```sh
       systemctl list-dependencies canopsis.service --type=service --no-pager | grep -E 'canopsis'
       ```

### Rotation des logs de MongoDB

MongoDB, la base de données utilisée par Canopsis produit également des fichiers journaux qu'il convient de limiter.

Ceci peut être réalisé grâce à cette commande :

```sh
cat > /etc/logrotate.d/canopsis-mongodb.conf << EOF
/var/log/mongodb/*.log {
       daily
       rotate 30
       copytruncate
       delaycompress
       compress
       notifempty
       missingok
}
EOF
```

### Rotation des logs de TimescaleDB

TimescaleDB est la base utilisé pour logger les actions, les tech metrics, etc.. Elle produit également des fichiers journaux qu'il convient de limiter.

Deux modes sont disponible, soit via Logrotate, soit via journald

=== "Logrotate"

       Configurer logrotate pour timescaledb :

       ```sh
       cat > /etc/logrotate.d/canopsis-timescaledb.conf << EOF
       /var/lib/pgsql/15/data/log/*.log {
              daily
              rotate 30
              copytruncate
              delaycompress
              compress
              notifempty
              missingok
       }
       EOF
       ```

=== "journald"

    Configurer timescaledb pour envoyer les logs dans journald :

    ```sh
    sed -i "s/^#\?logging_collector.*/logging_collector = off/" /var/lib/pgsql/15/data/postgresql.conf
    ```

    puis on redémarre le service :

    ```sh
    systemctl restart postgresql-15.service
    ```

### Rotation des logs de RabbitMQ
```sh
cat > /etc/logrotate.d/canopsis-rabbitmq.conf << EOF
/var/log/rabbitmq/*.log {
       daily
       rotate 30
       copytruncate
       delaycompress
       compress
       notifempty
       missingok
}
EOF
```

### Rotation des logs de Valkey

```sh
cat > /etc/logrotate.d/canopsis-redis.conf << EOF
/var/log/redis/*.log {
       daily
       rotate 30
       copytruncate
       delaycompress
       compress
       notifempty
       missingok
}
EOF
```