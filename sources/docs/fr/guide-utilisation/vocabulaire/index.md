# Guide Utilisateur

TODO : Compléter les définitions.  

## Section Vocabulaire

**Evénement :**  

Un *évènement* est un message arrivant dans Canopsis. Il est formatté en json et provient généralement d'une source externe ou d'un connecteur (email, snmp, etc.).
Lorsqu'un événement arrive il est envoyé vers le bac à événement puis traité, il devient donc un alarme.  

**Alarme :**  

Un *alarme* est le résultat de l'analyse des évènements. Elle historise et résume les changements d'état, les actions utilisateurs (acquittement, mise en pause, etc.).

**Etat :**  

**Statut :**  

**Bagot :**  

Un événement est considéré *Bagot* s'il est passé d'un état d'alerte à un état stable un nombre spécifique de fois sur une période donnée.  

**Météo :**  

La *météo des services* permet d'avoir une vue globale de plusieurs alarmes, par exemple en les regroupant par serveur, par type d'activité, ... C'est au choix, et complètement configurable :)

**Entité :**  

Une *entité* est une abstraction utilisée, entre autre, pour conserver des données statiques. Par exemple, l'entité associée à une alarme va contenir des informations qui ne varient pas d'une alarme à l'autre : emplacement du serveur, lien vers procédure liée à ce type d'alarme...

**Enrichissement :**

L'*enrichissement* est l'action de récupérer des informations supplémentaires dans un évènement pour vernir compléter le contexte (càd l'ensemble des entités).

**Ressource :**  

**Service :**  

**Composant :**  

**Connecteur :**  

**Nom de connecteur :**  
