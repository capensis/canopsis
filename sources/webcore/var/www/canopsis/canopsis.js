require.config({
	baseUrl: '/static/',
	paths: {
		'app': 'canopsis/core',
		'schemas': 'canopsis/schemas',
		'etc': 'canopsis/etc',
		'lib': 'webcore-libs/dev',
		'text': 'webcore-libs/dev/text',
		'jquery': 'webcore-libs/dev/jquery-1.10.2',
		'plugins': 'webcore-libs/plugins/plugin',
		'consolejs': 'webcore-libs/console.js/console',
		'ember-cloaking': 'canopsis/core/lib/wrappers/ember-cloaking',
		'codemirror': 'webcore-libs/codemirror/codemirror',
		'colorpicker': 'webcore-libs/colorpicker/js/spectrum',
		'colorselector': 'webcore-libs/colorselector/js/bootstrap-colorselector',
		'seeds': 'webcore-libs/seeds',
		'jquery.encoding.digests.sha1': 'webcore-libs/jQuery.encoding.digests.sha1',
		'jquery.md5': 'webcore-libs/jquery.md5',
		'handlebars': 'webcore-libs/dev/handlebars-v1.3.0',
		'ember': 'canopsis/core/lib/wrappers/ember',
		'mmenu': 'canopsis/core/lib/wrappers/mmenu',
		'jsonselect': 'canopsis/core/lib/wrappers/jsonselect',
		'gridster': 'webcore-libs/dev/gridster/jquery.gridster',
		'timepicker': 'webcore-libs/dev/timepicker/bootstrap-datetimepicker.min',
		'moment': 'webcore-libs/dev/moment-with-langs.min',
		'ember-data': 'canopsis/core/lib/wrappers/ember-data',
		'ember-listview': 'webcore-libs/dev/ember-list-view',
		'ember-widgets': 'webcore-libs/ember-widgets/js/ember-widgets',
		'bootstrap': 'webcore-libs/bootstrap/current/js/bootstrap.min',
		'jqueryui': 'webcore-libs/dev/jquery-ui-1.10.3',
		'bootbox': 'webcore-libs/dev/bootbox.min',
		'jquerydatatables': 'webcore-libs/dev/jquery.dataTables',
		'bootstrapdatatables': 'webcore-libs/dev/dataTables.bootstrap',
		'colreorder': 'webcore-libs/dev/ColReorder',
		'adminLTE': 'canopsis/core/lib/wrappers/adminLTE',
		'utils': 'canopsis/core/lib/loaders/utils',
		'lodash': 'webcore-libs/dev/lodash.compat',
		'css3-mediaqueries': 'webcore-libs/min/css3-mediaqueries'
	},

	shim: {

		'jquerydatatables': {
			deps: ['jquery']
		},

		'icheck': {
			deps: ['jquery']
		},

		'bootstrapdatatables': {
			deps: ['jquery', 'jquerydatatables']
		},

		'jqueryui': {
			deps: ['jquery']
		},

		'consolejs': {
			deps: ["ember"]
		},

		'ember': {
			deps: ['jquery', 'handlebars']
		},

		'ember-cloaking': {
			deps: ['ember']
		},

		'ember-data': {
			deps: ['ember']
		},

		'ember-listview': {
			deps: ['ember']
		},

		'ember-widgets': {
			deps: ['ember', 'lodash', 'jqueryui', 'ember-listview']
		},

		'bootstrap': {
			deps: ['jquery']
		},

		'colorpicker': {
			deps: ['jquery']
		},

		'gridster': {
			deps: ['jquery']
		},

		'timepicker': {
			deps: ['jquery']
		}
	}
});

define(["canopsis/file_loader"], function () {
	require(['canopsis/main']);
});