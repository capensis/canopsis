# Installation de Canopsis avec Helm

Cette procédure décrit l'installation de Canopsis avec Helm.

## Prérequis

### Cluster Kubernetes

Pour déployer Canopsis en utilisant Helm, il est nécessaire d'avoir un cluster [Kubernetes](https://kubernetes.io/) opérationnel. 

### Helm

[Helm](https://helm.sh/docs/intro/install/) est un gestionnaire de paquets pour Kubernetes.

### Kubectl
[Kubectl](https://kubernetes.io/fr/docs/tasks/tools/install-kubectl/) est un outil en ligne de commande pour contrôler des clusters Kubernetes.

### Git
[git](https://git-scm.com/downloads) est un système de contrôle de version open source

!!! information
    Il est possible de récupérer le chart de deux manières :

    - En utilisant le **repo** Helm Canopsis ;
    
    - En utilisant les **sources** du chart.


## Compatibilité

Des tests de compatibilités ont été réalisés sur les versions suivantes : 

| Kubernetes | kubectl | helm   |
|:----------:|:-------:|:-----: |
|   1.26     |  1.27.4 | 0.23.0 |
|   1.30     |  1.30.1 | 3.15.0 |

## Récupérer le chart en utilisant le repo Helm

### L'accès au [repo](https://git.canopsis.net/helm/charts/charts-repo)

Les charts Helm sont un moyen de déployer Canopsis Pro et les connecteurs sélectionnés, réservé aux clients abonnés à Canopsis Pro. Le dépôt est donc privé. L'utiliser implique une authentification et une autorisation appropriées.

Tout utilisateur ayant besoin d'accéder au dépôt des charts Helm de Canopsis doit être ajouté en tant que membre de ce projet, avec au moins le rôle de Reporter.

Il est recommandé d'utiliser un token d'accès utilisateur au de niveau GitLab évitez de divulguer votre mot de passe personnel lors de la configuration du dépôt dans le client Helm : possibilité de définir une date d'expiration, meilleur contrôle des permissions suffisantes.

Le champ d'application du token doit être uniquement read_api.

[Comment créer un token d'accès Gitlab ?](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

### Ajouter le repo
Exporter votre token Gitlab :
```sh
TOKEN=your-user-token-for-helm-app
```

Ajouter le repo en version **stable** :
```sh
helm repo add --username oauth2 --password $TOKEN canopsis \
https://git.canopsis.net/api/v4/projects/603/packages/helm/stable
```

!!! Remarque
    ___En option___, ajouter le repo en version devel :
    ```sh
    helm repo add --username oauth2 --password $TOKEN canopsis-devel https://git.canopsis.net/api/v4/projects/603/packages/helm/devel

    ```

Mettre à jour les repos :
```sh
helm repo update
```

Vérifier si les repos ont bien été mis à jour :
```sh
helm search repo canopsis
helm search repo canopsis-devel
```


## Récupérer le chart en utilisant les sources Helm

=== "Cloner le repo en utilisant le protocole SSH"
    ```sh
    git clone git@git.canopsis.net:helm/charts/canopsis-pro.git
    ```

=== "Cloner le repo en utilisant le protocole HTTPS"
    ```sh
    git clone https://git.canopsis.net/helm/charts/canopsis-pro.git
    ```

Se rendre dans le dossier des sources :
```
cd canopsis-pro
```

Si votre installation de Helm ne connaît pas encore le dépot Bitmani, ajouter le :
```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Builder les dépendances :
```
helm dependency build
```

## Configuration de l'environnement Kubernetes

Créer un **namespace canopsis** :
```
kubectl create namespace canopsis
```

Définir le **namespace canopsis** comme namespace **par défaut** :
```
kubectl config set-context --current --namespace=canopsis
```

Se connecter avec ses identifiants sur le **repo Docker** Gitlab **docker.canopsis.net** :
```
docker login docker.canopsis.net
```

Créez le **secret** Kubernetes à partir de vos identifiants Docker. Vous pouvez lui donner le nom que vous souhaitez (dans l'exemple ci-dessous, "canopsisregistry"). Dans tous les cas, vous le référencerez dans la valeur de imagePullSecrets à l'étape suivante.

!!! Remarque
    La création dudit secret est également détaillée dans la documentation de Kubernetes : [Créez un Secret basé sur les identifiants existants](https://kubernetes.io/fr/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials)

```
kubectl create secret generic canopsisregistry \
    --from-file=.dockerconfigjson=$HOME/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```
### Déploiement nginx-ingress controller

Le contrôleur NGINX Ingress permet de gérer l'accès aux applications hébergées dans un cluster Kubernetes en exposant ces applications au monde extérieur à travers des points d'entrée HTTP et HTTPS, tout en offrant des fonctionnalités avancées comme le routage basé sur les chemins et les hôtes, la répartition de charge, la terminaison SSL et plus encore.

Déploiement de l'ingress controller, ici nginx : 
```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

Vous pouvez surveillez la progression du déploiement avec la commande : 
```sh
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

Une fois le déploiment terminé, la commande vous rend la main.

Par la suite, 2 choix s'offrent à vous pour exposer l'IHM Canopsis : 

  - Ingress-controller ( voir le point ci-dessous concernant l'ingress)
  - [Port forwarding](#acceder-a-linterface-web)


Si vous faites le choix de l'ingress, il vous faudra disposer d'un nom de domaine, ainsi qu'un certificat SSL associé.

Ici, nous utiliserons le nom de domaine **canopsis.k8s.lan**. La résolution DNS se fera en modifiant notre fichier ```/etc/hosts```

Editer ledit fichier :
```sh
sudo vim /etc/hosts
```
Y ajouter la ligne:
```sh
127.0.0.1 canopsis.k8s.lan
```

Pour générer le certificat SSL auto-signé, créer une clé privée : 
```sh
openssl genrsa -out canopsis.k8s.lan.key 2048
```

Créer une demande de signature de certificat (CSR):
```sh
openssl req -new -key canopsis.k8s.lan.key -out canopsis.k8s.lan.csr
```
Des informations vous sera demandés comme ci-dessous : 
```sh
Country Name (2 letter code) [XX]:FR
State or Province Name (full name) []:Île-de-France
Locality Name (eg, city) [Default City]:Paris
Organization Name (eg, company) [Default Company Ltd]:Capensis
Organizational Unit Name (eg, section) []:IT Department
Common Name (eg, your name or your server's hostname) []:canopsis.k8s.lan
Email Address []:admin@canopsis.k8s.lan
```

Ensuite, générer le certificat auto-signé : 
```sh
openssl x509 -req -days 365 -in canopsis.k8s.lan.csr -signkey canopsis.k8s.lan.key -out canopsis.k8s.lan.crt
```

Pour que le certificat puisse être utilisé par l'ingress, il est nécessaire de créer un secret en procédant de la manière suivante :
```sh
kubectl create secret tls canopsis-front-tls-secret --cert=canopsis.k8s.lan.crt --key=canopsis.k8s.lan.key
```

## Surcharge des valeurs du fichier Helm


=== "Lab - sans ingress"
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

=== "Lab - avec ingress"
    Exemple de fichier ```customer-value.yaml``` : 
    ```
    imagePullSecrets:
      - name: canopsisregistry

    ingress:
      enabled: true
      hosts:
        - host: canopsis.k8s.lan
          paths:
            - path: /
              pathType: Prefix
      tls:
        - secretName: canopsis-front-tls-secret
          hosts:
            - canopsis.k8s.lan

    mongodb:
      enabled: true
    rabbitmq:
      enabled: true
    redis:
      enabled: true
    timescaledb:
      enabled: true
    ```

=== "Prod - sans ingress"
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
      backendUrls:
        mongodb: mongodb://cpsmongo:canopsis@mongodb:27017/canopsis
        rabbitmq: amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis
        redis: redis://:canopsis@redis:6379/0
        timescaledb: postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis

      ```
    En production, il est recommandé d'exécuter les services backend de Canopsis (bases de données, RabbitMQ) sur des serveurs dédiés, et non dans des conteneurs. Il faut donc veillez à définir les URL backend appropriées comme ci-dessus. 

=== "Prod - avec ingress"
    Exemple de fichier ```customer-value.yaml``` : 
      ```
      imagePullSecrets:
        - name: canopsisregistry

      ingress:
        enabled: true
        hosts:
          - host: canopsis.k8s.lan
            paths:
              - path: /
                pathType: Prefix
        tls:
          - secretName: canopsis-front-tls-secret
            hosts:
              - canopsis.k8s.lan

      mongodb:
        enabled: true
      rabbitmq:
        enabled: true
      redis:
        enabled: true
      timescaledb:
        enabled: true
      backendUrls:
        mongodb: mongodb://cpsmongo:canopsis@mongodb:27017/canopsis
        rabbitmq: amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis
        redis: redis://:canopsis@redis:6379/0
        timescaledb: postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis

      ```
    En production, il est recommandé d'exécuter les services backend de Canopsis (bases de données, RabbitMQ) sur des serveurs dédiés, et non dans des conteneurs. Il faut donc veillez à définir les URL backend appropriées comme ci-dessus. 


!!! Remarque
    Vous pouvez également remplacer tout autre paramètre activé dans le fichier des valeurs du chart, comme indiqué dans le [README](https://git.canopsis.net/helm/charts/canopsis-pro/-/tree/develop?ref_type=heads) du chart. Ce qui précède est le minimum nécessaire pour obtenir un laboratoire Helm Canopsis Pro entièrement fonctionnel lorsque vous souhaitez tester l'ensemble dans Helm/K8S (services backend inclus).

    **Si vous n'avez pas accès au README, merci de vous rapprocher de votre référent.**

## Déploiement

Choisissez un nom de version correspondant à votre instance Canopsis Pro (cano0, cano1, canopsislab, ...), ici nous utiliserons "canopsis-lab". Ensuite, vous devrez l'utiliser de manière appropriée dans toutes les commandes Helm pour cette instance. 

Les commandes d'exemple ci-dessous utiliseront la variable ${RELEASE_NAME} à cette fin.
Le nom de la version est à votre choix ; il permet plusieurs déploiements du même chart pour différentes instances sur le même cluster Kubernetes.

!!! Information
    Les noms de version doivent ressembler à des noms DNS :

    - uniquement des caractères alphanumériques en minuscules ou '-';
    - doivent commencer et se terminer par un caractère alphanumérique ;
    - longueur maximale de 53 caractères.

Définir le nom de votre instance : 
```sh
export RELEASE_NAME="canopsis-lab"
```

=== "Initier le déploiement depuis le repo Helm" 
    ```sh
    helm install ${RELEASE_NAME} canopsis/canopsis-pro -f customer-value.yaml
    ```

=== "Initier le déploiement depuis les sources" 
    ```sh
    helm install ${RELEASE_NAME} canopsis-pro -f customer-value.yaml
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
```
kubectl delete pvc --all
```