# SeLinux  & Firewall  

**TODO (DWU) :** il y a un autre fichier dans ce dépôt qui décrit exactement la même procédure. Décider sur lequel on travaille.

## Selinux : 

Il est nécessaire de désactiver SeLinux pour installer correctement Canopsis.  
Dans tous les cas la commande `canoctl deploy` s'en charge automatiquement.  
Vous retrouverez cette commande lors de [l'installation par paquets](installation-paquets.md) ou de [l'installation via conteneurs](installation-conteneurs.md).  

## Firewall

Au niveau du Firewall, certains ports devront être ouverts. Voici une matrice des flux de Canopsis V3 qui répertorie les différents ports utilisés.  
Pour plus de détail, [rendez-vous ici](../administration-avancee/configuration-parefeu-et-selinux.md)

![CanoV3Flux](img/matrice-flux-canopsis.png)

