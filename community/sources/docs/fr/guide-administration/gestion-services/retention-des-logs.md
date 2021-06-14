# Rétention des fichiers journaux

## Rotation des logs de Canopsis

Canopsis génère des fichiers journaux afin d'assurer la traçabilité des actions effectuées sur l'interface web ou par les moteurs.

Par défaut, aucune limite de place n'est mise en œuvre sur ces journaux.

Pour mettre en place une stratégie de rétention, et ainsi éviter une saturation d'espace disque disponible, vous pouvez appliquer la configuration logrotate suivante :

```sh
cat > /etc/logrotate.d/canopsis.conf << EOF
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

Vous pouvez personnaliser les valeurs `daily` et `rotate` pour ajuster la fréquence et la durée de la rétention.

## Rotation des logs de MongoDB

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

De la même façon que pour la rétention des journaux de Canopsis, celle de MongoDB peut être personnalisée en fonction de vos ressources.
