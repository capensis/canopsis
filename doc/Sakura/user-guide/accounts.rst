Accounts
========

Accounts are used to manage authentication all over Canopsis, from the UI to the Web API.

Authentication
--------------

Users' identity can be verified by two means.

* Login/Password : Asked at the creation of the account, the purpose of this identication system is mainly to log in to the UI
* Auth Key : a string that is generated at the creation of an user's account. It is mostly used to query the Web api. (see web API)

Access control
--------------

Canopsis manage groups in an Unix similar way. Elements of canopsis (views, selectors, consolidations...) have several keys to handle rights :

[TODO]

Users may not have the rights to view, modify or delete an element. When they have a permission problem, they will be notified as they process the operation :

* In the UI, a popup will appear at the top left corner.
* Using the web API, the request will be a 403 "Not authorized" error.

.. NOTE ::
  [TODO screenshots]

(see groups)