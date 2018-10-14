# AMQP

Connection et publication sur le bus AMQP.

## AmqpConnection

```python
from canopsis.common.amqp import AmqpConnection

ac = AmqpConnection('amqp://[user:pass]@host:port/vhost')
ac.connect()
```

Une fois connecté vous avez accès aux attributs `connection` et `channel`. C’est avec `channel` que vous allez traiter si vous voulez utiliser les objets pika en direct.

```python
ac.disconnect()
```

Vous pouvez aussi utiliser le mot-clef `with` :

```python
with AmqpConnection('url...') as ac:
    ac.channel...
```

## AmqpPublisher

Un objet rudimentaire permettant d’envoyer un event très simplement et de façon bloquante.

 Easy to use synchronous AMQP publisher.

```python
from canopsis.common.amqp import AmqpPublisher

url = 'amqp://cpsrabbit:canopsis@localhost/canopsis'

logger = logging.getLogger("...")

evt = {...}
with AmqpConnection(url) as apc:
    pub = AmqpPublisher(apc, logger)
    pub.canopsis_event(evt, 'canopsis.events')
```

La méthode `canopsis_event` enverra un event à partir du dictionnaire en créant la `routing_key` nécessaire via `canopsis.event.get_routingkey`.

Vous disposez de `json_document` pour envoyer un document JSON brut à partir d’un dictionnaire Python.
