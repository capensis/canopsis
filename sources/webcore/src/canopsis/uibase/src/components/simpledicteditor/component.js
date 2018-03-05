/*
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
    name: 'component-simpledicteditor',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component simpledicteditor
         * @description A simple key-value editor. Keys and values can only be edited as strings
         */
        var component = Ember.Component.extend({
            /**
             * @property content
             * @type object
             * @description the edited dictionnary
             */
            content: undefined,

            actions: {
                /**
                 * @method actions_addKey
                 * @description Adds the key to the edited dictionnary
                 */
                addKey: function() {
                    var newValue = get(this, 'newValue');
                    var newKey = get(this, 'newKey');

                    if(isNone(get(this, 'content'))) {
                        set(this, 'content', {});
                    }

                    if(!isNone(newValue) && !isNone(newKey)) {
                        set(this, 'content.' + newKey, newValue);
                        this.contentChanged();

                        set(this, 'newKey');
                        set(this, 'newValue');
                    }
                },

                /**
                 * @method actions_select
                 */
                select: function() {
                    console.warn('select', arguments);
                },
                /**
                 * @method actions_removeKey
                 * @param {string} keyToDelete
                 * @description Removes the key to the edited dictionnary
                 */
                removeKey: function(keyToDelete) {
                    delete this.content[keyToDelete];

                    this.contentChanged();
                }
            },

            /**
             * @method didInsertElement
             */
            didInsertElement: function() {
                this._super.apply(this);
                this.contentChanged();
            },

            /**
             * @method contentChanged
             */
            contentChanged: function(){
                var buffer = Ember.A();
                var content = get(this, 'content');

                if(content) {
                    var keys = Ember.keys(content);

                    console.warn(content, keys);

                    for (var i = 0, l = keys.length; i < l; i++) {
                        var currentKey = keys[i];
                        buffer.pushObject(Ember.Object.create({
                            key: currentKey,
                            value: content[currentKey]
                        }));
                    }
                    console.warn('dictAsArray', buffer);

                    set(this, 'dictAsArray', buffer);
                }
            }
        });

        application.register('component:component-simpledicteditor', component);
    }
});
