Translation
===========

This documentation is about how to translate UI into different languages with transifex system.

* first it is necessary to install transifex client : ``apt-get install transifex-client`` for debian users.

* Then you have to run ``tx init`` command.

.. role:: latex(code)
	Creating .tx folder...
	Transifex instance [https://www.transifex.com]:
	Creating skeleton...
	Creating config file...
	No authentication data found.
	No entry found for host https://www.transifex.com. Creating...
	Please enter your transifex username: username@mymail.org
	Password:
	Updating /home/utopman/.transifexrc file...
	Done.

* You now have your transifex configured update your  .po file with the new translation you wish update e.g. add

|	msgid "Test"
|	msgstr "Test"

at the end of your .po translation file.


it is time to push the lang file with ``tx push -s path_to_your_canopsis/sources/webcore/var/www/canopsis/resources/locales/lang-en.po``

Then the new translation should be available in transifex web UI for translation. You have to connect and translate new strings.

To complete the translation task, the new translated po file in each language and put them on the static server file where the build install will be able to download static po files on the following server ``http://repo.canopsis.org/locales/``.
