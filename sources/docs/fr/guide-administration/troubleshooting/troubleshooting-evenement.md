# Vérification d'événements

Voici un scénario de vérification d'événements avec les différentes commandes associées pour valider l'envoie.

## Communication avec Rabbitmq

Il faut dans un premier temps vérifier que la communication entre l'instance concernée et Rabbitmq est établie, cela ce fait grace à la commande suivante :

```
sudo tcpdump -vvv dst IP_OU_FQDN_CANO and port PORT
```

Si le traffic est bon de ce côté et que le port reçoit bien des informations vous pouvez passer l'étape suivante.
Sinon c'est qu'il y a une problème au niveau de ce dernier que vous pourrez diagnostiquer avec quelques commandes telles que **netstat** ou encore **ps**.

## amqp2tty

Dans un second temps, la vérification va passer par amqp2tty, une documentation complète à son sujet est [disponible ici](amqp2tty.md).  

## Vérification du JSON et de son contenu

### JSON 

Il se peut que votre JSON ne soit pas bien formatté. Pensez à vérifier celui-ci à laide, par exemple, d'un outil en ligne.

### Vérification du contenu de l'évènement

```
Un évènement est un message arrivant dans Canopsis. Il est formatté en json et provient généralement d'une source externe ou d'un connecteur (email, snmp, etc.).
Lorsqu'un événement arrive il est envoyé vers le bac à événement puis traité, il devient donc un alarme.
```
Cf: [Vocabulaire](../../guide-utilisation/vocabulaire/index.md)  

Ces évènements, formaté en Json, doivent-être composés de plusieurs champs obligatoires ou non. 
Un listing détaillé de ces champs est [disponible ici](../../guide-developpement/struct-event.md).  