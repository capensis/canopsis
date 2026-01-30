# Gestion des fichiers journaux

## Consultation des logs

=== "RPM (EL8/EL9)"

       L'écriture des logs de Canopsis est gérée par `journald`, la rotation des logs est donc gérée nativement sans devoir passer par logrotate.
!!!Warning
    Attention néanmoins à activer la persistance de journald pour garder les logs après un reboot du serveur ( cf : https://docs.redhat.com/fr/documentation/red_hat_enterprise_linux/7/html/system_administrators_guide/s1-using_the_journal#s2-Enabling_Persistent_Storage )

       
    * Pour visualiser les logs de Canopsis, il faut utiliser la commande suivante

       ```sh
       journalctl -u 'canopsis*' -f
       ```

       Pour récupérer les logs d'un seul service, remplacer `canopsis*` par le nom d'un des services de Canopsis.  
       Pour voir la liste des services de Canopsis :
       
       ```sh
       systemctl list-dependencies canopsis.service --type=service --no-pager | grep -E 'canopsis'
       ```
       
=== "Docker Compose"

       Pour voir les logs de Canopsis dans Docker Compose, il faut utiliser les commandes suivantes.

       * Voir la liste des services de Canopsis : 
       ```sh
       docker compose config --services
       ```

       * Récupérer le nom du service et consulter les logs : 
       ```sh
       docker compose logs [service]
       ```

       * Pour afficher les logs en temps réel :
       ```sh
       docker compose logs -f [service]
       ```

=== "Helm"

       * Pour voir les logs de Canopsis dans Kubernetes utiliser la commande `kubectl`

       * Voir la liste des pods actifs de Canopsis : 
       ```sh
       kubectl get pods -n canopsis
       ```

       * Récupérer le nom du pod et consulter les logs : 
       ```sh
       kubectl logs -n canopsis [nom du pod]
       ```

       * Pour voir les logs en temps réel :
       ```sh
       kubectl logs -n canopsis -f [nom du pod]
       ```

## Extraire les logs

=== "RPM (EL8/EL9)"

       * Extraire les logs concernant des erreurs : 
       ```sh
       journalctl -u 'canopsis*' -f | grep -i "err" > canopsis-journald-$(date +%F_%H-%M).log
       ```

       * Récupérer les logs à partir d'une date :
       ```sh
       journalctl -u 'canopsis*' --since "2026-01-28 10:00" -f > canopsis-journald-$(date +%F_%H-%M).log
       ```
       
=== "Docker Compose"

       * Extraire les logs concernant des erreurs :
       ```sh
       docker compose logs [service] | grep -i "err" > canopsis-dockercompose-$(date +%F_%H-%M).log
       ```

       * Récupérer les logs à partir d'une date :
       ```sh
       docker compose logs --since $(date -d "2026-01-28 10:00" +%s) [service] > canopsis-dockercompose-$(date +%F_%H-%M).log
       ```
       *Docker utilise les timestamp au format UNIX, vous pouvez vous aider de sites comme [epochconverter.com](https://www.epochconverter.com/) pour convertir les dates en timestamp UNIX.*

=== "Helm"

       * Extraire les logs concernant des erreurs :
       ```sh
       kubectl logs -n canopsis [nom du pod] | grep -i err > canopsis-helm-$(date +%F_%H-%M).log
       ```

       * Récupérer les logs à partir d'une date :
       ```sh
       kubectl logs -n canopsis [nom-du-pod] --since-time="2026-01-28T10:00:00Z" > canopsis-helm-$(date +%F_%H-%M).log
       ```

## Rotation des logs

### Docker Compose

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
### RPM (EL8, EL9)

#### Rotation des logs de MongoDB

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

#### Rotation des logs de TimescaleDB

TimescaleDB est la base utilisé pour logger les actions, les tech metrics, etc.. Elle produit également des fichiers journaux qu'il convient de limiter.

Deux modes sont disponible, soit via Logrotate, soit via journald

=== "Logrotate"

       Configurer logrotate pour timescaledb :

       ```sh
       cat > /etc/logrotate.d/canopsis-timescaledb.conf << EOF
       /var/lib/pgsql/17/data/log/*.log {
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
    sed -i "s/^#\?logging_collector.*/logging_collector = off/" /var/lib/pgsql/17/data/postgresql.conf
    ```

    puis on redémarre le service :

    ```sh
    systemctl restart postgresql-17.service
    ```

#### Rotation des logs de RabbitMQ

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

#### Rotation des logs de Valkey

```sh
cat > /etc/logrotate.d/canopsis-valkey.conf << EOF
/var/log/valkey/*.log {
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