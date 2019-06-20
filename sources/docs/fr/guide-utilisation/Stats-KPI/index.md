# Indicateurs statistiques et KPI

Canopsis fournit des indicateurs statistiques et des indicateurs de performance (KPI).

## Compteurs

| Nom            | Description     |
|----------------|-----------------|
| Alarmes Créées | C'est le nombre d’alarmes créées pendant une période donnée (hors alarmes pendant plage de maintenance). |
| Alarmes Résolues  | Il s'agit du nombre d’alarmes résolues (y compris les alarmes annulées) pendant une période donnée (hors alarmes pendant plage de maintenance). |
| Alarmes Annulées  | Retourne le nombre d’alarmes annulées pendant une période donnée (hors alarmes pendant plage de maintenance).  |
|  Alarmes Acquittées |  C'est le nombre d’alarmes acquittées pendant une période donnée (hors alarmes pendant plage de maintenance). |
| Nombre d’alarmes en cours pendant la période  | Il s'agit du nombre d’alarmes qui étaient actives pendant une période donnée (hors alarmes pendant plage de maintenance).  |
| Nombre d’alarmes actuellement en cours  | Retourne le nombre d’alarmes qui sont actives à un instant donné.  |
| Nombre d’alarmes acquittées actuellement en cours  | C'est le nombre d’alarmes actives et qui ont été acquittées à un instant donné.  |
| Nombre d’alarmes non acquittées actuellement en cours  |  Remonte le nombre d’alarmes actives et qui n’ont pas été acquittées à un instant donné. |

## Taux

| Nom            | Description     |
|----------------|-----------------|
| Taux d’Ack conforme SLA  | Fournit le pourcentage d’alarmes qui ont été acquittées en un temps conforme à un SLA (relatif à une configuration par widget et hors alarmes pendant plage de maintenance)  |
| Taux de résolution conforme SLA  |  Il s'agit du pourcentage d’alarmes qui ont été résolues (y compris annulées) en un temps conforme à un SLA (relatif à une configuration par widget et hors alarmes pendant plage de maintenance). |

## Temps

| Nom            | Description     |
|----------------|-----------------|
|  Proportion du temps dans un état |  Retourne le pourcentage de temps passé par une entité dans un ou plusieurs états (relatif à une configuration par widget et hors alarmes pendant plage de maintenance).   |
|  Temps passé dans un état | C'est le temps passé par une entité dans un ou plusieurs états (relatif à une configuration par widget et hors alarmes pendant plage de maintenance).  |
| Temps moyen entre les pannes (MTBF)  | Temps moyen entre les pannes sur une entité (hors alarmes pendant plage de maintenance)  |


## État

| Nom            | Description     |
|----------------|-----------------|
|  État courant | Affiche l'état actuel d’une entité : 0 pour OK, 1 pour Mineur, 2 pour Majeur et 3 pour Critique (comme dans la météo de service). |
