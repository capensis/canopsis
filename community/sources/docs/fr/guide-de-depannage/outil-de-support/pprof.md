# Pprof

`pprof` est un outil intégré aux moteurs Go de Canopsis permettant d’analyser finement leurs performances.

Il expose différents profils (CPU, mémoire, goroutines, contention…), très utiles pour identifier des lenteurs, des pics de charge inhabituels ou un comportement anormal d’un moteur.

L’interface pprof ouvre un petit serveur HTTP sur lequel ces profils peuvent être consultés en direct ou exportés pour une analyse plus poussée.

Par défaut, pprof est désactivé.

Cette page décrit la marche à suivre pour l’activer sur un moteur Go

=== "RPM"

    Exemple ici avec le moteur **action** :

    ```bash
    sudo systemctl edit canopsis-engine-go@engine-action.service
    ```

    Y ajouter les variables d’environnement suivantes :

    ```ini
    [Service]
    Environment="CPS_DEBUG_WEB_PPROF_ENABLE=1"
    ```

    Par défaut :

    * **HOST** : `localhost`
    * **PORT** : `6060`

    Recharger systemd

    ```bash
    sudo systemctl daemon-reload
    ```

    Redémarrer le moteur

    Toujours avec l’exemple du moteur action :

    ```bash
    sudo systemctl restart canopsis-engine-go@engine-action.service
    ```

    Tester

    Par exemple, exécuter la commande suivante :

    ```bash
    curl -s http://localhost:6060/debug/pprof/allocs?debug=1
    ```

=== "Docker"

    Dans un déploiement Docker, l’activation de pprof se fait en surchargeant la configuration du moteur concerné via le fichier `docker-compose.override.yml`.
    L’objectif est d’activer pprof, éventuellement modifier son port ou son interface d’écoute, puis exposer le port correspondant pour y accéder depuis l’extérieur.

    Modifier votre fichier `docker-compose.override.yml` :

    ```yaml
    action:
      command: /engine-action
      environment:
        CPS_DEBUG_WEB_PPROF_ENABLE: "1"
      ports:
        - "6060:6060"
    ```

    Dans cet exemple :

    * `CPS_DEBUG_WEB_PPROF_ENABLE` active l’endpoint pprof.
    * Le moteur écoute par défaut sur `localhost:6060` **dans le conteneur**.
    * Le port `6060` est exposé en local → accessible via `http://localhost:6060/debug/pprof`.

    Appliquer la configuration :

    ```bash
    docker compose down action
    docker compose up action -d
    ```

    Une fois le service relancé, pprof est accessible via :

    ```
    http://localhost:<port>/debug/pprof/
    ```

=== "Helm"

    Lorsque Canopsis est déployé via Helm, l’activation de pprof se fait en ajoutant les variables d’environnement nécessaires dans la section `env:` du moteur souhaité.
    
    Les charts des engines ne prévoient pas l’ouverture de ports supplémentaires, pprof sera donc accessible uniquement depuis l’intérieur du cluster ou via un port-forward.

    Dans votre fichier `customer-values.yaml`, ajouter les variables d’environnement dans le bloc `env:` du moteur voulu.

    Exemple pour l’engine **action** :

    ```yaml
    engine:
      action:
        env:
          - name: CPS_DEBUG_WEB_PPROF_ENABLE
            value: "1"
    ```
    
    Notes :

    * Par défaut pprof écoute sur `localhost:6060` dans le pod.
   
    Redéployer l’engine

    ```bash
    helm upgrade canopsis-prod -f customer-values.yaml canopsis/canopsis-pro
    ```

    Comme les ports ne sont pas exposés automatiquement, utiliser un port-forward :

    ```bash
    kubectl port-forward pod/<nom-du-pod> 6060:6060
    ```

    Puis ouvrir :

    ```
    http://localhost:6060/debug/pprof/
    ```