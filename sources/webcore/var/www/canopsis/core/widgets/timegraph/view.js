define([
	'app/application',
	'components/widgets/timegraph/controller'
], function(Application, WidgetTimegraphController) {
	Application.WidgetTimegraphView = Application.WidgetView.extend({
		controller: WidgetTimegraphController
	});

	return Application.WidgetTimegraphView;
});