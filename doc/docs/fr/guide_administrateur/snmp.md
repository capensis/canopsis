## SNMP

Utilisation de SNMP avec Canopsis.

Nécessite l’installation de la brique `CAT`.

### Installation

Sur debian, ajouter les contrib et non-free, puis
```
apt-get install snmp snmp-mibs-downloader
```
(déjà fait en installant CAT)

Pour ajouter des MIBS (comme Nagios), les copier dans le dossier  ~/.snmp/mibs ou /usr/share/mibs/ietf/

### Envoie de traps
Trap de test:
```
/usr/bin/snmptrap -v 2c -c public localhost '' NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456
```

Exemple de trap avec une mib specifique:
```
/usr/bin/snmptrap -v 2c -c public localhost '' NAGIOS-NOTIFY-MIB::nSvcEvent nHostname s "uncomposant" nSvcDesc s "uneresource" nSvcStateID i 0
```

### Installation et utilisation du connecteur

```
apt-get install python-pip

git clone https://git.canopsis.net/cat/connector-snmp2canopsis.git

cd connector-snmp2canopsis

pip install -r requirements.txt

python2 setup.py install

SNMP_DEBUG=1 snmp2canopsis
```

SNMP_DEBUG=1 permet d'activer le mode verbeux du connecteur.
