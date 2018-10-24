## SNMP

Utilisation de SNMP avec Canopsis.

Nécessite l’installation de la brique `CAT`.

L’accès se fait en direct par l’url suivante : `http://<canopsis_addr:port>/en/static/canopsis/index.html#/userview/view.snmprule`

### Pré-requis

Sur debian, ajouter les contrib et non-free, puis
```bash
apt-get install snmp snmp-mibs-downloader
```
(déjà fait en installant CAT)

Pour ajouter des MIBS (comme Nagios), les copier dans le dossier  `~/.snmp/mibs` ou `/usr/share/mibs/ietf/`.

### Installation et utilisation du connecteur

```bash
apt-get install python-pip

git clone https://git.canopsis.net/cat/connector-snmp2canopsis.git

cd connector-snmp2canopsis

pip install -r requirements.txt

python2 setup.py install

SNMP_DEBUG=1 snmp2canopsis -c snmp2canopsis.conf
```

`SNMP_DEBUG=1` permet d'activer le mode verbeux du connecteur.

### Exemples de trap

Trap standard de test:
```bash
/usr/bin/snmptrap -v 2c -c public localhost '' NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456
```

Exemple de trap avec une mib specifique (Nagios):
```bash
/usr/bin/snmptrap -v 2c -c public localhost '' NAGIOS-NOTIFY-MIB::nSvcEvent nHostname s "uncomposant" nSvcDesc s "uneresource" nSvcStateID i 0
```

### Exemple d'événement envoyé

On envoi un trap de test:
```bash
/usr/bin/snmptrap -v 2c -c public localhost '' NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456
```

Et voila l'évènement généré par le connecteur, et envoyé vers Canopsis:
```python
{
    'component': '127.0.0.1',
    'connector': 'snmp',
    'connector_name': 'snmp2canopsis',
    'event_type': 'trap',
    'snmp_trap_oid': '1.3.6.1.4.1.8072.2.3.0.1',
    'snmp_vars': {
        '1.3.6.1.2.1.1.3.0': '9308932',
        '1.3.6.1.4.1.8072.2.3.2.1': '123456',
        '1.3.6.1.6.3.1.1.4.1.0': '1.3.6.1.4.1.8072.2.3.0.1'
    },
    'snmp_version': '2c',
    'source_type': 'component',
    'state': 3,
    'state_type': 1,
    'timestamp': 1521538609.3371
}
```
