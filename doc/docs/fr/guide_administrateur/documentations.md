# Gestion de la documentation utilisateurs

La documentation utilisateurs est une documentation délivrée directement depuis
canopsis pour aider les utilisateurs dans leur utilisation de canopsis.

## Ajout de page documentation
Pour créer une nouvelle page de documentation, il fait créer un nouveau
fichier au format markdown dans le répertoire
``${CPS_HOME}/var/www/documentation`` de toutes vos instances de webserver.

Attention, il faut **absolument** que le fichier soit encodé en **UTF-8**. Dans
le cas contraire, la génération de la documentation échouera ! Et le fichier
doit absolument avoir pour extension ``.md`` sans cela, la page de documentation
sera inaccessible.

## Consultation de la documentation

Pour consulter la documentation, il suffit de faire une requête **HTTP GET**
sur la route suivante : ```/api/v2/documentation/<NAME>```, avec ``<NAME>``
le nom du de la page de documentation sans l'extension.
