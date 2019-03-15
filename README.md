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



# Déploiement du cluster k8s

⚠️❗Actuellement, le volume docker-mongo mappe le répertoire **/Datas/Test-kubernetes/docker/mongo** ❗⚠️

Il faut donc penser à modifier le fichier **deploy-cano.yaml** pour pointer vers le chemin complet du dossier **docker/mongo**

```vim
108    - name: docker-mongo         
109       hostPath:                                                                
110           path: /Datas/Test-kubernetes/docker/mongo
```

Une fois la modification effectuée, on peut déployer le cluster:

```bash
microk8s.kubectl create -f deploy-cano.yaml
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

