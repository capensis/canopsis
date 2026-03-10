# Procédure de migration dans un environnement Kubernetes/Helm

## Sommaire

[TOC]

## Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

Les commandes ci-dessous s'appliquent uniquement si votre instance MongoDB est hébergée sur votre cluster Kubernetes.
Dans le cas contraire, veuillez vous référer à l'onglet "Paquets RHEL 8".

```sh
export MONGODB_ROOT_PASSWORD=$(kubectl get secret canopsis-mongodb -o jsonpath='{.data.mongodb-root-password}' | base64 --decode)

kubectl exec canopsis-mongodb-0 -- mongosh -u root -p $MONGODB_ROOT_PASSWORD --eval 'db.adminCommand({ getParameter: 1, featureCompatibilityVersion: 1 })'
```

Le retour doit être de la forme `{ "featureCompatibilityVersion" : { "version" : "7.0" }, "ok" : 1 }`
Si ce n'est pas le cas, vous ne pouvez pas continuer la mise à jour.

## Arrêt de l'environnement en cours d'exécution

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

```sh
kubectl delete deployments --all
```

## Mise à jour de l'applicatif Canopsis

!!! information "Information"

    Canopsis 25.10 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations.  

!!! warning "Warning - Helm"

    Dans le cadre de la mise à jour vers Canopsis 25.10, les dépendances Helm pour MongoDB, RabbitMQ et Valkey ne s’appuient plus sur les charts Bitnami mais sur nos propres charts maintenus en interne. 

    La migration ne concerne donc pas uniquement Canopsis lui-même, mais également ces composants sous-jacents.

    Une attention particulière doit être portée aux valeurs de configuration et aux paramètres de persistance afin de garantir une transition fluide et sans perte de données.

## Mise à jour de Valkey (Helm uniquement)

Valkey ne s’appuie plus sur les charts Bitnami, mais sur nos propres charts maintenus en interne.
Il est donc nécessaire de supprimer le StatefulSet ainsi que le PVC associé.

Supprimer le Statefulset Valkey

```sh
kubectl get statefulset |grep valkey| awk {'print $1'}| xargs kubectl delete statefulset
```

Supprimer le PVC Valkey

```sh
kubectl get pvc |grep valkey| awk {'print $1'}| xargs kubectl delete pvc
```

## Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.15.1 à 2.21.4.  
En plus de la mise à jour de TimescaleDB lui-même, le système de gestion de base de données PostreSQL doit être mis à jour de la version 15 à la version 17.

Deux étapes sont à suivre :

1. Mise à jour de TimescaleDB 2.15.1 vers 2.21.4
2. Mise à jour de PostgreSQL 15 vers 17

Sauvegarder la base de données Canopsis :

```sh
kubectl exec canopsis-timescaledb-0 -- pg_dump postgresql://cpspostgres:canopsis@canopsis-timescaledb:5432/canopsis -Ft -f /tmp/postgres_canopsis_dump.tar

kubectl cp canopsis-timescaledb-0:/tmp/postgres_canopsis_dump.tar postgres_canopsis_dump.tar
```

Supprimer le Statefulset ainsi que le PVC :

```sh
kubectl get statefulset --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete statefulset

kubectl get pvc --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete pvc
```

Déploiement de PostgreSQL 17 - TimescaleDB 2.21.4

```sh
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
    name: canopsis-timescaledb
spec:
    selector:
    matchLabels:
        app.kubernetes.io/name: canopsis-pro
    serviceName: canopsis-timescaledb-headless
    updateStrategy:
    type: RollingUpdate
    template:
    metadata:
        labels:
        app.kubernetes.io/name: canopsis-pro
    spec:
        containers:
        - name: timescaledb
            image: docker.io/timescale/timescaledb:2.21.4-pg17
            ports:
            - containerPort: 5432
            env:
            - name: TIMESCALEDB_TELEMETRY
                value: "off"
            - name: POSTGRES_DB
                value: "canopsis"
            - name: POSTGRES_USER
                value: "cpspostgres"
            - name: POSTGRES_PASSWORD
                valueFrom:
                secretKeyRef:
                    name: canopsis-timescaledb
                    key: timescaledb-password
            readinessProbe:
            exec:
                command:
                - /bin/bash
                - -c
                - pg_isready -d \$POSTGRES_DB -U \$POSTGRES_USER
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 5
            volumeMounts:
            - name: datadir
                mountPath: /var/lib/postgresql/data
        imagePullSecrets:
        - name: canopsisregistry
    volumeClaimTemplates:
    - metadata:
        name: datadir
        annotations:
            helm.sh/resource-policy: "keep"
        spec:
        accessModes:
            - ReadWriteOnce
        resources:
            requests:
            storage: 8Gi
EOF
```

