# JUnit

!!! Info
    Disponible uniquement en édition Pro. (Canopsis ≥ [4.3.0](../../../../notes-de-version/4.3.0.md))

Ce module est capable de recevoir des résultats d'exécution de scénarios au format XML [JUnit](https://fr.wikipedia.org/wiki/JUnit).

Il comprend :

* Un récepteur (via API) de fichiers XML au format JUnit.
* Un moteur capable de parser, générer des alarmes à partir des résultats reçus.
* Un widget pour l'interface graphique capable de présenter les résultats sous diverses formes.

![junit-screenshot1](../../../../notes-de-version/img/4.3.0-junit-screenshot1.png){: .link width=80%"}

## Sommaire

1. [Installation](#installation)<br>
 A. [Docker](#docker)<br>
 B. [Paquet](#paquet)<br>
2. [Configuration](#configuration)<br>
3. [Création du widget](#creation-du-widget)<br>
 A. [Vue](#vue)<br>
 B. [Widget](#widget)<br>
4. [Publication de résultat](#publication-de-resultat)<br>

## Installation

### Docker

Depuis la version [4.4.0](../../../../notes-de-version/4.4.0.md), le connecteur `connector-junit` est installé et démarré avec Canopsis.

### Paquet CentOS 7

<!--- TODO --->

## Configuration

La configuration peut être changée dans le fichier `canopsis-pro.toml`.

*Configuration par défaut :*
```
[Canopsis.file]
…
# Local storage for Junit artifacts.
Junit = "/opt/canopsis/var/lib/junit-files"
# Temporary local storage for Junit data which are uploaded by API.
JunitApi = "/tmp/canopsis/junit"
```

## Création du widget

### Vue

Cliquer sur le bouton **Paramètres**, puis sur le bouton **Créer une vue** :

![Création vue 1/2](./img/vue1.png)

Remplir le formulaire :

![Création vue 2/2](./img/vue2.png)

### Widget

Aller dans la vue créée. Dans le menu latéral, cliquer sur **Ajouter un widget** :

![Création widget 1/3](./img/widget1.png)

Sélectionner le *widget* **Scénarios JUnit** :

![Création widget 2/3](./img/widget2.png)

Configurer le scénario en activant la réponse de l’API :

![Création widget 3/3](./img/widget3.png)

État du widget à sa création :

![État du nouveau widget](./img/widget4.png)

## Publication de résultat

Publier les résultats via l’API de Canopsis :

![API](./img/api.png)

*Par exemple via cURL :*
```
curl --location \
--request POST 'http://<canopsis.url>:8082/api/v4/junit/upload/reports?authkey=<01234567-89ab-cdef-0123456789ab>' \
--header 'Content-Type: multipart/form-data' \
--form 'files=@"</path/to/result.xml>"'
```

Le retour de l'API doit être le suivant en cas de succès :
```
{"upload_errors":{}}
```

Le résultat sera affiché lors du rafraichissement du *widget* :

![Tableau de bord](./img/tableaudebord.png)

*Exemple de résultat publié :*

![Résultat 1/3](./img/resultat1.png)

![Résultat 2/3](./img/resultat2.png)

![Résultat 3/3](./img/resultat3.png)
