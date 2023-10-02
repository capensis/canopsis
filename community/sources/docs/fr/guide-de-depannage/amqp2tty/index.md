# amqp2tty - Analyse temps réel des flux issus des connecteurs ou des relais AMQP

La commande `amqp2tty` permet de se connecter en ligne de commande sur l'exchange `canopsis.events` et ainsi d'afficher les évènements bruts qui circulent.

!!! attention
    Sur un même environnement Canopsis, **lancez seulement une instance d'`amqp2tty` à la fois.**
    Une fois votre investigation terminée, arrêtez le programme.

    En effet, si plusieurs `amqp2tty` sont connectés au bus de messages, chaque instance ne verra qu'une partie du flux de messages.
    C'est ce qui peut se passer si plusieurs utilisateurs interviennent en parallèle sur la plate-forme ou lorsqu'un collègue en laisse tourner un exemplaire dans une session en arrière-plan.

    Ceci est lié au principe du programme `amqp2tty` et au fonctionnement normal de RabbitMQ lorsque plusieurs *consumers* sont connectés à une même *queue*.

    Ce point est toutefois sans impact sur les évènements véritablement traités par les moteurs Canopsis.

## Depuis un environnement paquets

Voici un exemple d'utilisation de la commande, qui cherche des évènements en provenance de Centreon. Elle doit être exécutée depuis un nœud Canopsis :

```sh
set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
/opt/canopsis/bin/amqp2tty | grep -i centreon
```

et son résultat :
```json
{
    "connector": "centreon",
    "connector_name": "centreoninstance1",
    "event_type": "check",
    "source_type": "component",
    "component": "test_composant",
    "address": "x.x.x.x",
    "output": "(Process Timeout)",
    "state": 1,
    "state_type": 1,
    "scheduled": true,
    "check_type": 0,
    "current_attempt": 1,
    "max_attempts": 5,
    "execution_time": 9.544648,
    "latency": 0.597,
    "command_name": "/usr/lib64/nagios/plugins/check_icmp -H x.x.x.x -w 3000.0,80% -c 5000.0,100% -p 1",
    "component_alias": "alias composant",
    "hostgroups": ["HG1"],
    "timestamp": 1528812075
}
```

## Depuis un environnement Docker Compose

Dans un environnement Docker Compose, il suffit de lancer le conteneur `amqp2tty`
dans le réseau Docker de Canopsis et de lui indiquer l'URL de RabbitMQ :

=== "Docker Compose Community"

    ```sh
    docker run --rm -e CPS_AMQP_URL=amqp://cpsrabbit:canopsis@rabbitmq/canopsis \
    	--network=canopsis-community_default \
    	docker.canopsis.net/docker/community/amqp2tty:<VERSION CANOPSIS>
    ```

=== "Docker Compose Pro"

    ```sh
    docker run --rm -e CPS_AMQP_URL=amqp://cpsrabbit:canopsis@rabbitmq/canopsis \
    	--network=canopsis-pro_default \
    	docker.canopsis.net/docker/community/amqp2tty:<VERSION CANOPSIS>
    ```


