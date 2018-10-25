# Installation / SeLinux  & Firewall  

**TODO (DWU) :** il y a un autre fichier dans ce dépôt qui décrit exactement la même procédure. Décider sur lequel on travaille.

### Selinux : 

Il est nécessaire de désactiver SeLinux pour installer correctement Canopsis. Dans tous les cas la commande `canoctl deploy` s'en charge automatiquement.  
Vous retrouverez cette commande lors de [l'installation](/doc-ce/Guide Administrateur/Installation/Packages.md)

### Firewall

Au niveau du Firewall, certains ports devront être ouverts. Voici une matrice des flux de Canopsis V3 qui répertorie les différents ports utilisés. Pour plus de détail, rendez-vous [ici](/doc-ce/Guide Administrateur/Administration%20avanc%C3%A9e/Firewall%20&%20SeLinux%20avanc%C3%A9.md)

![CanoV3Flux](/doc-ce/Guide%20Administrateur/Installation/Images/Matrice%20des%20flux%20Canopsis%203.png)

