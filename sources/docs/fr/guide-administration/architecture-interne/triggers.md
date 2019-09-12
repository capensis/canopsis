# Triggers (Go)

Dans Canopsis, le traitement des événements par le moteur [axe](../moteurs/moteur-axe) peut déclencher des actions qui correspondent à des `triggers`.

Ces `triggers` peuvent servir comme point de déclenchement pour les [actions](../moteurs/moteur-action.md) et les [webhooks](../moteurs/moteur-axe-webhooks.md).

Les triggers possibles sont : `"ack"`, `"ackremove"`, `"assocticket"`, `"cancel"`, `"changestate"`, `"comment"`, `"create"`, `"declareticket"`, `"done"`, `"resolve"`, `"snooze"`, `"statedec"`, `"stateinc"`, `"uncancel"`, et `"unsnooze"`.

| Nom               | Description                                              |
|:----------------- |:-------------------------------------------------------- |
| `"ack"`           | Acquittement d'une alerte                                |
| `"ackremove"`     | Suppression de l'acquittement                            |
| `"assocticket"`   | Association d'un ticket à l'alarme                       |
| `"cancel"`        | Annulation de l'évènement                                |
| `"changestate"`   | Modification et verrouillage de la criticité de l'alarme |
| `"comment"`       | Envoi d'un commentaire                                   |
| `"create"`        | Création de l'évènement                                  |
| `"declareticket"` | Déclaration d'un ticket à l'alarme                       |
| `"done"`          | Fin de l'alarme                                          |
| `"resolve"`       | Résolution de l'alarme                                   |
| `"snooze"`        | Report de l'alarme                                       |
| `"statedec"`      | Diminution de la criticité de l'alarme                   |
| `"stateinc"`      | Augmentation de la criticité de l'alarme                 |
| `"uncancel"`      | Retablissement de l'alarme                               |
| `"unsnooze"`      | Fin du report de l'alarme                                |
