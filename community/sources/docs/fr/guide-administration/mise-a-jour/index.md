# Mise à jour de Canopsis

Cette procédure décrit la mise à jour d'une instance mono-nœud de Canopsis.

Lisez l'ensemble de ce document avant de procéder à toute manipulation.

## Versions concernées par cette procédure

Cette procédure concerne uniquement les mises à jour « régulières » de l'outil, c'est-à-dire celles où le premier chiffre du numéro de version n'a pas changé (3.25.0 vers 3.26.0, 3.26.0 vers 3.26.1…).

Les mises à jour « majeures » de Canopsis (3.48.0 vers 4.0.0) ne peuvent pas et **ne doivent pas** être réalisées à l'aide de ce seul document. Pour ces mises à jour en particulier, le Guide de migration lié à cette nouvelle version majeure doit être suivi.

Consultez le [document des numéros de version de Canopsis](numeros-version-canopsis.md) pour en savoir plus au sujet du *versionnage sémantique* utilisé par Canopsis.

Les environnements n'ayant pas été installés en suivant l'une des [méthodes d'installation officielles de Canopsis](../installation/index.md#methodes-dinstallation-de-canopsis), notamment les environnements de type Docker Swarm ou en paquets multi-nœuds, ne sont pas non plus pris en charge par cette procédure.

## Procédure de mise à jour

Vous devez tout d'abord lire **chacune** des [notes de version](../../index.md#notes-de-version) publiée entre votre version actuelle et celle que vous ciblez.

Par exemple, si vous effectuez une mise à jour de Canopsis 3.38.0 à 3.40.0, vous devez :

*  consulter et appliquer toute procédure donnée dans les notes de version de Canopsis 3.39.0,
*  puis celles de Canopsis 3.39.1,
*  puis celles de Canopsis 3.40.0,
*  puis suivre le reste de cette procédure, selon votre méthode d'installation (paquets ou Docker Compose).

Si vous bénéficiez d'un développement spécifique (modules ou add-ons ayant été spécifiquement développés pour votre installation), assurez-vous de suivre toute procédure complémentaire vous ayant été communiquée.

### Mise à jour en installation par paquets

!!! attention
    Cette mise à jour causera une **interruption de service** de Canopsis et des composants qui lui sont associés, durant son déroulement.
    
    Vous pouvez notamment utiliser la fonctionnalité de [diffusion de messages](../../guide-utilisation/interface/broadcast-messages.md) afin de prévenir vos utilisateurs en amont.

L'ensemble des commandes suivantes doivent être réalisées avec l'utilisateur `root`.

Appliquez la mise à jour des paquets Canopsis :

```sh
yum --disablerepo="*" --enablerepo="canopsis*" update
```

Vous devez ensuite finaliser la mise à jour avec les commandes suivantes, fonction de votre édition de Canopsis (Community ou Pro) :

=== "Canopsis Community"

    ```sh
    su - canopsis -c "canopsinit --canopsis-edition core"
    ```

=== "Canopsis Pro"

    ```sh
    su - canopsis -c "canopsinit --canopsis-edition cat"
    ```

Puis, après avoir pris en compte toute éventuelle remarque des notes de version au sujet du fichier `canopsis.toml`, appliquez les changements de configuration :

```bash
set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
/opt/canopsis/bin/canopsis-reconfigure
``` 

Puis, redémarrez l'ensemble des moteurs Canopsis :

```sh
systemctl daemon-reload
canoctl restart
```

Ne pas oublier d'appliquer toute éventuelle procédure supplémentaire décrite dans chacune des [notes de version](../../index.md#notes-de-version) qui vous concerne.

### Mise à jour en environnement Docker Compose

Après avoir suivi les notes de version, resynchronisez l'ensemble de vos fichiers Docker Compose avec les fichiers de référence correspondant à la version voulue : <https://git.canopsis.net/canopsis/canopsis-community/-/tree/develop/community/docker-compose>.

Puis, exécutez la commande suivante :

```sh
docker-compose up -d
```

Ne pas oublier d'appliquer toute éventuelle procédure supplémentaire décrite dans chacune des [notes de version](../../index.md#notes-de-version) qui vous concerne, y commpris toute éventuelle remarque des notes de version au sujet du fichier `canopsis.toml`.

## Après la mise à jour

Durant le temps de coupure des services Canopsis, RabbitMQ se sera chargé de mettre en attente vos [évènements](../../guide-utilisation/vocabulaire/index.md#evenement). Ils seront alors « dépilés » et traités normalement par les moteurs Canopsis, dès leur redémarrage.

Cette accumulation d'évènements en attente peut, néanmoins, provoquer une latence des traitements, ou une augmentation de la consommation des ressources, en raison du rattrapage à effectuer. Cette incidence reste temporaire. Nous vous conseillons de [surveiller l'interface d'administration de RabbitMQ](../../guide-de-depannage/rabbitmq-webui.md) juste avant, durant et après la mise à jour, afin de mesurer l'état de « retour à la normale » de votre plateforme lors d'une période de maintenance de l'outil.

En revanche, tout appel fait aux API Canopsis durant cette période de maintenance n'aura pas été temporisé et devra donc être renouvelé s'il a échoué.

Une fois que le service est rétabli, vous pouvez vous connecter à nouveau sur l'interface Canopsis pour valider que tout fonctionne correctement.
