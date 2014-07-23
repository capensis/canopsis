define([
	'ember-data',
	'app/application'
], function(DS, Application) {
	Application.Timegraph = Application.Crecord.extend({
		name: DS.attr('string'),
		description: DS.attr('string'),
		timegraphAttribute: DS.attr('string')
	});

	return Application.Timegraph;
});