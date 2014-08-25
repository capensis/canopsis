Run\_report
===========

With ``run_report`` command you can generate report (via webservice) to
Canopsis.

Authentication key
------------------

Firstly, you must add one user in ``CPS_report_admin`` group and
``CPS_view_admin`` in secondary group.

Secondly, get user's ``Authentication key`` in web UI
(``Build -> Accounts -> right-click -> Authentication key``)

How to find my view id ?
------------------------

You can use command:

\`run\_report -l
aa5b19862b6e19329999e89a21aa4aa864ab8339cc0b4eb57506bfa2

View ID Name
------------

view.\ *default*.dashboard Dashboard

view.root.1355336764158-7 test

view.root.1355391394043-3 VW\_Canopsis\`

Usage
-----

``usage: run_report [-h] [-o OUTPUT_FILE] [-f FILENAME] [-s SERVER] [--no-file]                   [-l] [-d]                   authkey [time_interval] [view_id]``

Examples
--------

Now you can use command:

run\_report **aa5b19862b6e19329999e89a21aa4aa864ab8339cc0b4eb57506bfa2**
**1d** **view.root.1355336764158-7**

-  view.root.1355336764158-7\_\_2012-12-12\_2012-12-13.pdf saved

-  **aa5b19862b6e19329999e89a21aa4aa864ab8339cc0b4eb57506bfa2** : Key
   Authentification of Report User

-  **1d** : Periodicity

-  **view.root.1355336764158-7** : View ID


