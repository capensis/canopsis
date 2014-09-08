Interconnexion CLI
==================
Cet onglet permet de définir techniquement la méthode d'interconnexion via la ligne de commande		

Description du type de connecteur	Ce type d'interconnexion permet de publier un événement dans le bus AMQP par l'intermédiaire de l'API REST de Canopsis

|run_manager|

Renseignements sur les sources : Fichier Excel / CSV
----------------------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Scripts Canopsis","Canopsis met à disposition 2 scripts d'émission d'événements compatibles Windows et Unix","https://github.com/capensis/canopsis/wiki"


Matrice des flux
----------------
.. csv-table::
   :header: "Source", "Destination", "Protocole","Ports","Remarques"
   :widths: 15, 20, 15,15,15

	"","Application","Webserver Canopsis","HTTP(s)","Prévoir une ouverture de flux par poller nagios"

.. |run_manager| image:: ../_static/images/connectors/InterconnexionCLI.png
