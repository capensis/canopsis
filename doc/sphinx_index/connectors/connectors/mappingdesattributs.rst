Mapping des attributs
=====================

L'idée de cet onglet est de faire correspondre les attributs de l'application que vous souhaitez interconnecter avec les attributs Canopsis

Matrice des flux
----------------
.. csv-table::
   :header: "Attribut Canopsis","Explications","Obligatoire ?","Attribut applicatif correspondant","Exemple concrèt Fichier de log, requête SQL, Appel web service, Format de log, etc."
   :widths: 330, 20, 15, 15, 15

	"_id","Réservé",,,
	"event_id","Réservé",,,
	"connector : Connector type (gelf, nagios, snmp, ...)","Type de connecteur gelf, nagios, snmp, etc.","X",,
	"connector_name : Connector name (nagios1, nagios2 ...)","Nom du connecteur Valeur libre","X",,
	"event_type: Event type (check, log, trap, ...)","Type d'événement check, log, trap","X",,
	"source_type : Source type ('component' or 'resource')","Type de source 'component' ou 'resource'","X",,
	"component : Component name","Nom du composant","X",,
	"resource :  Ressource name","Nom de la ressource","X si source type = 'resource'",,
	"timestamp : UNIX seconds timestamp (UTC)","Timestamp au format UNIX Epoch",,,
	"state : State (0 (Ok), 1 (Warning), 2 (Critical), 3 (Unknown))","Etat
	0 -> OK
	1 -> WARNING
	2 -> CRITICAL
	3 -> UNKNOWN","X  si event_type = 'check'"
		 "state_type : State type (O (Soft), 1 (Hard))","Type d'état
	0 -> Soft
	1 -> Hard",,,
	"scheduled : (optional) True if this is a scheduled event","Réservé",,,
	"last_state_change : (reserved) Last timestamp after state change","Réservé",,,
	"previous_state : (reserved) Previous state (after change)","Réservé",,,
	"output : Event message","Message de l'événement",,,
	"long_output : Event long message","Message Long de l'événement",,,
	"tags : Event Tags (default: [])","Tags",,,
	"display_name : The name to display (customization purpose)","Nom sympathique",,,