# Erreur de type `ShortStringTooLong`

## Contexte

Cette erreur peut survenir dans plusieurs situations, par exemple :

- lors de l'envoi d'un événement en utilisant l'API `event` ou un connecteur. Dans ce cas l'API renvoie une erreur 500.
- lors d'une action sur une alarme dans le [Bac à alarmes](../../guide-utilisation/interface/widgets/bac-a-alarmes/index.md) (ack, snooze…). Dans ce cas un message d'erreur de type `Something went wrong` s'affiche dans l'interface.

L'erreur est alors visible en consultant les logs de Canopsis. Le détail du message est du type `exceptions.ShortStringTooLong(encoded_value)\nShortStringTooLong:` suivi de la chaîne de caractères qui a provoqué l'erreur.

## Cause

Comme expliqué dans la documentation des [limitations des évènements](../../guide-utilisation/limitations/index.md#limitation-des-evenements), ceci est dû à la longueur de la clé de routage RabbitMQ générée par votre évènement.

## Solution

La clé de routage étant constituée de la façon suivante `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`, essayez de réduire la longueur d'une ou plusieurs de ces valeurs.
