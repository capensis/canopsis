# Triggers (Go)

Dans Canopsis, le traitement des [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) par le moteur [`engine-axe`](../moteurs/moteur-axe.md), des actions automatisées par le moteur [`engine-action`](../moteurs/moteur-action.md) et des webhooks par le moteur [`engine-webhook`](../moteurs/moteur-webhook.md) peuvent déclencher des `triggers`.

Ces `triggers` peuvent servir comme point de déclenchement pour les [actions automatisées](../moteurs/moteur-action.md) et les [webhooks](../moteurs/moteur-webhook.md).


| Nom                      | Description                                                                                         |
|:------------------------ |:--------------------------------------------------------------------------------------------------- |
|`create`  | Création d'alarme
|`statedec`| Diminution de la criticité
|`changestate`| Changement et verrouillage de la criticité
|`stateinc`| Augmentation de la criticité
|`changestatus`| Changement de statut (flapping, bagot, ...)
|`ack`| Acquittement d'une alarme
|`ackremove`| Suppression de l'acquittement d'une alarme
|`cancel`| Annulation d'une alarme
|`uncancel`| Annulation de l'annulation d'une alarme
|`comment`| Commentaire sur une alarme
|`done`||
|`declareticket`| Déclaration de ticket depuis l'interface graphique
|`declareticketwebhook`| Déclaration de ticket depuis un webhook
|`assocticket`| Association de ticket sur une alarme
|`snooze`| Mise en veille d'une alarme
|`unsnooze`| Sortie de veille d'une alarme
|`pbhenter`| Comportement périodique démarré
|`activate`| Activation d'une alarme
|`resolve`| Résolution d'une alarme
|`pbhleave`| Comportement périodique terminé
|`instructionfail`| Consigne manuelle en erreur
|`autoinstructionfail`| Consigne automatique en erreur
|`instructionjobfail`| Job de remédiation en erreur
|`instructioncomplete`| Consigne manuelle terminée
|`autoinstructioncomplete`| Consigne automatique terminée
