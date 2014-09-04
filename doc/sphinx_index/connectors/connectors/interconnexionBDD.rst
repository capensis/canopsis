Interconnexion BDD
==================

Cet onglet permet de définir techniquement la méthode d'interconnexion via des mécanismes de sélection de données en base		
		
Description du type de connecteur	Ce type d'interconnexion permet de publier un événement dans le bus AMQP par l'intermédiaire d'un connecteur qui effectue des sélections en base de données.

|run_manager|

Renseignements sur les sources : Fichier Excel / CSV
-----------------------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Nom du fichier","Nom du fichier + emplacement",
	"Format du fichier","Décrire le découpage du fichier",
	"Mise à disposition","Quelles sont les méthodes de mise à disposition du fichier ? Copie, Partage réseau, FTP, Autre",


Renseignements sur les sources : Base de données
------------------------------------------------
.. csv-table::
   :header: "Item", "Commentaires", "Valeurs"
   :widths: 15, 20, 15

	"Technologie","S'agit-il d'une base de données de type : Oracle, SQL Server, MySQL",
	"Version",,
	"Nom de la base de données","Nom de la base de données/instance",
	"Port du listener","Fonction du type de base de données",
	"Compte d'authentification",,
	"Tables concernées",,
	"Exemples de requêtes","L'ensemble des requêtes finales est à positionner dans l'onglet 'mapping' pour chaque attribut",
	"Commentaires",,

Matrice des flux
----------------
.. csv-table::
   :header: "Source", "Destination", "Protocole","Ports","Remarques"
   :widths: 15, 20, 15,15,15

	"Connecteur EDC","Application","Listener de BDD","Fonction de la BDD",
	"Connecteur EDC","Bus Canopsis","AMQP","5672",

.. |run_manager| image:: ../_static/images/connectors/InterconnecionBDD.png
