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
    name:'ConsolemanagerMixin',
    after: ['MixinFactory', 'FormsUtils', 'DataUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var formUtils = container.lookupFactory('utility:forms');
        var dataUtils = container.lookupFactory('utility:data');

        var get = Ember.get,
            set = Ember.set,
            __ = Ember.String.loc;


        //TODO move this to development brick
        /**
         * Mixin allowing console and various js runtime settings
         *
         * @class ConsolemanagerMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('consolemanager', {
            init: function() {
                this.partials.statusbar.pushObject('consolemanagerstatusmenu' );
                this._super();
            },

            actions: {
                /**
                 * Shows a form to edit runtime settings
                 *
                 * @event showConsoleSettings
                 */
                showConsoleSettings: function(){
                    var jsruntimeconfigrecord = dataUtils.getStore().createRecord('jsruntimeconfiguration', {
                        id: 0,
                        selected_tags: window.console.tags._selectedTags,
                        colors: window.console.style._colors
                    });

                    var editForm = formUtils.showNew('modelform', jsruntimeconfigrecord, { title: __('Edit JS runtime configuration'), inspectedItemType: "jsruntimeconfiguration" });
                    console.log("editForm deferred", editForm.submit);
                    editForm.submit.done(function() {
                        console.log("jsruntimeconfigrecord saved", jsruntimeconfigrecord);
                    });
                    editForm.submit.always(function() {
                        console.log("jsruntimeconfigrecord always", jsruntimeconfigrecord);
                        window.console.tags._selectedTags = jsruntimeconfigrecord.get('selected_tags');
                        window.console.style._colors = jsruntimeconfigrecord.get('colors');
                        window.console.settings.save();
                        jsruntimeconfigrecord.unloadRecord();
                    });

                }
            },

            verbosity_mode: function() {
                return __("custom");
            }.property()
        });

        application.register('mixin:consolemanager', mixin);
    }
});
