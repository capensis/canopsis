# Présentation

Les pbehaviors (pour periodical behaviors) sont des évènements de calendrier
récurrents ou non, qui permettent de mettre en pause la surveillance d'une
alarme pendant une période donnée (« downtime ») pour des maintenances ou des
astreintes par exemple.

## API
Dans cette section nous allons voir en détail leur gestion. Pour cela
[rendez-vous ici](caracterisation.md) pour découvrir leur caractérisation des PBehavior et [ICI](Utilisation.md) pour savoir comment les utiliser.

## Fichier de configuration
Le fichier de configuration des `pbehaviors`
(`/opt/canopsis/etc/pbehavior/manager.conf`) ne contient qu'un champ
`default_timezone` utilisé pour définir la timezone par défaut à utiliser
dans le calcul des dates d'exécution d'un `pbehavior`.
