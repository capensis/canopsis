# Triggers (Go)

Dans Canopsis, le traitement des [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) par le moteur [`engine-axe`](../moteurs/moteur-axe.md), des actions automatisées par le moteur [`engine-action`](../moteurs/moteur-action.md) et des webhooks par le moteur [`engine-webhook`](../moteurs/moteur-webhook.md) peuvent déclencher des `triggers`.

Ces `triggers` peuvent servir comme point de déclenchement pour les [actions automatisées](../moteurs/moteur-action.md) et les [webhooks](../moteurs/moteur-webhook.md).

Les triggers possibles sont : `"ack"`, `"ackremove"`, `"assocticket"`, `"activate"`, `"cancel"`, `"changestate"`, `"comment"`, `"create"`, `"declareticket"`, `"declareticketwebhook"`, `"done"`,  `"pbhenter"`,  `"pbhleave"`, `"resolve"`, `"snooze"`, `"statedec"`, `"stateinc"`, `"statusdec"`, `"statusinc"`, `"uncancel"`, et `"unsnooze"`.

!!! note
    Les triggers `declareticketwebhook`, `resolve` et `unsnooze` ne correspondent pas à un évènement mais à un traitement interne par Canopsis

| Nom                      | Description                                                                                         | Déclenché par des [évènements](../../guide-utilisation/vocabulaire/index.md#evenement) |
|:------------------------ |:--------------------------------------------------------------------------------------------------- | ---------------------------- |
| `"ack"`                  | Acquittement d'une [alerte](../../guide-utilisation/vocabulaire/index.md#alarme)                    | ✅                           |
| `"ackremove"`            | Suppression de l'acquittement                                                                       | ✅                           |
| `"activate"`            | Passage d'une alerte en mode activation ( soit dès sa création si il n'y a pas de pbehavior ou de snooze en cours, ou soit en sortie de Snooze ou de Pbehavior retardant ainsi le début de l'alerte à prendre en compte dans par exemple un calcul de SLA                                                                       | ❌                           |
| `"assocticket"`          | Association d'un ticket à l'alarme                                                                  | ✅                           |
| `"cancel"`               | Annulation de l'évènement                                                                           | ✅                           |
| `"changestate"`          | Modification et verrouillage de la [criticité](../../guide-utilisation/vocabulaire/index.md#criticité) de l'alarme | ✅                           |
| `"comment"`              | Envoi d'un commentaire                                                                              | ✅                           |
| `"create"`               | Création de l'évènement                                                                             | ✅                           |
| `"declareticket"`        | Action du bac à alarmes de déclaration d'un ticket                                                  | ✅                           |
| `"declareticketwebhook"` | Déclaration d'un ticket à l'alarme par un webhook                                                   | ❌                           |
| `"done"`                 | Fin de l'alarme                                                                                     | ✅                           |
| `"pbhenter"` | Démarrage d'un pbehavior                                                   | ❌                           |
| `"pbhleave"` | Sortie d'un pbehavior                                                   | ❌                           |
| `"resolve"`              | Résolution de l'alarme                                                                              | ❌                           |
| `"snooze"`               | Mise en veille de l'alarme                                                                          | ✅                           |
| `"statedec"`             | Diminution de la [criticité](../../guide-utilisation/vocabulaire/index.md#criticité) de l'alarme    | ✅                           |
| `"stateinc"`             | Augmentation de la [criticité](../../guide-utilisation/vocabulaire/index.md#criticité) de l'alarme  | ✅                           |
| `"statusdec"`            | Diminution du [statut](../../guide-utilisation/vocabulaire/index.md#statut) de l'alarme             | ✅                           |
| `"statusinc"`            | Augmentation du [statut](../../guide-utilisation/vocabulaire/index.md#statut) de l'alarme           | ✅                           |
| `"uncancel"`             | Rétablissement de l'alarme                                                                          | ✅                           |
| `"unsnooze"`             | Sortie de veille de l'alarme                                                                        | ❌                           |

!!! attention
   Le trigger `"declareticketwebhook"` peut dans certains cas générer une boucle infinie si il n'est pas utilisé correctement, il faut donc être attentif lors de son utilisation. 
