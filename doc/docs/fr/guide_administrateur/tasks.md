# Tasks Canopsis

Les `tasks` sont des engines un peu spéciaux qui ne recevront d’évènements que s’ils sont explicitement envoyés via l’event filter.

## Configuration

### Activation de la tâche

Dans le fichier `etc/amqp2engines.conf` il suffit d’ajouter une section vide de la forme :

```ini
[engine:task_TASKNAME]
```

Exemple :

```ini
[engine:task_mail]
```

Avec les unités `systemd` :

```bash
systemctl enable canopsis-engine@task_TASKNAME-task_TASKNAME

# Ou pour une task CAT
systemctl enable canopsis-engine-cat@task_TASKNAME-task_TASKNAME
```

### Création d’une task

Dans les paramètres de Canopsis, dans `Scheduled Jobs` créer une task en choisissant son type. Remplir les formulaires, puis cliquer sur `Next` afin de passer à la suite.

Dans l’onglet `Information` vous pouvez choisir entre une tâche `notification` ou `scheduled`. Pour une tâche devant s’exécuter sur réception d’un évènement, choisissez `notification`.

Pour une tâche devant s’exécuter périodiquement, choisissez `scheduled` puis paramètrez la RRULE (règle de récurrence) dans l’onglet `Schedule`.

Cliquez sur `Save changes`.

### Règle d’event filter

Si votre task est bien de type `notification`, vous pourrez alors aller creer une règle d’event filter :

 * Completez le filtre pour ne *matcher* que les évènements concernés par la tâche
 * Dans `Actions`, sélectionnez `execjob`.
 * Sélectionnez dans la liste la task que vous voulez utiliser pour cet event
 * Cliquez sur `Add action`
 * Cliquez enfin sur `Save changes`

Une minute plus tard les règles seront alors rechargées par l’event filter et les évènements seront alors envoyés aux tasks adéquates.