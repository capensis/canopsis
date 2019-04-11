*Ce fichier est à intégrer à la documentation de `che`*

L'event-filter peut utiliser des sources de données externes pour enrichir les
évènements. Ces sources externes (à l'exception de `entity`) sont des plugins
disponibles dans Canopsis CAT.

Les plugins doivent-être placés dans un dossier accessible par le moteur `che`.

L'exécutable `engine-che` accepte une option `-dataSourceDirectory` permettant
de préciser le dossier contenant les plugins. Par défaut, ce dossier est celui
contenant `engine-che`.

## Docker

Les plugins doivent-être ajoutés dans un volume dans l'image docker, et leur
emplacement doit-être précisé dans la commande. Par exemple, avec
`docker-compose` :

```yaml
  che:
    image: canopsis/engine-che:${CANOPSIS_IMAGE_TAG}
    env_file:
      - compose.env
    restart: unless-stopped
    command: /engine-che -dataSourceDirectory /data-source-plugins
    volumes:
      - "./plugins:/data-source-plugins"
```
