# Guide Administrateur : Installation

## Section : Installation / Conteneur

### Pré-requis

#### Ports

Il faut vérifier les ports utilisés par `docker-compose.yml`, si certains ports sont utilisés sur votre machines veuillez les libérer pour le bon déroulement de l'installation.

### Installation

Pour l'installation dockerisé de Canopsis, le porcédure est la suivante :

- clôner le dépôt Canopsis : https://git.canopsis.net/canopsis/canopsis
- Dans ce dépot un fichier `docker-compose.yml` est présent. Il va servir à la création de votre Canopsis en version Dockerisé.
- faire la commande suivante : `docker-compose up -d`

### Vérification

La vérification va passer par la commande `docker-compose ps`, elle va lister les conteneurs présents. Voici la liste des conteneur devant-être présent :

**TODO**