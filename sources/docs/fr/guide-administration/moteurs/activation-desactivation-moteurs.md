#  Activation et désactivation des moteurs

Moteurs Go :
```sh
systemctl start/stop canopsis-engine-go\@\*
```

Moteurs Python
```sh
systemctl start/stop canopsis-engine@cleaner-cleaner_alerts.service
systemctl start/stop canopsis-engine@dynamic-alerts.service
systemctl start/stop canopsis-engine@cleaner-cleaner_events.service
```

Vérification de leur bon lancement :
```sh
systemctl status canopsis-engine-go\@\* -l
```

Arrêt d'un moteur en particulier :
```sh
systemctl stop canopsis-engine-go\@engine-stat
```
