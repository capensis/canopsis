# Description d'une alarme

Une alarme résume l'état d'un élément du SI sur lequel nous avons effectué des actions de l'apparition du problème à sa résolution.  

Voici une liste non exhaustive des différents cas où une alarme peut apparaître :

*  un **évènement de contrôle avec un état non OK**, déclenchant l'alarme  
*  un **évènement de contrôle, avec un état OK**, mettant fin à l'alarme  
*  un évènement de **déclaration de ticket** ou un évènement d'*association de ticket*  
*  un évènement d'**état changeant**  
*  un évènement d'**annulation** d'alarme  
*  un évènement de **restauration** d'alarme  
*  un évènement **snooze** d'alarme  
*  un évènement de **commentaire** d'alarme  
*  un ou plusieurs évènements de **contrôle avec un état distinct non OK**  
*  un ou plusieurs **évènements de reconnaissance**  
*  
Cet ensemble d'évènements s'appelle, dans Canopsis, un *cycle d'alarme* et est associé à une entité contextuelle.

## Etapes d'Alarmes

Le cycle d'alarme ne peut être terminé qu'après un statut défini sur 0 si la période de battement potentiel s'est écoulée.   
À chaque étape de l'alarme , le cycle d'alarme peut transporter une information parmi celles ci-dessous:  
  
*  l'alarme est en cours (on going)  
*  l'alarme sonne (flapping)  
*  l'alarme est furtive (stealthy)  
*  l'alarme a été acquittée  
*  l'alarme a été associée à un ticket  
*  un ticket a été déclaré pour l'alarme  
*  l'alarme a été annulée  
*  l'alarme a été rétablie à partir de son état annulé  
*  l'état d'alarme a augmenté  
*  l'état d'alarme a diminué  
*  l'alarme a été snoozed  
*  le nombre d'étapes a atteint une limite stricte  
*  l'alarme a été commentée  

Chaque étape **DOIT** être historisée dans son cycle d'alarme correspondant. Et une fois l'alarme terminée, le cycle **DOIT** être fermé et archivé.

## Les différents états

Les différents états que peut avoir une alarme sont :  

*  0 - Info
*  1 - Minor
*  2 - Major
*  3 - Critical

Ces états sont **variables** et **visibles** dans le [bac à alarmes](actions.md).  

## Les différents statuts

Les différents statuts que peut avoir une alarme sont :

*  0 - Off
*  1 - On going
*  2 - Stealthy
*  3 - Bagot
*  4 - Cancel

### Off

Un évènement est considéré **Off** s'il est stable. (c'est-à-dire que _Criticity_ est stable à 0).

### On going

Un évènement est considéré **On going** si sa _criticité_ est dans un état d'alerte. (> 0).

### Stealthy

Un évènement est considéré **Stealthy** si sa _criticité_ est passée d'alerte à stable dans un délai spécifié.  
Si la _criticité_ de cet évènement est modifiée à nouveau dans le délai spécifié, il est toujours considéré **Stealthy**.  
Un évènement restera **Stealthy** pendant une durée spécifiée et passera à **Off** si le dernier état était 0, **On Going** s'il s'agissait d'une alerte ou **Bagot** s'il se qualifie en tant que tel.

### Bagot

Un évènement est considéré Bagot s'il est passé d'un état d'alerte à un état stable un nombre spécifique de fois sur une période donnée.

### Cancel

Un évènement est considéré **cancel** si l'utilisateur l'a signalé comme tel à partir de l'interface utilisateur.
Un évènement marqué comme **cancel** changera d'état s'il passe d'un état d'alerte à un état stable.
De plus, l'utilisateur peut spécifier si l'évènement doit changer d'état si sa _criticité_ change dans les différents états d'alerte ou uniquement entre les états d'alerte et les états stables.

## Timeline

La timeline représente tous les changements qui ont été faits sur l'alarme ; c'est à dire ses changments d'état, de statut, d'output, ...  
Cette timeline est limitée en taille : s'il y a trop d'états, un "crop state" sera généré ; c'est un bloc résumant combien de changements ont été effectués.

Dans le cas où le recadrage par paliers ne suffirait pas, un fonction existe empêchant une alarme de devenir trop grosse. elle contrôle simplement le nombre maximum d'étapes qu'une alarme peut avoir.
Si une étape doit être ajoutée alors que la limite est atteinte, elle doit être supprimée et ne peut pas être récupérée.
La seule étape à prendre en compte lorsqu'une alarme a atteint sa limite absolue est une annulation d'alarme.
  
Le nombre limite d'étapes à conserver est configurable. Cette valeur peut être mise à jour à tout moment et les alarmes qui ont été gelées doivent continuer à enregistrer les étapes si cette limite a été étendue.
