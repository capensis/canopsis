Module ctemplate
================

Handlebars templating in Python, using PyBars module.

Usage
-----

.. code-block:: python

	from ctemplate import CTemplate

	tmpl = CTemplate(u'''
	<h1>{{store}}</h1>

	<ul>
	{{#foreach items}}
		<li>{{name}} : {{author}}</li>
	{{/foreach}}
	</ul>
	''')

	rendered = tmpl({
		'store': 'Books',
		'items': [
			{'name': 'Lord of the Rings', 'author': 'J.R.R. Tolkien'},
			{'name': 'Game of Thrones', 'author': 'George R.R. Martin'},
			{'name': 'Foundation', 'author': 'Isaac Asimov'}
		]
	})

API
---

``ctemplate.CTemplate.__init__(source)``

 * ``source`` : template source as **unicode** string

``ctemplate.CTemplate.register_helper(name, handler)``

 * ``name`` : helper's name as **unicode** string
 * ``handler`` : helper's handler, must be a callable object, see *PyBars* documentation for more informations
