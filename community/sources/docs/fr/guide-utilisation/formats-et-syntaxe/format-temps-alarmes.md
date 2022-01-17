# Format des temps des alarmes

Le tableau ci-dessous récapitule les différents champs à caractère temporel [structurant une alarme](../../guide-developpement/base-de-donnees/periodical-alarm.md).

| Champ  | Format | Description |
|--------|--------|-------------|
| `t`    | Unix Timestamp | Date de création de l'événement |
| `v.state.t` | Unix Timestamp | Date du dernier changement d'état |
| `v.status.t` | Unix Timestamp | Date du dernier changement de statut |
| `v.creation_date` | Unix Timestamp | Date de création de l'alarme |
| `v.last_update_date` | Unix Timestamp | Date de la dernière modification de l'alarme (changement de criticité, pose d'un acquittement…). Par défaut égale à `v.last_event_date`. Pour activer la dissociation des deux variables, il est nécessaire de configurer l'option `EnableLastEventDate = true` dans le fichier `canopsis.toml`. |
| `v.last_event_date` | Unix Timestamp | Date du dernier événement reçu pour cette alarme même si cela n'a pas généré de changement |
| `v.snooze_duration` | Secondes | Durée des mises en veille |
| `v.pbh_inactive_duration` | Secondes | Durée des comportements périodiques |
| `v.duration` | Secondes | Durée totale de l'alarme |
| `v.active_duration` | Secondes | Durée réelle de l'alarme : `v.active_duration` = `v.duration` - `v.snooze_duration` - `v.pbh_inactive_duration` |
| `v.current_state_duration` | Secondes | Durée du dernier statut |


