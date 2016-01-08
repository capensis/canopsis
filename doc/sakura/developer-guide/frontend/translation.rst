.. _dev-frontend-translation:

Translation
===========

UI translation
--------------

This documentation is about how to translate UI into different languages with transifex system.

* first it is necessary to install transifex client (here is for Debian users :

.. code-block:: bash

        apt-get install transifex-client

* Then you have to run ``tx init`` command. This will lead to the following trace :

.. code-block:: bash

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

* You now have your transifex configured update your .po file that will be sent to the transifex server. Po file definitions can be extracted from database by issuing the following command once logged as canopsis user.

.. code-block:: bash

        i18n extract

This command will fetch new untranslated definitions from database to a ``.po`` file. This .po file has to be uploaded to Transifex in order to allow project translators to start translate files.

Issue the following commant to upload the ``todo.po`` file

.. code-block:: bash

        tx push -s path_to_your_canopsis/locale/todo.po

Then the new translation should be available in transifex web UI for translation. You have to connect and translate new strings.

To complete the translation task, once all fieds translated via Transifex, go download the ``.po`` file from language you translated it and put it in the following location (let assume you translate canopsis ui to french)

.. code-block:: bash

        path_to_your_canopsis/locale/fr/ui_lang.po

Then this file will be reachable from the webserver that will be able to parse the file enabling tranlations into the whole UI.

Documentation translation
-------------------------

In order to translate documentation into another language the shinx project must be setup as described in the Sphinx documentation `Sphinx internationalization <http://sphinx-doc.org/intl.html>`_.
What is left to you is to produce or update ``.po`` files in the language documentation will be translated. This can be done by issuing the following command :

.. code-block:: bash

        make gettext

Then update the locale dir :

.. code-block:: bash

        sphinx-intl update -p _build/locale -l fr

Then translation can be done by updating po files into the locale folder. Translation must be done as following. For each original paragraph, a **msgid** is a reference to the string to translate and you have to fill the **msgstr** accordingly depending on  the language you wish translate the documentation.

When the `msgstr` are translated, you have to generate the documentation in the language you chose by typing the following commands in the canopsis's documentation folder:

.. code-block:: bash

        sphinx-intl build
        make -e SPHINXOPTS="-D language='fr'" html

That's all.
