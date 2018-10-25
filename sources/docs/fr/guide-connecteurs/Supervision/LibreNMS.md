# Guide Connecteur 

## Section : Supervision / LibreNMS, rules, alerts and transport

# Installation & Configuration

Instaler `php-pecl-amqp` et l'activer dans la configuration PHP.

Copiez le fichier connecteur dans `/opt/librenms/includes/alerts/` et configurez `/opt/librenms/config.php`:

```php
$config['alert']['transports']['canopsis'] = Array();
$config['alert']['transports']['canopsis']['user'] = 'cpsrabbit';
$config['alert']['transports']['canopsis']['password'] = 'canopsis';
$config['alert']['transports']['canopsis']['vhost'] = 'canopsis';
$config['alert']['transports']['canopsis']['host'] = '10.25.190.159';
$config['alert']['transports']['canopsis']['port'] = 5672;
$config['alert']['transports']['canopsis']['exchange_name'] = 'canopsis.events';
$config['alert']['transports']['canopsis']['connector_name'] = 'LibreNMS';
$config['alert']['transports']['canopsis']['debug_noconnect'] = false;
```


## Comment ça marche ?

Tout d'abord, gardez à l'esprit qu'un connecteur d'alerte librenms n'envoie que des alertes.

Voici le processus complet d'une vérification LibreNMS menant à une alerte:


Le poller LibreNMS recueillera les données, principalement SNMP, de vos hôtes.
Ces données seront vérifiées par rapport à un ensemble de règles, configurables via l'interface Web ou dans un fichier de configuration.
Si une règle correspond, une alerte sera ajoutée et prête à être envoyée.
LibreNMS, via le script `alerts.php`, enverra des alertes sur chaque transport activé configuré.


## L'envoi d'alerte fonctionne comme ceci:


LibreNMS parcourt toutes les alertes
LibreNMS effectuera une boucle sur chaque transport activé pour chaque alerte: `$config['alert']['transports']['transport_name'] != empty`

Le fichier PHP de transport complet est lu, stocké dans une chaîne déclarant une fonction temporaire, puis évalué (!).
Cette fonction est exécutée avec certains paramètres donnés: `$obj` pour l'alerte, `$opts` pour la configuration du transport.


Veillez à éviter les longs délais d'attente lors de l'envoi d'alertes, car chaque traitement d'alerte est bloquant.


Environnement de test

1. Désactiver alert cron:, `/etc/cron.d/librenms` commentez l'appel de la ligne `/opt/librenms/alerts.php`.
2. Créez / configurez une machine (virtuelle ou non) pour qu'elle dispose des interfaces / ports / requis, de sorte que vous puissiez changer leur état.
3. Vous pouvez éventuellement installer Canopsis ou le configurer `debug_noconnectsur` sur `true`. Cela ignorera la connexion AMQP et imprimera des alertes sur la sortie standard du `alerts.php` script.
4. Créez vos devices, règles
5. Run `poller.php -h <test_host_name>`

Run `alerts.php`