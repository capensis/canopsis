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
    name: 'MixineditdropdownView',
    initialize: function(container, application) {
        var schemasRegistry = window.schemasRegistry;

        var set = Ember.set,
            get = Ember.get,
            isNone = Ember.isNone;

        //TODO @gwen check if it's possible to remove this class

        /**
         * @class MixineditdropdownView
         * @extends Ember.View
         * @constructor
         */
        var view = Ember.View.extend({
            tagName: 'span',
            templateName: 'mixineditdropdown',

            hasEditableMixins: function () {
                return get(this, 'editableEnabledMixins.length') || get(this, 'wrapperMixins.length');
            }.property('editableEnabledMixins', 'wrapperMixins'),

            wrapperMixins: function () {
                var mixins = Ember.A();
                if (get(this, 'isGridLayout')) {
                    mixins.pushObject({'name': 'gridlayout'});
                }
                return mixins;
            }.property('isGridLayout'),

            editableEnabledMixins: function () {
                var mixins = get(this, 'mixins');
                var editableMixins = Ember.A();
                if(mixins) {
                    for (var i = 0; i < mixins.length; i++) {
                        if(schemasRegistry.getByName(mixins[i].name.camelize())) {
                            editableMixins.pushObject(mixins[i]);
                        }
                    }
                }

                return editableMixins;
            }.property('mixins')
        });

        application.register('view:mixineditdropdown', view);
    }
});
