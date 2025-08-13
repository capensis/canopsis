!!! important
    Les fixtures doivent être importées uniquement sur une installation vierge de Canopsis.
    
    L'import ne fonctionnera pas si des données sont déjà présentes dans la base de données MongoDB.
    
    De plus, `reconfigure` embarquant des fixtures prédéfinies ; il faut donc s’assurer qu’une même collection ne soit pas présente dans deux fichiers `.yml`.
    
    Dans certains cas, il pourra être nécessaire de vider ou de supprimer certains fichiers fournis par défaut.

!!! warning "Attention"
    L'extension `.yml` doit être respectée, `.yaml` n'est pas supportée.

!!! note "Note"
    Selon le contexte, il faudra filtrer la sortie de l'export.
    Par exemple, dans le cas de la collection `user`, ne garder uniquement que les users locaux ou comptes de service.

=== "Docker Compose"

    Placer le retour de l'export dans un fichier, par exemple, dans `./files/canopsis/reconfigure/fixtures/fixtures.yml`

    Editer votre fichier `docker-compose.override.yml` et surcharger le service reconfigure comme ci-dessous : 

    ```yaml
    reconfigure:
    volumes:
        - ./files/canopsis/reconfigure/fixtures/fixtures.yml:/opt/canopsis/share/database/fixtures/fixtures.yml:ro

    ```

    Exécuter le déploiement du docker compose :

    ```console
    docker compose up -d
    ```

=== "Paquets EL"

    Placer le retour de l'export dans un fichier, par exemple, dans `/opt/canopsis/share/database/fixtures/fixtures.yml`

    Provisionner Canopsis :

    ```
    set -o allexport; source /opt/canopsis/etc/go-engines-vars.conf; /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -edition pro
    ```

=== "Helm"

    En helm, il est possible de réaliser l'import de 2 manières.

    **Avec un configmap :**

    Créer un configmap, par exemple `cps-fixtures` à l'aide du fichier contenant le retour de l'export précédemment réalisé. 
    
    ```console
    kubectl create configmap cps-fixtures --from-file=fixtures.yml=fixtures.yml
    ```

    Surcharger le job `reconfigure` comme ci-dessous :

    ```yaml
    reconfigure:
      additionalVolumes:
        - name: fixtures-volume # defined name into chart
          configMap:
            name: cps-fixtures
      additionalVolumeMounts:
        - name: fixtures-volume 
          mountPath: "/opt/canopsis/share/database/fixtures/fixtures.yml"
          subPath: "fixtures.yml"
    ```

    Déployer votre instance de Canopsis : 

    ```console
    helm install canopsis-prod canopsis/canopsis-pro -f customer-values.yml
    ```

    **Avec un initcontainer, qui permet de récupérer l'export directement dans un fichier nommé `all_fixtures.yml` :**

    Surcharger le job `reconfigure` comme ci-dessous :

    ```yaml
    reconfigure:
      useFixtures: true
      initContainers:
        - name: fixture-download
          image: curlimages/curl:latest
          imagePullPolicy: Always
          command:
            - sh
            - -c
            - >
              curl -X POST https://demo.canopsis.net/api/v4/export-configuration
              -o /opt/canopsis/share/database/fixtures/all_fixtures.yml
              -H 'accept: application/x-yaml'
              -H 'Content-Type: application/json'
              -u root:root
              -d '{"export":[
              "configuration",
              "user",
              "role",
              "permission",
              "pbehavior",
              "pbehavior_type",
              "pbehavior_reason",
              "pbehavior_exception",
              "scenario",
              "metaalarm",
              "idle_rule",
              "eventfilter",
              "dynamic_infos",
              "playlist",
              "state_settings",
              "broadcast",
              "associative_table",
              "notification",
              "view",
              "view_tab",
              "widget",
              "widget_filter",
              "widget_template",
              "view_group",
              "instruction",
              "job_config",
              "job",
              "resolve_rule",
              "flapping_rule",
              "user_preferences",
              "kpi_filter",
              "pattern",
              "declare_ticket_rule",
              "link_rule",
              "map"]}'
          volumeMounts:
            - name: fixtures-volume
              mountPath: /opt/canopsis/share/database/fixtures
    ```

    Déployer votre instance de Canopsis : 

    ```console
    helm install canopsis-prod canopsis/canopsis-pro -f customer-values.yml
    ```
    
Dans les logs du reconfigure, les lignes suivantes doivent être présentes :

```
../cmd/canopsis-reconfigure/main.go:199 > start mongo fixtures
../cmd/canopsis-reconfigure/main.go:215 > finish mongo fixtures
```