# WebServer

## Gunicorn

Le lancement de `gunicorn` n’utilise pas de fichier de configuration.

La rotation des logs avec `gunicorn` se fait en deux étapes :

 * `mv` du fichier de log `var/log/webserver-access.log`
 * `kill -USR1 <pid master>` où `<pid master>` doit être le PID du processus maître de `gunicorn`.

## uWSGI

`uWSGI` est lancé avec un fichier de configuration `etc/uwsgi-webserver.ini`.

La rotation des logs avec `uWSGI` se fait en deux étapes :

 * `mv` du fichier de log `var/log/webserver-access.log`
 * `kill -HUP <pid master>` où `<pid master>` doit être le PID du processus maître de `uWSGI`.
