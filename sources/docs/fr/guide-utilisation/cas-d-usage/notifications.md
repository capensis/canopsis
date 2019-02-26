# Notifications vers un outil tiers

Canopsis est capable de *réagir* en fonction de critères objectifs de notifier un outil tiers en appelant une API HTTP.  
Ce guide vous propose d'intergir avec Mattermost et avec Logstash.  

!!! Warning
    Dans tous les cas, ces possibilités ne sont offertes que par l'utilisation des moteurs GO dans l'édition CAT de Canopsis.  

Une option du moteur **axe** va vous permettre de prendre en charge cette fonctionnalité.

````
-postProcessorsDirectory /plugins/axepostprocessor/
````

La documentation complète est disponible [ici](../../../guide-administration/webhooks/)


## Mattermost

Nous partons du principe que vous possédez une URL Mattermost valide pour publier du contenu.  

Prenons le cas d'usage suivant : 

!!! note ""
    Je souhaite notifier Mattermost depuis Canopsis dans les cas suivants :
    Création d'alarme, Mise à jour d'état, Ack d'un utilisateur

Pour cela, RDV sur le menu des **webhooks**

![Menu Webhooks](./img/notification_mattermost_menu.png "Menu Webhooks")  


Vous devez créer un règle comme suit : 

**Choix des triggers**

![Choix triggers](./img/notification_mattermost_choix_trigger.png "Choix des triggers")  

**Le pattern d'événements sur lesquels on applique la règle**

![Pattern d'événements](./img/notification_mattermost_edit_pattern.png "Pattern d'événements")  

**La requête HTTP à exécuter**

![Requête HTTP](./img/notification_mattermost_request.png "Requête HTTP")  

Pour vérifier le résultat, nous considérons l'événement suivant :  

````json
{
  "resource": "ressource-doc3", 
  "event_type": "check", 
  "component": "composant-doc3", 
  "connector": "cas-d-usage", 
  "source_type": "resource", 
  "state": 2, 
  "connector_name": "cas-usage-notification-mattermost", 
  "output": "doc cas d'usage"
}

````

Une fois publié, vous pourrez consulter votre canal Mattermost

![Resultats](./img/notification_mattermost_resultat.png "Résultats Mattermost")  
