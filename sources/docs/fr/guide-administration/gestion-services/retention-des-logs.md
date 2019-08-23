# Rétention des fichiers journaux

## Canopsis

Canopsis génère des fichiers journaux afin d'assurer la traçabilité des actions
effectuées sur l'interface web ou par les moteurs.

Par défaut, aucune limite de place n'est mise en œuvre sur ces journaux.

Pour mettre en place une stratégie de rétention, et ainsi éviter une saturation
d'espace disque disponible, vous pouvez appliquer la commande suivante :

```bash
cat << EOF | sudo tee /etc/logrotate.d/canopsis.conf
/opt/canopsis/var/log/*.log /opt/canopsis/var/log/engines/*.log {
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

Vous pouvez personnaliser les valeurs `daily` et `rotate` pour ajuster la
fréquence et la durée de la rétention.

## MongoDB

MongoDB, la base de données utilisée par Canopsis produit également des fichiers
journaux qu'il convient de limiter.

Ceci peut être réalisé grâce à cette commande :

```bash
cat << EOF | sudo tee /etc/logrotate.d/mondodb-server.conf
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

De la même façon que pour la rétention des journaux de Canopsis, celle de
MongoDB peut être personnalisée en fonction de vos ressources.
