Spécifications connecteur Canopsis
==================================
 

Ce document permet de spécifier techniquement les différentes
possibilités d'interconnexion entre une application et l'outil Canopsis.

Pour ce faire, il vous est demandé de saisir certaines informations de
la manière suivante :

 

- Onglet Description

- Onglet Interconnexion pour le type sélectionné

- Onglet Mapping

- Onglet Conclusion

.. csv-table::
   :header: "Application ciblée", "Version"
   :widths: 15, 15

	"", ""

.. csv-table::
   :header: "Type d'interconnexion  possible", "Flux" , " ", "Fichiers", "", "BDD", "","","","CLI"
   :widths: 15, 15, 15, 15, 15, 15, 15, 15, 15, 15

	"", "AMQP", "API","Trap SNMP", "Log", "Excel/CSV","Oracle","SQL Server", "MySql", "Send_Event"
	"Cochez pour selectionner", "", "", "",  "", "", "",  "", "", ""


Si le type d'interconnexion possible n'est pas listé, vous devez le
spécifier.

 Une fois ce type renseigné, RDV sur l'onglet correspondant au type
d'interconnexion pour saisir de plus amples informations
