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
    name: 'UserviewController',
    after: ['FormsUtils', 'HashUtils', 'InspectableItemMixin'],
    initialize: function(container, application) {

        var formUtils = container.lookupFactory('utility:forms');
        var hashUtils = container.lookupFactory('utility:hash');
        var schemasregistry = window.schemasRegistry;
        var InspectableItem = container.lookupFactory('mixin:inspectable-item');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;


        /**
         * @class UserviewController
         * @extends Ember.ObjectController
         * @constructor
         * @uses InspectableItemMixin
         * @uses Ember.Evented
         */
        var controller = Ember.ObjectController.extend(InspectableItem, Ember.Evented, {
            needs: ['application'],

            actions: {
                /**
                 * @event insertWidget
                 * @param {ContainerWidget} containerController
                 *
                 * Pop up a widget chooser form, and when the form sequence is over, inserts the widget into the container specified as first parameter
                 */
                insertWidget: function(containerController) {
                    console.log("insertWidget", containerController);
                    var widgetChooserForm = formUtils.showNew('widgetform', this);

                    var userviewController = this;

                    widgetChooserForm.submit.done(function(form, newWidgetWrappers) {
                        void (form);
                        var newWidgetWrapper = newWidgetWrappers[0];

                        console.log('onsubmit, adding widgetwrapper to containerwidget', newWidgetWrapper, containerController);
                        console.log('containerwidget items', get(containerController, 'model.items.content'));
                        //FIXME wrapper does not seems to have a widget
                        get(containerController, 'model.items.content').pushObject(newWidgetWrapper);

                        console.log("saving view");
                        get(userviewController, 'model').save();
                    });

                    widgetChooserForm.submit.fail(function() {
                        console.info('Widget insertion canceled');
                    });
                },

                /**
                 * @event duplicateWidgetAndContent
                 * @param {WidgetWrapper} widgetwrapperModel
                 * @param {ContainerWidget} containerModel
                 */
                duplicateWidgetAndContent: function(widgetwrapperModel, containerModel) {
                    copyWidget(this, widgetwrapperModel, containerModel);

                    this.get('model').save();

                },

                /**
                 * @event refreshView
                 */
                refreshView: function() {
                    console.log('refresh view');
                    this.trigger('refreshView');
                }
            }
        });


        /**
         * @method copyWidget
         * @param {UserviewController} self
         * @param {Widgetwrapper} widgetwrapperModel
         * @param {ContainerWidget} containerModel
         */
        copyWidget = function(self, widgetwrapperModel, containerModel) {
            var widgetwrapperJson = cleanRecord(widgetwrapperModel.toJSON());
            widgetwrapperJson.widget = null;
            var duplicatedWrapper = self.store.createRecord('widgetwrapper', widgetwrapperJson);

            var widgetJson = cleanRecord(widgetwrapperModel.get('widget').toJSON());
            var mixins = Ember.copy(get(widgetJson, 'mixins'));
            set(mixins, 'id', '')
            set(mixins, '_id', '');

            var duplicatedWidget = self.store.createRecord(widgetJson.xtype, widgetJson);

            if(!isNone(widgetwrapperModel.get('widget.items'))) {
                var items = widgetwrapperModel.get('widget.items.content');

                for (var i = 0; i < items.length; i++) {
                    var subWrapperModel = items[i];

                    copyWidget(self, subWrapperModel, duplicatedWidget);
                }
            }

            duplicatedWrapper.set('widget',  duplicatedWidget);
            containerModel.get('items.content').pushObject(duplicatedWrapper);
        }

        /**
         * @function cleanRecord
         * @param {Object} recordJSON
         * @return {Object} the cleaned record
         */
        function cleanRecord(recordJSON) {
            for (var key in recordJSON) {
                var item = recordJSON[key];
                //see if the key need to be cleaned
                if(key === 'id' || key === '_id' || key === 'widgetId' || key === 'preference_id') {
                    delete recordJSON[key];
                }

                //if this item is an object, then recurse into it
                //to remove empty arrays in it too
                if (typeof item === 'object') {
                    cleanRecord(item);
                }
            }

            if(recordJSON !== null && recordJSON !== undefined) {
                recordJSON['id'] = hashUtils.generateId(recordJSON.xtype || recordJSON.crecord_type || 'item');
                recordJSON['_id'] = recordJSON['id'];
            }

            return recordJSON;
        }

        application.register('controller:userview', controller);
    }
});
