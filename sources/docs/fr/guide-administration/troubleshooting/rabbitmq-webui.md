# UI RabbitMQ

**TODO (DWU) :** aller plus loin avec DWU.
**TODO (MG) :** Vérification de l'exactitude de l'existant.

Dans le contexte d'une utilisation de Canopsis, RabbitMQ va vous servir à avoir une vision globale sur le bon fonctionnement de vos moteurs.
Vous pourrez y voir plusieurs informations utiles mais nous allons nous concentrer dans cette documentation à la section "Queues" qui nous montre le nombre de messages
en cours de traitement sur les diférents moteurs de Canopsis.

**Rappel :**
L'interface de RabbitMQ est accessible via l'URL ```http://localhost:15672/```

## Queues

Rendez-vous ici :

![img1](img/section_queues.png)

Vous y retrouverez un tableau comme celui ci :

![img2](img/tab1.png)

La première colonne *Overview* vous présente plusieurs informations tel que :

- **Virtual host :** Nom de la machine sur laquelle le moteur est présent.  
- **Name :** Nom du moteur.  
- **Feature :** Montre si l'architecture est en HA ou non.  
- **State :** Etat du moteur, peut être *running* ou *idle* (fontionnement dégradé).  

La seconde, *Messages*, vous présente :

- **Ready :** Nombre de messages près à être ack.
- **Unacked :** Nombre de messgaes qui ne sont pas encore ack.
- **Total :** Nombre de message total.

La troisème, *Message rates* permet d'avoir une idée sur les performances du moteur à gerer les files arrivantes. Trois stats permettent de juger l'efficacité en messages/secondes :

- **Incomming :** Nombre de messgae arrivants dans le moteur.
- **Deliver / get :** 
- **ack :**

Le but étant d'avoir une section "Messages" remplie de 0. Si ce n'est pas le cas, cela veut dire qu'un des moteurs de Canopsis est dans un état dégradé et n'assure plus sa gestion de files.  
Dans ce cas plusieurs piste de résolutions sont possibles : 

**TODO**

## Policy

Afin d’éviter de remplir inutilement les queues de RabbitMQ, il est possible de mettre en place une policy.

La procédure est la suivante :

![img](img/rabbitmq_policy.png)

Ensuite, vous devez voir apparaître votre policy sur les queues dans l’onglet Queues.


## Aller plus loins

[Cette documentation](https://www.cloudamqp.com/blog/2015-05-27-part3-rabbitmq-for-beginners_the-management-interface.html#overview) peut vous permettre d'avoir plus de détails sur le foncitonnement général de l'UI de Rabitmq.