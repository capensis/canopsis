Report header Logo
==================

Create your custom files
------------------------

In order to customize canopsis you must create your own files in
/opt/canopsis/var/wkhtmltopdf/. The easiest way is to make a copy of
header.html and footer.html then rename them, for exemple
myCompanyHeader.html , myCompanyFooter.html .

::

    cp /opt/canopsis/var/wkhtmltopdf/header.html /opt/canopsis/var/wkhtmltopdf/myCompanyHeader.html
    cp /opt/canopsis/var/wkhtmltopdf/footer.html /opt/canopsis/var/wkhtmltopdf/myCompanyFooter.html

Now put your logo (if you want one) in the same folder, your logo must
be around 50px tall to fit in the header.

::

    cp my_company_log.png /opt/canopsis/var/wkhtmltopdf/my_company_log.png

Edit header
------------

In header and footer page you can put all the html you want. It's
rendered by wkhtmltopdf as a part of the webpage. In the marker change
the src name for the one of your file, and delete the 'style' attribut :

Before:

::

    <img src="logo_canopsis.png" style="height:58;width:150"/><br/>

After:

::

    <img src="my_company_logo.png"/><br/>

Edit wkhtmltopdf\_wrapper.json
------------------------------

In this file you need to change to line, first the header line, just
replace the file name.

::

    vim /opt/canopsis/etc/wkhtmltopdf_wrapper.json

Before:

::

    "header" : "--header-html /opt/canopsis/var/wkhtmltopdf/header.html",

After:

::

    "header": "--header-html /opt/canopsis/var/wkhtmltopdf/myCompanyHeader.html",

Same process for the footer file

Before :

::

    "footer": "--footer-html /opt/canopsis/var/wkhtmltopdf/footer.html",

After :

::

    "footer": "--footer-html /opt/canopsis/var/wkhtmltopdf/myCompanyFooter.html"