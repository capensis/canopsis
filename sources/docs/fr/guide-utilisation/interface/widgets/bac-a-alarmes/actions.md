# Les action du bac à alarmes

Lorsqu'un événement arrive il est envoyé vers le bac à événement puis traité, il devient donc un alarme (Une alarme est le résultat de l'analyse des évènements. cf [vocabulaire](../../vocabulaire/index.md)).  

Les différentes actions possibles sur cette alarme sont :  


## Accuser reception

![ack](img/ack.png)  

Deux choix possibles : ACK et Fast ACK

L'un permet de voir les détails généraux de l'événement et de lier un numéro de ticket et d'écrire une note. Il permets d'accuser réception, où d'accuser réception et de reporter un incident.  
Le second lui permet d'accuser réception sans pour autant remplir ces informations.  

## Gérer l'incident

Un menu apparaît donc à la place de l'ACK : 

![menu2](img/menu2.png)  

Il permet de :

- Déclarer un ticket
- Associer un ticket
- Annuler l'alarme
 
Et bien plus encore si vous cliquez sur les trois points :

![menu3](img/menu3.png)  

- Annuler l'ACK (cancel ACK) permet de revenir à la première étape.
- "Snooze alarm", comme son nom l'indique, permet de reporter l'alarme à plus tard. Vous pouvez choisir la durée et le laps de temps (minutes, heures, jours, semaines, mois, années).  
Cette icône ![snooze](img/snooze.png) apparaîtra donc, survolez la avec votre curseur pour avoir des détails sur l'auteur, la date de début et de fin.  
- Change criticity permet de changer l'état de l'évènement. Voici les différents états :  ![state](img/state.png)   
- "List periodical behavior" permet de lister nos différent PBehaviors. Voir la [section dédiée](Les-PBehaviors.md)  