# Sécurisation d'une installation de Canopsis et de ses composants

## SELinux

Ceci concerne majoritairement les environnements CentOS et assimilés.

SELinux n'est pas pris en charge par Canopsis.

Il est nécessaire de le mettre en mode permissif ou de le désactiver sur les nœuds où vous installez Canopsis :

```sh
setenforce 0
sed -i 's/^SELINUX=.*$/SELINUX=permissive/' /etc/selinux/config
```

Puis, redémarrer le système.

## Sécurisation réseau

Pour la configuration du filtrage réseau, se référer à la [matrice des flux réseaux](../matrice-des-flux-reseaux/index.md). Vous y trouverez le détail des flux réseau à autiriser ou restreindre.

