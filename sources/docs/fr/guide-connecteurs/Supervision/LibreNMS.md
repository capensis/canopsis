# Connecteur LibreNMS

## Fonctionnement général

Il est important de noter que ce connecteur n'envoie que des alertes.

Voici le processus complet d'une vérification LibreNMS menant à une alerte :

*  Le poller LibreNMS recueille les données (principalement SNMP) des hôtes.
*  Ces données sont vérifiées par rapport à un ensemble de règles, configurables via l'interface Web ou dans un fichier de configuration.
*  Si une règle correspond, une alerte est ajoutée et prête à être envoyée.
*  LibreNMS envoie, via le script `alerts.php` des alertes sur chaque transport activé configuré.

## Installation et configuration

Installer `php-pecl-amqp` et l'activer dans la configuration PHP.

Copier le fichier connecteur dans `/opt/librenms/includes/alerts/` et configurer `/opt/librenms/config.php`:

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

## Environnement de test

1.  Désactiver alert cron. Dans `/etc/cron.d/librenms`, commenter l'appel de la ligne `/opt/librenms/alerts.php`.
2.  Créer / configurer une machine pour qu'elle dispose des interfaces et ports requis, afin de pouvoir changer leur état.
3.  Vous pouvez éventuellement installer Canopsis ou configurer `debug_noconnect` sur `true`. Ceci ignorera la connexion AMQP et imprimera des alertes sur la sortie standard du scripts `alerts.php`.
4.  Créer vos devices et règles.
5.  Lancer `poller.php -h <hostname>`

Lancer `alerts.php`.
