Interconnexion FLUX
===================
Cet onglet permet de définir techniquement la méthode d'interconnexion via des mécanismes de flux


		
Description du type de connecteur	Ce type d'interconnexion permet de publier un événement dans le bus AMQP par l'intermédiaire d'un lien direct ou via une API applicative	
		
		
Type d'application	Choisir un type d'application parmi les 4 présents sur le schéma

|run_manager|

*Renseignements sur les sources : Application qui appelle un WS*
----------------------------------------------------------------

.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"API Canopsis", "Documentation de l'API Canopsis - Authentification - Publication d'un événement", "https://github.com/capensis/canopsis/wiki/API-Web"

*Renseignements sur les sources : Application qui propose un WS*
----------------------------------------------------------------

.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"API Applicative","Documentation de l'API nécessaire",""
	"Méthode d'authentification","Décrire le process ainsi que les requêtes d'authentification",""
	"Méthode de sélection des événements","Décrire le process ainsi que les requêtes de sélection des événements sur l'API", ""


*Renseignements sur les sources : Application Custom*
-----------------------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Bindings dans différents langages","Canopsis propose des échantillons de code pour publier des événements dans Canopsis dans les langages suivants :
	PHP
	Perl
	Python
	","https://github.com/capensis/canopsis/wiki/Send-Event-with-PHP 
		https://github.com/capensis/canopsis/wiki/Send-Event-with-Perl 
		https://github.com/capensis/canopsis/wiki/Send-Event-with-Python"

*Renseignements sur les sources : Superviseurs Nagios/Shinken*
--------------------------------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Superviseur","S'agit-il de Nagios ou Shinken ?",
	"Version superviseur",,


Matrice des flux
----------------
.. csv-table::
   :header: "Source", "Destination", "Protocole","Ports","Remarques"
   :widths: 15, 20, 15,15,15

	"Application","Webserver Canopsis","HTTP(s)","80,443",
	"Application","Bus Canopsis","AMQP","5672",
	"Connecteur","Application","HTTP(s)","80,443",
	"Connecteur","Webserver Canopsis","HTTP(s)","80,443"
	"Connecteur","Bus Canopsis","AMQP","5672",

.. |run_manager| image:: ../_static/images/connectors/InterconnecionFlux.png
