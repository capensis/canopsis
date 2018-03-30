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
    name: 'component-serieitemeditor',
    initialize: function(container, application) {

        var schemasregistry = window.schemasRegistry;


        var get = Ember.get,
            set = Ember.set;


        var component = Ember.Component.extend({
            init: function() {
                this._super(arguments);

                set(this, "componentDataStore", DS.Store.create({
                    container: get(this, "container")
                }));

                var modelname = 'stylizedserie';
                var model = schemasregistry.getByName(modelname).EmberModel.proto();
                console.log('Fetch model:', modelname, model);

                var item = {};
                var me = this;

                console.group('Create virtual attributes for serieitem:');

                model.eachAttribute(function(name, attr) {
                    var contentKey = 'content.value.' + name;
                    var itemKey = 'item.' + name + '.value';

                    var val = get(me, contentKey);
                    var defaultVal = get(attr, 'options.default');
                    var value = val || defaultVal;

                    item[name] = Ember.Object.create({
                        value: val || defaultVal,
                        model: attr
                    });

                    me.addObserver(itemKey, function() {
                        var val = get(me, itemKey);
                        set(me, contentKey, val);
                    });

                    //ensure initilize content
                    if (value !== undefined) {
                        set(me, contentKey, value);
                    }

                    console.log(name, val, defaultVal, item[name]);
                });

                console.groupEnd();

                set(this, 'item', Ember.Object.create(item));
            }
        });

        application.register('component:component-serieitemeditor', component);
    }
});
