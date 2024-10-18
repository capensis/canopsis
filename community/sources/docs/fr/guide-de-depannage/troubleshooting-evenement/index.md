# Vérification d'évènements

Voici un scénario de vérification d'évènements avec les différentes commandes associées pour valider l'envoi.

## Communication avec RabbitMQ

Il faut, dans un premier temps, vérifier que la communication entre l'instance concernée et RabbitMQ est établie. Cela se fait grâce à la commande suivante :

```sh
sudo tcpdump -vvv -A -i any -s 0 dst IP_OU_FQDN_CANOPSIS and port 5672
```

Rappel : le port `5672` est le port d'écoute par défaut de RabbitMQ.

Si le trafic est bon de ce côté et que le port reçoit bien des informations, vous pouvez passer l'étape suivante.

Sinon, c'est qu'il y a un problème au niveau de ce dernier que vous pourrez diagnostiquer avec quelques commandes telles que `netstat` ou encore `ps`.

## amqp2tty

Dans un second temps, la vérification va passer par `amqp2tty`, une documentation complète à son sujet est [disponible ici](../amqp2tty/index.md).  

## Vérification du JSON et de son contenu

### Syntaxe JSON 

Il se peut que votre JSON ne soit pas bien formaté. Pensez à vérifier celui-ci à l'aide, par exemple, d'un outil en ligne.

### Contenu de l'évènement

Un évènement est un message arrivant dans Canopsis. Il est formatté en JSON et provient généralement d'une source externe ou d'un connecteur (email, SNMP, etc.).

Ces évènements, formatés en JSON, doivent être composés de plusieurs champs, certains étant obligatoires, d'autres étant optionnels.

Voir aussi la page de [Vocabulaire](../../guide-utilisation/vocabulaire/index.md).
