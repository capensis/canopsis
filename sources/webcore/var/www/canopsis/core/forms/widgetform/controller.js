/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
	'ember',
	'app/application',
	'app/lib/factories/form',
	'utils',
	'app/lib/loaders/schema-manager',
	'app/controller/journal'
], function(Ember, Application, FormFactory, utils) {

	FormFactory('widgetform', {
		needs: ['journal'],

		title: "Select a widget",

		availableWidgets: function() {
			console.log("availableWidgets");
			var widgets = [];

			for (var key in Canopsis.widgets.all) {
				var currentWidget = Canopsis.widgets.all[key];
				currentWidget.name = key;
				widgets.push(currentWidget);
			}

			return widgets;
		}.property('Canopsis.widgets.all', "title"),

		actions: {
			show: function() {
				var widgets = [];

				for (var key in Canopsis.widgets.all) {
					var currentWidget = Canopsis.widgets.all[key];
					currentWidget.set('name', key);
					widgets.push(currentWidget);
				}

				this.set('availableWidgets', widgets);
				this._super();
			},

			submit: function(newWidgets) {
				var newWidget = newWidgets[0];

				this.get('controllers.journal').send('publish', 'create', 'widget');

				console.log("onWidgetChooserSubmit", arguments);

				console.group("attach new widget to widgetwrapper");

				console.log("newWidget", newWidget);
				console.log("widgetwrapper", this.newWidgetWrapper);
				// Ember.set(this, 'newWidgetWrapper.widget', newWidget);

				console.groupEnd();


				this._super(this.newWidgetWrapper);
			},

			selectWidget: function(widget) {
				var containerwidget = this.get('formContext.containerwidget');
				console.group('selectWidget', this, containerwidget, widget.name);

				var widgetId = utils.hash.generateId('widget_' + widget.name);

				//FIXME this works when "xtype" is "widget"
				var newWidget = utils.data.getStore().createRecord(widget.name, {
					'xtype': widget.name,
					'listed_crecord_type': 'account',
					'meta': {
						'embeddedRecord': true,
						'parentType': 'widgetwrapper'
					},
					'id': widgetId
				});

				this.newWidgetWrapper = this.get('formContext.content.store').push('widgetwrapper', {
					'id': utils.hash.generateId('widgetwrapper'),
					'xtype': 'widgetwrapper',
					'title': 'wrapper',
					'widget': widgetId,
					'widgetType': widget.name,
					'meta': {
						'embeddedRecord': true,
						'parentType': containerwidget.get('xtype'),
						'parentId': containerwidget.get('widgetId')
					}
				});

				console.log('newWidgetWrapper', this.newWidgetWrapper);

				console.log('newWidget', newWidget);
				console.log('formwrapper', this.get('formwrapper'));

				console.info('show embedded widget wizard');

				utils.forms.show('modelform', newWidget, {formParent: this});
				console.groupEnd();
			}
		},

		partials: {
			buttons: []
		},

		parentContainerWidget: Ember.required(),
		parentUserview: Ember.required()
	});

	return Application.WidgetformController;
});
