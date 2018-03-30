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
    name: 'component-elementidselectorwithoptions',
    after: ['SearchmethodsRegistry', 'HashUtils'],
    initialize: function(container, application) {
        var searchmethodsRegistry = container.lookupFactory('registry:searchmethods');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component elementidselectorwithoptions
         */
        var component = Ember.Component.extend({
            additionnalSelectionButtons: ['removeitembutton'],

            searchMethod: function () {
                var searchMethodKey = get(this, 'content.model.options.source.searchMethod');
                return searchmethodsRegistry.getByName(searchMethodKey);
            }.property('content.model.options.source.searchMethod'),

            helpModal: function () {
                return get(this, 'searchMethod.help');
            }.property('searchMethod'),

            /**
             * Columns available on the "avaiable data" table. This computed property gets the list of columns defined on the schema, and add a special column with the add button
             * @property availableDataColumns
             * @type array
             */
            availableDataColumns: function() {
                var columns = get(this, 'content.model.options.source.columns');
                var result = Ember.A();

                for (var i = 0; i < columns.length; i++) {
                    var currentColumn = columns[i];
                    var title = get(currentColumn, 'title'),
                        name = get(currentColumn, 'name');

                    result.pushObject({
                        title: __(title),
                        name: name
                    });
                }

                result.pushObject({
                    action: 'select',
                    actionAll: (get(this, 'multiselect') === true ? 'selectAll' : undefined),
                    title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
                    style: 'text-align: center;'
                });

                return result;
            }.property('content.model.options.source.columns'),

            init: function() {
                if (isNone(get(this, 'selectedData'))) {
                    set(this, 'selectedData', []);
                }

                if (isNone(get(this, 'searchPhrase'))) {
                    set(this, 'searchPhrase', null);
                }

                if (isNone(get(this, 'multiselect'))) {
                    set(this, 'multiselect', true);
                }

                set(this, 'listedAvailableItems', []);

                this._super.apply(this, arguments);

                var store = DS.Store.create({
                    container: this.get('container')
                });

                set(this, 'componentDataStore', store);

                var query = {filter: {_id: undefined}}, me = this;

                if (!isNone(get(this,'content'))) {
                    if (get(this, 'multiselect') === true) {
                        var content = get(this, 'content') || [];
                        query.filter._id = {'$in': content};
                    }
                    else {
                        query.filter._id = get(this, 'content');
                    }

                    var dataSchema = get(this, 'content.model.options.source.schema');

                    store.findQuery(dataSchema, query).then(function(result) {
                        var listedAvailableItems = get(result, 'content') || [];

                        console.log('Received data:', l, content, listedAvailableItems);

                        for(var i = 0, l = listedAvailableItems.length; i < l; i++) {
                            listedAvailableItems.pushObject(listedAvailableItems[i]);
                        }

                        set(me, 'listedAvailableItems', listedAvailableItems);
                    });
                }
            },

            //WARNING duplication
            //TODO put this on a dedicated util
            rendererFor: function(attribute) {
                var type = get(attribute, 'type');
                var role = get(attribute, 'options.role');
                if(get(attribute, 'model.options.role')) {
                    role = get(attribute, 'model.options.role');
                }
                var subRole = get(attribute, 'options.items.role');
                if(role === 'array' && !isNone(subRole)) {
                    role = subRole;
                }

                var rendererName;
                if (role) {
                    rendererName = 'renderer-' + role;
                } else {
                    rendererName = 'renderer-' + type;
                }

                if (Ember.TEMPLATES[rendererName] === undefined) {
                    rendererName = undefined;
                }

                return rendererName;
            },

            /**
             * @property selectedItemHeaderRendererType
             * The renderer type for the selected elements' header renderer
             */
            selectedItemHeaderRendererType: function () {
                return this.rendererFor(get(this, 'selectedItemHeaderRendererAttr'));
            }.property('selectedItemHeaderRendererAttr'),

            /**
             * @property selectedItemHeaderRendererAttr
             * The attribute for the selected elements' header renderer
             */
            selectedItemHeaderRendererAttr: function () {
                var attr = {
                    type: get(this, 'content.model.options.items.type'),
                    options: get(this, 'content.model.options.items')
                };

                return attr;
            }.property(),

            actions: {
                select: function(element) {
                    var selected = get(this, 'content.value') || [];

                    if (selected.indexOf(element) < 0) {
                        console.log('Select element:', element);


                        var polymorphicTypeKey = get(this, 'content.model.options.items.polymorphicTypeKey');

                        //TODO manage default values
                        var insertedElement = {};

                        var typeKeyName = get(this, 'content.model.options.items.typeKeyName');
                        var typeKeyValue;

                        //TODO stop using polymorphicTypeKey, use sourceMappingKeys instead
                        if(polymorphicTypeKey) {
                            typeKeyValue = get(element, polymorphicTypeKey);
                        } else {
                            var modelName = get(this, 'content.model.options.items.model');
                            if(modelName) {
                                typeKeyValue = get(this, 'content.model.options.items.model');
                            }
                        }

                        if(typeKeyValue && typeKeyName) {
                            insertedElement[typeKeyName] = typeKeyValue;
                        }

                        //Manage mapping : take properties from the available element and put them in the generated object
                        var sourceMappingKeys = get(this, 'content.model.options.sourceMappingKeys');

                        for (var i = 0; i < sourceMappingKeys.length; i++) {
                            var currentMappingOptions = sourceMappingKeys[i];
                            insertedElement[currentMappingOptions.objectKey] = get(element, currentMappingOptions.sourceItemKey);
                        }

                        if (get(this, 'multiselect') === true) {
                            selected.pushObject(insertedElement);
                        }
                        else {
                            selected = [insertedElement];
                        }
                    }

                    set(this, 'content.value', selected);
                },

                unselect: function(element) {
                    console.error('unselect', element);
                    var selected = get(this, 'content.value');

                    var idx = selected.indexOf(element);

                    if (idx >= 0) {
                        console.log('Unselect element:', element);
                        selected.removeAt(idx);
                    }

                    set(this, 'content.value', selected);
                },

                selectAll: function() {
                    if (get(this, 'multiselect') === true) {
                        var listedAvailableItems = get(this, 'listedAvailableItems');

                        if(listedAvailableItems.length) {
                            for (var i = 0, l = listedAvailableItems.length; i < l; i++) {
                                this.send('select', listedAvailableItems[i]);
                            }
                        }
                    }
                },

                unselectAll: function() {
                    set(this, 'selectedData', []);
                },

                search: function(search) {
                    if(search) {
                        var searchMethod = get(this, 'searchMethod');

                        if(searchMethod) {
                            var mfilter = searchMethod.buildFilter(search);
                            set(this, 'searchPhrase', mfilter);
                        }
                    }
                    else {
                        set(this, 'searchPhrase', null);
                    }
                }
            }
        });
        application.register('component:component-elementidselectorwithoptions', component);
    }
});
