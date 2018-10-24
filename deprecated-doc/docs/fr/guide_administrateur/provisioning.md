# Provisioning

Avec Docker, les images `canopsis/canopsis-prov` et `canopsis/canopsis-cat-prov` disposent d’un volume monté sur `/provisioning`.

Vous pouvez alors exécuter cette image avec l’option `-v source:/provisioning` de la commande `docker run`.

Sur votre machine hôte le dossier `source` pourra contenir :

 * Un dossier `mongo` contenant une arborescence comparable à ce qu’on trouve dans `sources/db-conf/opt/mongodb/load.d/`. Les fichiers seront copiés avec `rsync` au lancement de l’image.
