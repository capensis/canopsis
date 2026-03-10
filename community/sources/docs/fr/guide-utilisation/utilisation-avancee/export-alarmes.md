# Export d'alarmes au format CSV

L'utilitaire **export-alarms** permet d'exporter des alarmes Canopsis au format **CSV** via l'API.

Il est conçu pour être exécuté automatiquement (via un ordonnanceur de tâches) afin de produire des exports réguliers destinés à :

- Du reporting
- De l'archivage
- Des traitements externes
- L'alimentation d'outils de BI


## Principe de fonctionnement

`export-alarms` utilise l'API Canopsis pour :

1. Créer une tâche d'export d'alarmes
2. Attendre la génération du fichier
3. Récupérer les données
4. Ecrire le résultat dans un fichier CSV local

Chaque export est défini dans un **fichier de configuration YAML**.

Un fichier CSV est généré **par export configuré**.

Le nom du fichier suit le format :

    <label>-<timestamp>.csv

Exemple :

    opened-alarms-2026-03-01T12-00-00.csv

## Lancement de la commande

L'utilitaire peut être lancé directement :

``` bash
export-alarms --config /opt/canopsis-export-alarms/config.yml
```

Options disponibles :

L'option `-h` permet d'afficher toutes les options disponibles au lancement du service.

| Option            | Description                                   |
|-------------------|-----------------------------------------------|
| `-c`, `--config`  |  chemin du fichier de configuration           |
| `-d`, `--debug`   |  active le mode debug                         |
| `--logger`        |  destination des logs (`stderr`, `journald`)  |
| `-v`, `--version` |  affiche la version                           |
| `-h`, `--help`    |  affiche l'aide                               |


## Configuration

La configuration de l'outil est définie dans un **fichier YAML**.

### Section API

Cette section définit la connexion à l'API Canopsis.

``` yaml
api:
  host: http://localhost:8082
  http_timeout: 5s
  export_timeout: 1m
  username: ""
  password: ""
  insecure_skip_verify: false
```

!!! warning 
    Il est recommandé d'utiliser les variables d'environnement pour les identifiants.

Variables supportées :

    CPS_API_USERNAME
    CPS_API_PASSWORD
    CPS_API_URL
    CPS_API_HTTP_TIMEOUT


### Section CSV

Cette section configure la génération des fichiers CSV.

``` yaml
csv:
  separator: comma
  dir: /tmp/cps/export
  time_format: "YYYY-MM-DDThh:mm:ss"
```

Paramètres disponibles :

| Paramètre      | Description                                            |
|----------------|--------------------------------------------------------|
|  `separator`   |  séparateur CSV (`comma`, `semicolon`, `tab`, `space`) |
|  `dir`         |  répertoire de sortie                                  |
|  `time_format` |  format des dates dans le CSV                          |


## Définition des exports

Les exports sont définis dans la section `exports`.

Chaque export correspond à un fichier CSV généré.

Exemple :

``` yaml
exports:
  - label: newest
    open: true
    interval: 15m
    time_field: v.creation_date

    fields:
      - name: v.display_name
        label: Display name
      - name: v.state.val
        label: State
      - name: v.output
        label: Output
      - name: v.creation_date
        label: Creation
```

Paramètres disponibles :

|  Paramètre    | Description                                              |
| ------------- | -------------------------------------------------------- |
|  `label`      | préfixe du nom du fichier                                |
|  `open`       | `true` = alarmes ouvertes, `false` = alarmes résolues    |
|  `interval`   | fenêtre temporelle (`15m`, `1h`, `1d`, `1w`, `1M`, `6M`) |
|  `time_field` | champ utilisé pour le filtrage temporel                  |
|  `fields`     | liste des colonnes exportées                             |


## Exemple : export des alarmes résolues

``` yaml
- label: last_resolved
  open: false
  interval: 15m
  time_field: v.resolved

  fields:
    - name: v.display_name
      label: Display name
    - name: v.state.val
      label: State
    - name: v.output
      label: Output
    - name: v.creation_date
      label: Creation
    - name: v.resolved
      label: Resolved
    - name: v.duration
      label: Duration
```


## Filtres avancés

Les filtres avancés utilisent le **langage de patterns Canopsis**, qui permet de définir des règles de filtrage sur les alarmes, les entités ou les comportements périodiques.

Ce langage repose sur une structure de **conditions combinées avec des opérateurs logiques (ET / OU)** appliquées sur les différents attributs des objets Canopsis.

La description complète du langage, des champs disponibles et des types de conditions est disponible dans la page suivante :

[Langage des filtres et patterns](../../guide-developpement/filtres/index.md)

Cette documentation détaille notamment :

