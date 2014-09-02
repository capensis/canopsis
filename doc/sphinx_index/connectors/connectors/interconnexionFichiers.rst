Interconnexion Fichiers
=======================

Cet onglet permet de définir techniquement la méthode d'interconnexion via des Logs et Trap SNMP
Description du type de connecteur	Ce type d'interconnexion permet de publier un événement dans le bus AMQP par l'intermédiaire d'un connecteur qui réceptionne des fichiers de logs ainsi que des Trap SNMP

http://www.logstash.net

Type d'application	Choisir un type d'application parmi ceux proposés sur le schéma

|run_manager|

Renseignements sur les sources : Logs
--------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Nom des fichiers de Log","Liste des Logs + Emplacement FS",
	"Format des fichiers de Logs","Décrire le format des fichiers de logs à utiliser",
	"Agent de transmission de fichier","Quel agent peut-être utilisé ? rsyslog, nxlog, snare, autre",


Renseignements sur les sources : Traps SNMP
-------------------------------------------

.. |I1| replace:: Nom des fichiers MIB
.. |C1| replace:: MIB contenant les Trap SNMP à intercepter

.. |I2| replace:: Nom des objets de type NOTIFICATION-TYPE ou TRAP-TYPE
.. |C2| replace:: Dans la MIB, quels sont les objets à intercepter ?

.. |I3| replace:: Version SNMP / Communauté SNMP / Auth SNMP
.. |C3| replace:: Informations administratives SNMP

+-----------------------+--------------+---------+
| Item                  | Commentaires | Valeurs |
+=======================+==============+=========+
| |I1|                  |   |C1|       |         |
+-----------------------+--------------+---------+
| |I2|                  |   |C2| |     |    |    |
+-----------------------+--------------+---------+
| |I3|                  |   |C3|       |         |
+-----------------------+--------------+---------+

Matrice des flux
----------------
.. csv-table::
   :header: "Source", "Destination", "Protocole","Ports","Remarques"
   :widths: 15, 20, 15,15,15

	"Connecteur Logstash","Bus Canopsis","AMQP","5672",
	"Connecteur Logstash","Bus Canopsis","HTTP(s)","5672",
	"Application","Connecteur Logstash","TCP","5140",
	"Application","Connecteur Logstash","Trap SNMP (UDP)","162",

.. |run_manager| image:: ../../_static/images/connectors/InterconnecionFlux.png