Restauration du dump :

```sh
kubectl cp postgres_canopsis_dump.tar canopsis-timescaledb-0:/tmp

kubectl exec canopsis-timescaledb-0 -- pg_restore --dbname=postgresql://cpspostgres:canopsis@canopsis-timescaledb-0:5432/canopsis --no-owner -Ft -v /tmp/postgres_canopsis_dump.tar
```

Suppression du Statefulset :

```sh
kubectl get statefulset --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete statefulset
```

## Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 7.0 à 8.0.  

!!! warning Attention
        Ce bloc est réservé uniquement aux environnements impliquant MongoDB exécuté dans un environnement Kubernetes.
        
        Si ce n'est pas votre cas, référez-vous au bloc [Paquets RPM](#__tabbed_1_2)

Dump de la base de données Canopsis :
```sh
kubectl exec -n canopsis canopsis-mongodb-0 -- mongodump --uri="mongodb://cpsmongo:canopsis@localhost:27017/canopsis" --gzip --out /tmp/dump_canopsis.gz
```

Récupération en local du dump :
```sh
kubectl cp canopsis/canopsis-mongodb-0:/tmp/dump_canopsis.gz .
```

Arrêt des pods MongoDB :
```sh
    kubectl get statefulset --no-headers=true |grep mongodb| awk {'print $1'}| xargs kubectl delete statefulset
```

Suppresion des PVCs MongoDB :
```sh
kubectl get pvc --no-headers=true | awk '{print $1}' | grep mongodb | xargs kubectl delete pvc
```

!!! warning Attention
    Veillez à bien adapter la commande ci-dessous avec vos paramètres présent dans votre fichier de surcharges, par exemple customer-values.yml.

Mise à jour de MongoDB :
```sh
helm repo update

helm upgrade canopsis canopsis/mongodb \
--set enabled=true \
--set settings.rootUsername=root \
--set settings.rootPassword=root \
--set replicaSet.enabled=true \
--set replicaSet.name=rs0 \
--set replicaSet.secondaries=0 \
--set replicaSet.key="c29tZXJhbmRvbXN0cmluZzEyMzQ1Ng==" \
--set userDatabase.name=canopsis \
--set userDatabase.user=cpsmongo \
--set userDatabase.password=canopsis \
--set storage.keepPvc=true \
--set storage.requestedSize=8Gi \
--version 1.1.0
```

Lorsque le POD est UP, copie du dump de la DB sur l'instance 0 de MongoDB : 
```sh
kubectl cp ./canopsis canopsis-mongodb-0:/tmp/
```

Restauration du dump :
```sh
kubectl exec -n canopsis canopsis-mongodb-0 -- mongorestore -u cpsmongo --password canopsis --gzip --db canopsis /tmp/canopsis
``` 

Suppression du Statefulset :
```sh
kubectl get statefulset --no-headers=true | awk '{print $1}' | grep mongo | xargs kubectl delete statefulset
```

## Mise à jour de RabbitMQ

Dans cette version de Canopsis, le bus de données RabbitMQ passe de la version 4.0 à 4.1.  

Supprimer le volume associé à RabbitMQ

```sh
kubectl get pvc --no-headers=true | awk '{print $1}' | grep rabbitmq | xargs kubectl delete pvc
```

## Lancement du provisioning `canopsis-reconfigure`

### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 

Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

### Séparation des flux d’événements par initiateur

En version 25.04, les flags suivants avaient été dépréciés, ils sont à présent obsolètes.  
Toute référence doit être supprimée dans vos configurations. Les moteurs concernés ne démareront pas sans cela.  

* -publishQueue
* -consumeQueue
* -workers (remplacé par des flags spécifiques à chaque type de flux)

| Moteur               | Flags nouveaux     | Valeur par défaut | Flags obsolètes                             |
|----------------------|--------------------|-------------------|---------------------------------------------|
| engine-fifo          | -workers           | 10                | -publishQueue, -consumeQueue                |
| engine-che           | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
| engine-axe           | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-correlation   | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-dynamic-infos | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-action        | -externalWorkers   | 4                 | -workers                                    |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |

## Mise à jour et démarrage final de Canopsis


Définir le nom de votre instance

```sh
export RELEASE_NAME="canopsis-prod"
```

Mise à jour des repos helm

```sh
helm repo update
```

Mise à jour de Canopsis

```sh
helm upgrade ${RELEASE_NAME} canopsis/canopsis-pro -f customer-values.yaml
```

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/25.10.0-healthcheck.png)