## Watcher

Un watcher est un objet ayant un état calculé à partir de l’état d’entités filtrées depuis le context graph.

Le plus mauvais état sera associé au watcher.

Exemple avec trois entités :

 * `Host1/mongodb` : état "critical"
 * `Host1/webserver` : état "warning"
 * `Host1/disk_space` : état "ok"

En supposant que nous créons un Watcher sur ces trois entités, le watcher aura l’état "critical".

### Caractérisation

 * `_id` : `string` ID du watcher
 * `mfilter` : `string` filtre mongo sur les entités du context graph
 * `display_name` : FIXIT: utilisé où ?
 * `enable` : `bool` activer/désactiver le watcher
 * `description` : `string` inutilisé
 * `state_when_all_ack` : `string` inutilisé
 * `downtimes_as_ok` : inutilisé
 * `output_tpl` : `string` inutilisé
 * `state_algorithm` : inutilisé

L’état d’un watcher est recalculé dans les cas suivants :

 * Création d’un PBehavior
 * Modification / suppression d’un PBehavior
 * Changement d’état d’une alarme