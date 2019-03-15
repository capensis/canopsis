# Prérequis

## Installation de microk8s

```bash
sudo snap install microk8s --classic
sudo microk8s.start
```



## Installation du plugin DNS

```bash
microk8s.enable dns
```



## Installation de heml

```bash
sudo snap install heml --classic
heml init
```



# Déploiement du cluster k8s

⚠️❗Par défaut, le volume docker-mongo mappe le répertoire **/Datas/Test-kubernetes/docker/mongo** ❗⚠️

Il faut donc modifier la variable **MONGO_VOLUME** pour pointer vers le chemin complet du dossier **docker/mongo** de votre clone

```vim
vim canopsis/values.yaml

MONGO_VOLUME: /Datas/Test-kubernetes/docker/mongo
```

⚠️❗Par défaut, le PersistentVolume **task-pv-volume** mappe le répertoire **/tmp/mongo1** ❗⚠️

Si vous souhaitez modifier le path du volume persistent, il faut modifier la variable **MONGO_PERSISTENT_VOLUME** dans le fichier **canopsis/values.yaml**



Une fois les modifications effectuées, on peut déployer le cluster:

```bash
helm install --name canopsis ./canopsis
```



# Connexion à l'interface web de Canopsis

⚠️ Avant de pouvoir accéder à l'interface web, il faut s'assurer que le pod **provisionning** est bien à l'état **completed** ⚠️



```bash
microk8s.kubectl get pods
```

La commande doit retourner l'état suivant pour le pod provisionning :

```bash
provisionning                     0/1     Completed          0          15h
```



## Récupération de l'ip du Cluster

```bash
microk8s.config |grep server
```

La commande renvoit une sortie de la forme suivante :

```bash
server: http://adr:192.168.100.131:8080
```



## Récupération du port mappé par le service webserver

```bash
microk8s.kubectl get svc webserver
```

La commande renvoit une sortie de la forme suivante:

```bash
NAME        TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
webserver   NodePort   10.152.183.253   <none>        80:31270/TCP   16h
```



Nous pouvons donc maintenant nous connecter en saisissant l'url http://192.168.100.131:31270/