- la structure des patterns
- les champs disponibles selon le type de pattern (alarme, entité, événement, comportement périodique)
- les différents types de conditions (`eq`, `neq`, `regexp`, `exist`, `gt`, `lt`, etc.)
- les conditions spécifiques aux tableaux, aux durées et aux timestamps

Les exemples ci-dessous illustrent l'utilisation de ces patterns dans la configuration de l'outil `export-alarms`.

### Alarm pattern

``` yaml
alarm_pattern:
  - - field: v.state.val
      cond:
        type: eq
        value: 3
```

### Entity pattern

``` yaml
entity_pattern:
  - - field: type
      cond:
        type: eq
        value: resource
```


### Pbehavior pattern

``` yaml
pbehavior_pattern:
  - - field: pbehavior_info.canonical_type
      cond:
        type: eq
        value: active
```

### Recherche texte

``` yaml
search: "database"
```

## Gestion des fichiers

L'utilitaire ne gère **pas la rotation ni la suppression des fichiers**
générés.

La gestion de la rétention doit être assurée par :

- Un script externe
- Un cron dédié
- Un outil de gestion de fichiers

## Exécution

=== "RPM"

    ### RPM

    Si vous avez installé Canopsis via les paquets RPM, l'outil `export-alarms` est disponible directement sous forme de binaire.

    Un fichier d'exemple de configuration est fourni dans le paquet (`/opt/canopsis/share/config/export-alarms/config.yml.example`).
    
    Copiez le en `/opt/canopsis/share/config/export-alarms/config.yml` et adaptez-le à vos besoins.

    !!! warning "Important"

        Assurez-vous que le répertoire de destination existe (par défaut `/tmp/export-alarms`).
        Dans le cas contraire, le créer.

    #### Exécution manuelle

    Pour lancer l'export manuellement, utilisez la commande suivante :

    ```bash
    /opt/canopsis/bin/export-alarms
    ```

    Si l'exécution réussit, les fichiers CSV seront générés dans le répertoire configuré (par défaut `/tmp/export-alarms`).

    #### Exécution automatique

    Pour automatiser l'export, vous pouvez utiliser `cron` ou tout autre ordonnanceur de tâches afin d'exécuter la commande à intervalles réguliers.

    Par exemple, pour planifier un export quotidien à 8h00 via un fichier `crontab` :

    ```sh
    cat <<EOF > /etc/cron.d/export-alarms
    0 8 * * * root /opt/canopsis/bin/export-alarms
    EOF
    ```

    !!! info "Info"

        Les fichiers seront créés en root, pensez à adapter les permissions ou à exécuter la commande avec un utilisateur spécifique si nécessaire.

=== "Docker Compose"

    ### Docker

    Dans votre fichier `docker-compose.override.yml`, une section `export-alarms` est normalement préconfigurée comme suit :

    ```yaml
    export-alarms:
      <<: *initial_config_base
      user: "${UID}:${GID}"
      profiles:
        - support
      image: ${DOCKER_REPOSITORY}${CPS_EDITION}/export-alarms:${CANOPSIS_IMAGE_TAG}
      volumes:
        - ./files-pro/export-alarms/csv:/tmp/export-alarms
        - ./files-pro/export-alarms/config.yml:/opt/canopsis/share/config/export-alarms/config.yml
      command: /export-alarms
    ```

    #### Exécution manuelle

    Pour lancer le conteneur manuellement, exécutez la commande suivante :

    ```bash
    docker compose --profile support up export-alarms -d
    ```

    Une fois le service démarré, les fichiers CSV seront générés dans le répertoire hôte configuré, ici `./files-pro/export-alarms/csv`.

    #### Exécution automatique

    Pour automatiser cette tâche, vous pouvez configurer un `cron` sur l'hôte pour piloter le conteneur Docker.

    Exemple pour un export quotidien à 8h00 (pensez à adapter le chemin vers votre projet Canopsis) :

    ```sh
    cat <<EOF > /etc/cron.d/export-alarms
    0 8 * * * root cd /opt/canopsis && docker compose --profile support up export-alarms --force-recreate
    EOF
    ```

=== "Helm"

    ### Helm

    Lors d'une installation via Helm, le service `export-alarms` est géré nativement via un **CronJob** Kubernetes. Il n'y a pas d'exécution manuelle directe ; l'outil est conçu pour fonctionner de manière planifiée.

    #### Configuration de l'automatisation

    Pour activer et planifier l'export, modifiez ou surcharger votre fichier `values.yaml` :

    ```yaml
    exportAlarms:
      enabled: true
      schedule: "0 8 * * *" # À adapter selon la fréquence souhaitée
      storageClassName: "" # À renseigner pour persister et récupérer les CSV sur un volume partagé
    ```

    !!! warning "Important"

        Avant le déploiement, assurez-vous que la configuration de l'outil est correctement définie dans votre `config.yml` ou via un ConfigMap dédié.
