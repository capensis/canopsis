Browser 2 canopsis
==================

Information
-----------

.. csv-table:: Connector description
   :header: "Property", "Value"
   :widths: 20, 30

	"Name", "Browser 2 canopsis"
	"Function", "Generate metrics from a website url for canopsis system"
	"Status", "BETA"
	"Github repository", "http://github.com/capensis"
	"Version", "0.1"

.. csv-table:: Metric list
   :header: "metric", "information"
   :widths: 20, 30

   "metric", "information"


.. csv-table:: Configuration variables
   :header: "variable name", "usage"
   :widths: 20, 30

   "variable", "information"



* Setup instructions

* Use instructions

	This connector produces canopsis events from standard input whitch have to be feed by the `Browser Time tool <https://github.com/tobli/browsertime/releases>`_ .
	The command line tool in crontab should looks like:

	``java -jar browsertime-0.6-full.jar http://example.org -f json``

* Changelog

	* 24/04/2014: documentation and connector coding started


