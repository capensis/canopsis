({
	baseUrl: 'static/',
	paths: {
		'app': 'canopsis/core',
		'lib': 'webcore-libs/dev',
		'text': 'webcore-libs/dev/text',

		'moment': 'webcore-libs/moment.min',
		'jquery': 'webcore-libs/dev/jquery-1.10.2',
		'handlebars': 'webcore-libs/dev/handlebars-1.0.0',
		'ember': 'webcore-libs/dev/ember',
		'ember-data': 'webcore-libs/dev/ember-data',
		'bootstrap': 'webcore-libs/bootstrap/3.0.3/js/bootstrap.min'
	},

	shim: {
		'ember': {
			deps: ['jquery', 'handlebars']
		},

		'ember-data': {
			deps: ['ember']
		},

		'bootstrap': {
			deps: ['jquery']
		}
	},

	name: 'canopsis/canopsis',
	out: 'canopsis/canopsis.min.js'
})