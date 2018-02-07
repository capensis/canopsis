# Gestion des services de Canopsis

Sur une installation de Canopsis à l'aide du script ./build-install.sh, la gestion des services se fait via la commande hypcontrol.

## Hypcontrol
La commande hypcontrol permet de démarrer les services canopsis.


```bash
hypcontrol start
hypcontrol stop
hypcontrol status
hypcontrol restart
```

## Gestion des erreurs

La commande hypcontrol permet à premier niveau de vérifier l'état de chaque service de Canopsis.

- Si le service est indiqué en RUNNING alors le service est OK
- Si le service est indiqué comme étant un état autre que RUNNING, alors il faut aller voir les logs de Canopsis pour obtenir plus de détails. Ces derniers se trouvent dans /opt/canopsis/var/log/

## Avertissement

La commande Hypcontrol n'est pas fonctionnelle en cas d'installation Ansible.
Il faut dans ce cas utiliser les commandes [service](service.md)
