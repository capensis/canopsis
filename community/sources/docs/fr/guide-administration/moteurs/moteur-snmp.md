# Moteur `snmp` (Python, Pro)

!!! info
    Ce moteur n'est disponible que dans l'édition Pro de Canopsis.

Le moteur `snmp` permet de traiter les traps SNMP récupérés par le connecteur `snmp2canopsis` selon des règles prédéfinies par l'utilisateur.

## Fonctionnement

Ce schéma présente le cycle de vie d'un trap SNMP depuis son émission jusqu'à sa conversion en alarmes Canopsis.

![img1](img/Cycle_vie_trap_snmp.png)

Comme observé sur le schéma de cycle de vie, les traps SNMP bruts sont traduits par un moteur grâce à un jeu de règles à définir.

Pour rappel, le résultat de la traduction doit se concrétiser par la génération d'un message compréhensible par Canopsis.

Ce message doit comporter au minimum les informations suivantes :

*  `connector`
*  `connector_name`
*  `component`
*  `resource`
*  `state`
*  `output`

Le principal objectif est donc de déduire ces attributs à partir du tableau `snmp_vars` présent dans les traps bruts.

## Mise en place

### Activation du moteur

=== "Docker Compose, Canopsis ≥ 22.10"

    Lancez la *stack* complémentaire SNMP à l'aide d'une commande de type :

    ```console
    $ docker compose up -d
    ```

    Dans cette *stack*, le service Compose `snmp` correspond au moteur `snmp`.

=== "RPM, Canopsis 4.x uniquement"

    Sur le nœud des moteurs Canopsis :

    ```console
    # systemctl enable canopsis-engine-cat@snmp
    # systemctl start canopsis-engine-cat@snmp
    ```

### Activation du service SNMP dans l'interface web

=== "Canopsis ≥ 22.10"

    Assurez-vous que le service `api`, de la stack Docker Compose standard
    Canopsis Pro, dispose en variable d'environnement de l'adresse de votre
    service `oldapi`.

    Dans nos environnements de référence, cela se définit dans le fichier
    `compose.env` :

    ```bash
    CPS_OLD_API_URL=http://oldapi:8081
    ```

    Si vous avez dû modifier quelque chose, redéployez le service `api` :

    ```console
    $ docker compose up -d api
    ```

=== "Canopsis < 22.10"

    À la fin du fichier `/opt/canopsis/etc/oldapi.conf` (ou équivalent Docker), assurez-vous de la présence de la ligne suivante :

    ```ini
    canopsis_cat.webcore.services.snmprule = 1
    ```

    et redémarrer le service `oldapi`

    ```console
    # systemctl restart canopsis-service@canopsis-oldapi
    ```

    (ou équivalent en service Docker Compose)

### Traduction des traps

Pour créer des règles de transformations, rendez-vous sur l'interface
d'exploitation dédiée, par le menu « Exploitation » > « Règles SNMP ».

![Menu exploitation](img/menu_exploitation_snmprules.png)

!!! note
    L'accès à cette page est régi par le droit `models_exploitation_snmpRule` de type CRUD.
    Veillez à octroyer les permissions dans la matrice des droits ![Droit SNMPRULE](img/droit_snmprule.png)

Une règle de transformation consiste à convertir des `OID` en valeurs
compréhensibles et associer les attributs nécessaires à un évènement Canopsis
de type `check`.

Dans l'exemple du connecteur [`snmp2canopsis`](../../interconnexions/Supervision/SNMPtrap.md), nous souhaitons obtenir le message suivant :

```json
{
  "connector" : "snmp",
  "connector_name" : "snmp",
  "component" : "Equipement Impacte",
  "resource" : "Ressource Impactee",
  "output" : "Message de sortie du trap SNMP",
  "state" : 3
}
```

Pour cela, nous devons :

* Envoyer les MIB Nagios dans Canopsis
* Créer une règle de transformation
* Tester un trap et constater les résultats

### Envoi des MIB

Le paquet `snmp-mibs-downloader` peut être nécessaire sur la machine qui
exécute `oldapi`. Ce paquet embarque une bibliothèque de MIB et permet, au
besoin, de télécharger celles manquantes depuis le web.

Lors de l'upload des MIB, Canopsis concatène les fichiers uploadés dans l'ordre
dans lequel il les reçoit. Il faut donc être vigilant sur ce point.

Cela peut en plus dépendre du navigateur web utilisé. Par exemple, Firefox
upload les fichiers dans l'ordre dans lequel ils ont été sélectionnés dans la
fenêtre de sélection de fichiers. Chrome, quant à lui, upload les fichiers
sélectionnés dans l'ordre alphabétique.

Par exemple, si le fichier `nagios-root.mib` doit être traité avant le fichier
`NAGIOS-NOTIFY-MIB`, vous devrez soit les uploader dans cet ordre soit les
renommer respectivement en `NAGIOS1-ROOT-MIB` et `NAGIOS2-ROOT-MIB`.

!!! attention
    Lorsque plusieurs fichiers de MIB vont ensemble, comme dans cet exemple,
    il faut absolument les envoyer à Canopsis **au cours de la même opération
    d'upload**, tout en veillant à l'ordre de prise en compte, comme expliqué
    juste avant.

    Un envoi dans le bon ordre mais en deux opérations mènera à une
    interprétation incomplète et donc erronnée de la MIB. Les règles SNMP que
    vous saisiriez ensuite ne fonctionneraient pas.

On sélectionne les fichiers :

![img2](img/scenario_e1.png)

![img3](img/scenario_e2.png)

On vérifie que le traducteur a bien trouvé des objets de type « notification ».

![img4](img/scenario_e3.png)

**Création de la règle**

![img5](img/scenario_e4.png)

### Vérification

On exécute à nouveau l'envoi du trap SNMP :

```sh
/usr/bin/snmptrap -v 2c -c public IP_RECEPTEUR_SNMP '' NAGIOS-NOTIFY-MIB::nSvcEvent \
  nHostname s "Equipement Impacte" nSvcDesc s "Ressource Impactee" \
  nSvcStateID i 3 nSvcOutput s "Message de sortie du trap SNMP"
```

On contrôle le bac à alarmes :

![img6](img/scenario_e5.png)
