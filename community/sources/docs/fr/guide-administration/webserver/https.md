# Configuration HTTPS

Canopsis utilise NGINX comme serveur web.
À partir de Canopsis 4.4.0, il est recommandé d’activer HTTPS car la fonctionnalité de healthcheck en est dépendante.
Votre SI doit supporter HTTP/1.1 ainsi que les Web sockets.

## Docker

### Certificats SSL

HTTPS ayant besoin de certificats SSL pour fonctionner, deux options s’offrent à vous :

* Créer ou obtenir un certificat SSL valide vous-même (Via votre propre autorité de certification par exemple). (Recommandé)
* Utiliser un certificat auto signé.

#### Certificats auto signés

Dans le cas de certificats auto signé vous n’avez rien à faire, le conteneur NGINX va en générer automatiquement lors de son démarrage.

!!! warning
	Vous devrez renouveler les certificats auto-signés tout les deux ans.

!!! info
	L’utilisation de certificats auto-signés provoquera l’affichage d’un message dans votre navigateur lors de la connexion à Canopsis.
	Vous devrez *manuellement* ajouter une exception sur chaque navigateur ou vous voudrez accéder à Canopsis.



#### Certificats personnalisés

Dans le cas où vous voulez utiliser vos propres certificats, ils devront êtres au [format X.509](https://fr.wikipedia.org/wiki/X.509) et montés dans le conteneur NGINX au bon endroit :

| Fichier | Point de montage dans le conteneur |
| ------- | ---------------------------------- |
| Clef publique du certificat | `/etc/nginx/ssl/cert.crt` |
| Clef privée du certificat | `/etc/nginx/ssl/key.key` |

Exemple avec `./mon_certificat.crt` comme clef publique et `./ma_clef_privee.key` comme clef privée :
```
docker-compose.yml
```
```yaml hl_lines="15 16"
...

  nginx:
    image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}nginx:${CANOPSIS_IMAGE_TAG}
    ports:
      - "80:80"
      - "443:443"
    env_file:
      - compose.env
    depends_on:
      - "api"
    restart: unless-stopped
    volumes:
      #- nginxcerts:/etc/nginx/ssl
	  - ./mon_certificat.cert:/etc/nginx/ssl/cert.crt:ro
	  - ./ma_clef_privee.key:/etc/nginx/ssl/key.key:ro
...
```

### Configuration du conteneur NGINX

Pour activer HTTPS dans le conteneur vous devrez modifier ses variables d’environnement dans le fichier `compose.env`.

| Variable | Description |
| -------- | ----------- |
| CPS_SERVER_NAME | FQDN sur lequel Canopsis sera disponible (ex : `capensis.mon-si.fr`) |
| CPS_ENABLE_HTTPS | État d’activation de HTTPS. Si cette variable vaut `true` HTTPS sera activé, dans tous les autres cas il sera désactivé |

Exemple pour activer HTTPS avec `capensis.mon-si.fr` en FQDN :
```
compose.env
```
```
...
CPS_SERVER_NAME=capensis.mon-si.fr
CPS_ENABLE_HTTPS=true
...
```
