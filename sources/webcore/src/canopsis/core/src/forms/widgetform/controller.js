/**
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'WidgetForm',
    after: ['FormFactory', 'WidgetsRegistry', 'FormsUtils', 'HashUtils'],
    initialize: function(container, application) {
        var FormFactory = container.lookupFactory('factory:form');
        var widgets = container.lookupFactory('registry:widgets');
        var formsUtils = container.lookupFactory('utility:forms');
        var hashUtils = container.lookupFactory('utility:hash');

        var get = Ember.get,
            set = Ember.set;


        var form = FormFactory('widgetform', {
            needs: ['application'],

            title: "Select a widget",

            widgets: widgets,

            availableWidgets: function() {
                console.log("availableWidgets");
                var widgets = [];

                for (var i = 0, l = widgets.all.length; i < l; i++) {
                    var currentWidget = widgets.all[i];

                    widgets.pushObject(currentWidget);
                }

                return widgets;
            }.property('widgets.all', "title"),

            actions: {
                submit: function(newWidgets) {
                    var newWidget = newWidgets[0];

                    console.log("onWidgetChooserSubmit", arguments);

                    console.group("attach new widget to widgetwrapper");

                    console.log("newWidget", newWidget);
                    console.log("widgetwrapper", this.newWidgetWrapper);

                    console.groupEnd();


                    this._super(this.newWidgetWrapper);
                },

                selectItem: function(widgetName) {
                    console.log('selectItem', arguments);

                    var containerwidget = get(this, 'formContext.containerwidget');
                    console.group('selectWidget', this, containerwidget, widgetName);

                    var store = get(this, 'formContext.containerwidget').store;
                    console.log('store to use', get(this, 'formContext.containerwidget').store);
                    var widgetId = hashUtils.generateId('widget_' + widgetName);

                    //FIXME this works when "xtype" is "widget"
                    var newWidget = store.createRecord(widgetName, {
                        'xtype': widgetName,
                        'listed_crecord_type': 'account',
                        'meta': {
                            'embeddedRecord': true,
                            'parentType': 'widgetwrapper'
                        },
                        'id': widgetId
                    });

                    this.newWidgetWrapper = store.push('widgetwrapper', {
                        'id': hashUtils.generateId('widgetwrapper'),
                        'xtype': 'widgetwrapper',
                        'title': 'wrapper',
                        'widget': newWidget,
                        'meta': {
                            'embeddedRecord': true,
                            'parentType': get(containerwidget, 'xtype'),
                            'parentId': get(containerwidget, 'id')
                        }
                    });

                    console.log('newWidgetWrapper', this.newWidgetWrapper);

                    console.log('newWidget', newWidget);
                    console.log('formwrapper', get(this, 'formwrapper'));

                    console.info('show embedded widget wizard');

                    formsUtils.showNew('modelform', newWidget, { formParent: this, title: "Add new " + widgetName });
                    console.groupEnd();
                }
            },

            partials: {
                buttons: []
            },

            parentContainerWidget: Ember.required(),
            parentUserview: Ember.required()
        });

        application.register('form:widgetform', form);
    }
});
