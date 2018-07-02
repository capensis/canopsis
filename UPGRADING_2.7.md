# Mise à jour Canopsis 2.7

## Redis

Canopsis nécessite maintenant l'installation d'un serveur redis pour pouvoir fonctionner.

Une fois installé, il faut exporter la variable d'environnement suivante dans les conteneurs docker :
```bash
CPS_REDIS_URL=redis://adresse_de_redis:6379/0
```
