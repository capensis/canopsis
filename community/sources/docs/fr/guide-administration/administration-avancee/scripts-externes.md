Charger des scripts JavaScript externes

!!! info
    Cette fonctionnalité permet d’injecter des scripts JavaScript externes (par exemple Matomo, TagCommander, AppDynamics, Umami, etc.) directement dans l’interface web de Canopsis.
    Cela permet notamment :
    
    * de réaliser des tests de type APM (Application Performance Monitoring),
    * de collecter des statistiques déportées,
    * ou d’ajouter des fonctionnalités spécifiques de suivi utilisateur.

## Principe

Un fichier `external.js` est mis à disposition via nginx et chargé automatiquement par l’UI de Canopsis.

Ce fichier contient un tableau global `window.EXTERNAL_SCRIPTS`, dans lequel chaque élément est une balise `<script>` complète (inline ou avec attribut src).

Exemple minimal :

```js
window.EXTERNAL_SCRIPTS = [
  '<script type="text/javascript" src="https://host.com/some-script.js"></script>'
];
```

L’interface insérera ensuite dynamiquement ces balises dans le DOM.

## Déploiement

### Monter un fichier external.js

Dans votre docker-compose.override.yml :

```yaml
nginx:
  image: canopsis/nginx:externaljs
  ports:
    - "80:80"
    - "443:443"
  env_file:
    - compose.env
  restart: unless-stopped
  volumes:
    - nginxcerts:/etc/nginx/ssl
    - ./files/nginx/external.js:/opt/canopsis/srv/www/scripts/external.js
```

### Exemple de contenu external.js

```js
/**
 * We may put only script tags into this array (NOT SRC only the WHOLE script tags)
 *
 * @example window.EXTERNAL_SCRIPTS = [
 *   '<script type="text/javascript" src="https://host.com/some-script.js"></script>'
 *   '<script type="text/javascript">console.log("Inline script")</script>'
 *   `<script>
 *      console.log("Multi-line inline script too")
 *    </script>`
 * ];
 */

window.EXTERNAL_SCRIPTS = [
   '<script defer src="https://cloud.umami.is/script.js" data-website-id="64775368-f8ba-4ae9-9de7-7d2af7cddafe"></script>'
];
```

## Bonnes pratiques

* Toujours charger les scripts en asynchrone (async ou defer) pour éviter tout blocage de l’UI.
* Vérifier les politiques de sécurité de votre organisation (CSP, RGPD, etc.) avant d’ajouter des trackers.

## Exemples d’usage

**TagCommander :**

```js
window.EXTERNAL_SCRIPTS = [
'<script type="text/javascript" src="//cdn.tagcommander.com/xxx/Canopsis_xx.js" async></script>'
];
```

**Umami :**

```js
window.EXTERNAL_SCRIPTS = [
   '<script defer src="https://cloud.umami.is/script.js" data-website-id="xxx"></script>'
];
```
