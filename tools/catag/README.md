## Catag - Canopsis Tag

### Prérequis

 * Avoir les droits pour : ajouter **et** supprimer des tags sur les projets à configurer.
 * Créer un token d’accès GitLab : https://git.canopsis.net/profile/personal_access_tokens

### Build

```
glide install && go build .
```

### Utilisation

Le fichier INI contient les dépôts et la branche ou le commit sur lequel mettre le tag.

```
./catag -tag <tag> -token <gitlab access token>
```
