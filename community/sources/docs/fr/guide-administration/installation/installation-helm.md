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

___En option___ Ajouter le repo en version devel :
```sh
helm repo add --username oauth2 --password $TOKEN canopsis-devel https://git.canopsis.net/api/v4/projects/603/packages/helm/devel

```

Mettre à jour les repos
```sh
helm repo update
```

## Récupérer le chart en utilisant les sources Helm

Cloner le repo en utilisant le protocole **SSH** :
```sh
git clone git@git.canopsis.net:helm/charts/canopsis-pro.git
```
**ou**

Cloner le repo en utilisant le protocole **HTTPS** :
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

## Surcharge des valeurs du fichier Helm

Créer un fichier ```support-value.yaml``` avec au moins le contenu suivant : 
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

!!! Remarque
    Vous pouvez également remplacer tout autre paramètre activé dans le fichier des valeurs du chart, comme indiqué dans le README du chart. Ce qui précède est le minimum nécessaire pour obtenir un laboratoire Helm Canopsis Pro entièrement fonctionnel lorsque vous souhaitez tester l'ensemble dans Helm/K8S (services backend inclus).

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
```
export RELEASE_NAME="canopsis-lab"
```

Initier le déploiement depuis le repo Helm : 
```
helm install ${RELEASE_NAME} canopsis/canopsis-pro -f support-values.yaml
```
**ou**

Initier le déploiement depuis les sources : 
```
helm install ${RELEASE_NAME} canopsis-pro -f support-values.yaml
```
  
Superviser le déploiement :
```
watch kubectl get pod
```

### Accéder à l'interface Web :
```
kubectl port-forward svc/${RELEASE_NAME}-nginx 8443:443
```

!!! Information
    Cette commande ouvrira le port 8443 en local et redirigera les connexions vers le port 443 du service Nginx

### Désinstaller Canopsis : 
```
helm uninstall ${RELEASE_NAME}
```

### Supprimer les volumes persistants
```
kubectl delete pvc --all
```