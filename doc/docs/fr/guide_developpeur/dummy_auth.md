# Dummy authentification

La dummy authentification est un plugin qui permet de tester l'authentification
externe de canopsis.

## Activation

Pour activer la dummy authentification, il suffit d'ajouter
`canopsis.auth.dummy` dans la liste provider de la section *auth*.

Et de charger les routes du module em ajoutant la ligne
`canopsis.webcore.services.dummy=1` dans la section *webservices*.

## Utilisation

Le plugin dummy authentification accepte n'importe quel couple
compte utilisateur/mot de passe à l'exception de l'utilisateur *davros* et
*mr.rabbit*.

Pour l'utiliser, il faut essayer de s'authentifier sus le formulaire de login
de canopsis pour être redirigé vers la page de login de dummy authentification.
