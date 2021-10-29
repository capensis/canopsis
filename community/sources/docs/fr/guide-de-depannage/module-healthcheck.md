# Module Healthcheck

Le module `Healthcheck` inclus dans Canopsis Edition Pro permet de connaitre l'état de chaque composant technique de Canopsis.  
Les composants éligibles sont : 

* [API Canopsis](#api-canopsis)
* [Cache Redis](#cache-redis)
* [Base de données MongoDB](#base-de-donnees-mongodb)
* [Bus RabbitMQ](#bus-rabbitmq)
* [Moteurs Canopsis](#moteurs-canopsis)

Le menu `Administration->Bilan de santé` vous permet  d'accéder aux résultats du module Healthcheck.

![module-healthcheck-general](img/module-healthcheck-general.png)

Chaque pastille représente l'état du composant indiqué et peut être verte, rouge, ou grise.

![module-healthcheck-general2](img/module-healthcheck-general2.png)

Par ailleurs, l'état global ou l'état de chaque composant est mis à disposition dans le bandeau de l'application.

![module-healthcheck-general3](img/module-healthcheck-general3.png)
![module-healthcheck-general4](img/module-healthcheck-general4.png)

!!! Note
    Le contrôle des composants est effectué une fois par minute.

## API Canopsis

Lorsque l'API de Canopsis est indisponible, vous ne pouvez plus accéder à l'interface graphique.  
La réponse technique lors de [l'appel à l'API](../../guide-developpement/swagger/#/healthcheck/healthcheck-get) sera de la forme

```html
<html>
<head><title>502 Bad Gateway</title></head>
<body>
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

Il s'agit en réalité de la réponse renvoyée par le reverse proxy `nginx`.

**Remédiation associée**

=== "Paquets CentOS 7"

    ```sh
    systemctl restart canopsis-service@canopsis-api.service
    ```

=== "Docker Compose"

    ```sh
    docker-compose restart api
    ```

Si cela ne permet pas de rétablir le service, une analyse de logs est nécessaire.

## Cache Redis

Lorsque le système de cache `redis` est indisponible, la pastille associée passe au rouge.  
Canopsis ne peut alors pas fonctionner correctement.

[La réponse renvoyée par l'API Canopsis](../../guide-developpement/swagger/#/healthcheck/healthcheck-get) vous précisera alors le dysfonctionnement de Redis

**Remédiation associée**

=== "Paquets CentOS 7"

    ```sh
    systemctl restart redis.service
    ```

=== "Docker Compose"

    ```sh
    docker-compose restart redis
    ```

Si cela ne permet pas de rétablir le service, une analyse de logs est nécessaire.


## Base de données MongoDB

Lorsque la base de données `mongoDB` est indisponible, la pastille associée passe au rouge.  
Canopsis ne peut alors pas fonctionner correctement.

[La réponse renvoyée par l'API Canopsis](../../guide-developpement/swagger/#/healthcheck/healthcheck-get) vous précisera alors le dysfonctionnement de MongoDB.

**Remédiation associée**

=== "Paquets CentOS 7"

    ```sh
    systemctl restart mongod.service
    ```

=== "Docker Compose"

    ```sh
    docker-compose restart mongodb
    ```

Si cela ne permet pas de rétablir le service, une analyse de logs est nécessaire.

## Bus RabbitMQ

Lorsque la base de données `rabbitmq` est indisponible, la pastille associée passe au rouge.  
Canopsis ne peut alors pas fonctionner correctement.

[La réponse renvoyée par l'API Canopsis](../../guide-developpement/swagger/#/healthcheck/healthcheck-get) vous précisera alors le dysfonctionnement de RabbitMQ.

**Remédiation associée**

=== "Paquets CentOS 7"

    ```sh
    systemctl restart rabbitmq.service
    ```

=== "Docker Compose"

    ```sh
    docker-compose restart rabbitmq
    ```

Si cela ne permet pas de rétablir le service, une analyse de logs est nécessaire.

## Moteurs Canopsis

Le contrôle de bon fonctionnement des moteurs comportent plusieurs volets.  

### L'enchainement des moteurs

La transmission d'événements d'un moteur à un moteur respecte un [ordre précis](../../guide-administration/moteurs/schema-enchainement-moteurs). 
Lorsque l'ordre est respecté, le schéma d'enchainement reste au vert.

![module-healthcheck-engines1](img/module-healthcheck-engines1.png)

Si cet ordre n'est pas celui constaté sur l'environnement contrôlé, alors le module de healthcheck réagit.

![module-healthcheck-engines2](img/module-healthcheck-engines2.png)

Dans le cas présent, on constate que le moteur `engine-correlation` n'est alimenté par aucun autre moteur.

**Remédiation associée**

Vous devez vous assurer que les options de lancement des moteurs sont correctement positionnées.


### L'empilement de messages dans les files d'attente

Chaque moteur de Canopsis est associé à une file d'attente de messages.  
Si le nombre de messages dépasse une limite, on peut en déduire que le moteur n'arrive pas à consommer les messages dans les temps.  
Il peut alors s'agir d'un problème de performance dégradée ou d'un potentiel plantage du moteur qui n'arriverait plus du tout à consommer les messages.  

Le paramètre définissant la limite est un paramètre global à toutes les files d'attentes de moteurs. Il est défini dans l'onglet `Paramètres` de la page de Bilan de santé.

![module-healthcheck-engines3](img/module-healthcheck-engines3.png)
![module-healthcheck-engines4](img/module-healthcheck-engines4.png)

**Remédiation associée**

Si le moteur en défaut continue à dépiler des messages dans sa file, ce signifie que la plate-forme n'est pas suffisament dimensionnée pour accepter ce trafic.  
Les pistes à envisager dans ce cas sont : 

* Les requêtes non ou mal indexées dans la base de données MongoDB (Collscan à rechercher dans les logs)
* Les ressources locales allouées sont consommées.
* Envisager une instance supplémentaire du moteur

### Le nombre d'instances minimal ou optimal d'un moteur n'est pas respecté

Le nombre d'instances minimal et optimal d'un moteur peut être défini dans le menu `Paramètres` du Bilan de santé.  
Si le nombre minimal d'instances n'est pas respecté, alors un avertissement est montré par le module.

Par exemple, avec la configuration suivante

![module-healthcheck-engines5](img/module-healthcheck-engines5.png)

Lorsque deux instances sont en cours d'exécution, le module indique 

![module-healthcheck-engines7](img/module-healthcheck-engines7.png)

Si une seule instance est en cours d'exécution, le module indique 

![module-healthcheck-engines6](img/module-healthcheck-engines6.png)

Le nombre d'instances optimal n'est pas exploitée dans l'interface graphique. Il est à exploiter au travers de l'API.
