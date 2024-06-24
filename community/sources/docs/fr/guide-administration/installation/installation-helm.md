# Installation de Canopsis avec Helm

Cette procédure décrit l'installation de Canopsis avec Helm.

## Prérequis

### Cluster Kubernetes

Pour déployer Canopsis en utilisant Helm, il est nécessaire d'avoir un cluster [Kubernetes](https://kubernetes.io/) opérationnel.

À ce jour, le déploiement a été testé sur la dernière version de `Kubernetes (1.30)` et de `Helm (3.15)`.

### Helm

[Helm](https://helm.sh/docs/intro/install/) est un gestionnaire de paquets pour Kubernetes.

### Kubectl
[Kubectl](https://kubernetes.io/fr/docs/tasks/tools/install-kubectl/) est un outil en ligne de commande pour contrôler des clusters Kubernetes.

### Git
[git](https://git-scm.com/downloads) est un logiciel de gestion de versions open source.


## Récupérer le chart en utilisant le repo Helm

### L'accès au repo 

Les charts Helm sont un moyen de déployer Canopsis Pro et certains connecteurs. Ce moyen de déploiement est réservé aux clients disposant d'une souscription Canopsis Pro.

Le repository est accessible via le lien suivant : [https://git.canopsis.net/helm/charts/charts-repo](https://git.canopsis.net/helm/charts/charts-repo)

Utilisez un token d'accès utilisateur dans GitLab associé à votre compte client. Ceci évitera de stocker votre mot de passe personnel lors de la configuration du dépôt dans `helm`. Le token doit être créé avec la permission `read_api` seulement.

[Comment créer un token d'accès Gitlab ?](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

### Ajouter le repo
Définir une variable avec votre token GitLab :
```sh
TOKEN=your-user-token-for-helm-app
```

Ajouter le repo en version **stable** :
```sh
helm repo add --username oauth2 --password $TOKEN canopsis \
     https://git.canopsis.net/api/v4/projects/603/packages/helm/stable
```

Mettre à jour les repos :
```sh
helm repo update
```

Vérifier les charts disponibles grâce au repo:
```sh
helm search repo canopsis
```

!!! Information
    ___En option___, ajouter le repo en version devel :
    ```sh
    helm repo add --username oauth2 --password $TOKEN canopsis-devel \
         https://git.canopsis.net/api/v4/projects/603/packages/helm/devel
    ```
    Mettre à jour les repos :
    ```sh
    helm repo update
    ```
    Vérifier les charts disponibles grâce au repo:
    ```sh
    helm search repo canopsis-devel
    ```

## Configuration de l'environnement Kubernetes

Créer un **namespace dédié** (exemple : `canopsis`) :
```
kubectl create namespace canopsis
```

Il est possible de définir le *namespace* `canopsis` par défaut pour la suite du maniement des commandes de déploiement :
```
kubectl config set-context --current --namespace=canopsis
```

Se connecter avec ses identifiants sur le **registry Docker** `docker.canopsis.net` :
```
docker login docker.canopsis.net
```

Créez un **secret** Kubernetes pour stocker les identifiants du registry Docker. Ceux-ci sont indispensables à la récupération des images soumises à souscription Canopsis Pro.

Dans l'exemple ci-dessous, l'objet *Secret* sera nommé `canopsisregistry`.

!!! Remarque
    La création dudit secret est également détaillée dans la documentation de Kubernetes : [Créez un Secret basé sur les identifiants existants](https://kubernetes.io/fr/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials)

```
kubectl create secret generic canopsisregistry \
    --from-file=.dockerconfigjson=$HOME/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```
## Ingress controller

Pour exposer Canopsis, vous pouvez utiliser un ingress adapté à vos besoins, permettant ainsi un accès direct depuis votre nom de domaine tout en garantissant la sécurité des communications avec un certificat SSL.

## Surcharge des valeurs du fichier Helm

!!! note "Remarque"
    Vous pouvez remplacer tout autre paramètre activé dans le fichier des valeurs du chart, comme indiqué dans le [README](https://git.canopsis.net/helm/charts/canopsis-pro/-/tree/develop?ref_type=heads) du chart. Ce qui suit est fourni à titre d'exemple ; il est important de personnaliser et de surcharger les valeurs du fichier Helm en fonction de vos besoins spécifiques et des configurations requises pour votre environnement. Assurez-vous de bien comprendre chaque paramètre afin de garantir que le déploiement répond à vos attentes et exigences opérationnelles.

    **Si vous n'avez pas accès au README, merci de vous rapprocher de votre référent client.**

=== "Prod"
    Exemple de fichier ```customer-value.yaml``` : 
      ```
      imagePullSecrets:
        - name: canopsisregistry

      backendUrls:
        mongodb: mongodb://cpsmongo:canopsis@mongodb:27017/canopsis
        rabbitmq: amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis
        redis: redis://:canopsis@redis:6379/0
        timescaledb: postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis

      ```
    
    !!! warning "Recommendation"
        En production, il est recommandé d'exécuter les services backend de Canopsis (bases de données, RabbitMQ) sur des serveurs dédiés, et non dans des conteneurs. Il faut donc veillez à définir les URL backend appropriées comme ci-dessus. 

=== "Lab"
    Exemple de fichier ```customer-value.yaml``` : 
    ```
    imagePullSecrets:
      - name: canopsisregistry

    mongodb:
      enabled: true
    rabbitmq:
      enabled: true
    redis:
      enabled: true
    timescaledb:
      enabled: true
    ```


## Déploiement

Choisissez un nom de *release* Helm correspondant à votre instance Canopsis Pro (cano0, cano1, canopsislab, ...), ici nous utiliserons "canopsis-prod". Ensuite, vous devrez l'utiliser de manière appropriée dans toutes les commandes Helm pour cette instance. 

Les commandes d'exemple ci-dessous utiliseront la variable ${RELEASE_NAME} à cette fin.
Le nom de la version est à votre choix ; il permet plusieurs déploiements du même chart pour différentes instances sur le même cluster Kubernetes.

!!! info "Information"
    Les noms de version doivent ressembler à des noms DNS :

    - uniquement des caractères alphanumériques en minuscules ou '-';
    - doivent commencer et se terminer par un caractère alphanumérique ;
    - longueur maximale de 53 caractères.

Définir le nom de votre instance : 
```sh
export RELEASE_NAME="canopsis-prod"
```

Initier le déploiement :
    ```sh
    helm install ${RELEASE_NAME} canopsis/canopsis-pro -f customer-value.yaml
    ```
  
Superviser le déploiement :
```
watch kubectl get pod
```

### Accéder à l'interface Web :

Si vous utilisez l'ingress, il sera directement possible d'accéder à l'IHM de Canopsis en y accédant via le nom de domaine configuré, par exemple : 
```sh
https://canopsis.k8s.lan
```

Ou si vous n'utilisez pas l'ingress :
```sh
kubectl port-forward svc/${RELEASE_NAME}-nginx 8443:443
```

!!! Information
    Cette commande ouvrira le port 8443 en local et redirigera les connexions vers le port 443 du service Nginx

## Désinstaller Canopsis : 
```
helm uninstall ${RELEASE_NAME}
```

### Supprimer les volumes persistants

!!! danger "Attention"
    Cette commande est destructrice, bien vérifier les PVC présents avant son exécution.

```
kubectl delete pvc --all
```

