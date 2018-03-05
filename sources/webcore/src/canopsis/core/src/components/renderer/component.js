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
    name: 'component-renderer',
    after: 'DebugUtils',
    initialize: function(container, application) {
        var debugUtils = container.lookupFactory('utility:debug');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        /**
         * @component renderer
         * @description Component displaying the correct renderer for an attribute.
         * It is possible to specify the renderer type to use. If not specified, it will try to get the correct type on its own.
         *
         * @example {{component-renderer attr=attr record=controller.formContext value=attr.value}}
         */
        var component = Ember.Component.extend({
            /**
             * @method init
             */
            init: function() {
                var record = get(this, 'record'),
                    attrName = get(this, 'attrName');

                if(!isNone(attrName)) {
                    console.group('Fetch attribute from record');

                    console.log('record:', record);
                    console.log('attrName:', attrName);

                    if(!isNone(record)) {
                        var attr = get(record, 'constructor.attributes.' + attrName),
                            value = get(record, attrName);

                        console.log('attr:', attr);
                        console.log('value:', value);

                        var role;
                        if (!isNone(attr)) {
                            role = get(attr, 'options.role');
                        }

                        if(!isNone(role)) {
                            var renderer = 'renderer-' + role;

                            if(!isNone(Ember.TEMPLATES[renderer])) {
                                console.log('rendererType:', renderer);
                                set(this, 'rendererType', renderer);
                            }
                        }

                        set(this, 'attr', attr);
                        set(this, 'value', value);
                    }

                    console.groupEnd();
                }

                this._super.apply(this, arguments);
            },

            /**
             * @property canopsisConfiguration
             * @type object
             * @description the canopsis frontend configuration object
             */
            canopsisConfiguration: canopsisConfiguration,

            /**
             * @property debug
             * @description whether the UI is in debug mode or not
             * @type boolean
             * @default Ember.computed.alias('canopsisConfiguration.DEBUG')
             */
            debug: Ember.computed.alias('canopsisConfiguration.DEBUG'),

            actions: {
                /**
                 * @method actions_inspect
                 * @description inspects the object in the console (see debugUtils for more info)
                 */
                inspect: function() {
                    debugUtils.inspectObject(this);
                },

                /**
                 * @method actions_do
                 * @param {string} action
                 * @description sends an action to the parent controller. Every parameter after the first one is bubbled to the parent controller action
                 */
                do: function(action) {
                    var params = [];
                    for (var i = 1, l = arguments.length; i < l; i++) {
                        params.push(arguments[i]);
                    }

                    get(this, 'parentView.controller').send(action, params);
                }
            },

            /**
             * @property tagName
             * @type string
             * @default
             */
            tagName: 'span',

            //TODO check why there is a property dependant on "shown_columns" in here. As it is a List Widget property, it does not seems relevant at all.
            /**
             * @property attr
             * @description the rendered attribute
             */
            attr: function() {
                var shown_columns = get(this, 'shown_columns');
                for (var i = 0, l = shown_columns.length; i < l; i++) {
                    if(shown_columns[i].field === get(this, 'field')) {
                        return shown_columns[i];
                    }
                }
            }.property('shown_columns')
        });

        application.register('component:component-renderer', component);
    }
});
