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
    name: 'component-classifiedcrecordselector',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set;

        /**
         * @component classifiedcrecordselector
         */
        var component = Ember.Component.extend({
            model: undefined,
            modelfilter: undefined,
            data: undefined,

            nameKeyDefault: 'crecord_name',
            nameKey: 'crecord_name',
            idKeyDefault: 'id',
            idKey: 'id',

            items: [],

            selectedValue: undefined,

            init: function() {
                this._super(arguments);

                set(this, 'selectionUnprepared', Ember.A());

                //FIXME is store destroyed?
                this.set('store_' + get(this, 'elementId'), DS.Store.create({
                    container: this.get('container')
                }));

                var initialContent = get(this, 'content');
                console.log('init', initialContent);

                this.setInitialContent(initialContent);

                this.refreshContent();
            },

            actions: {
                do: function(action, item) {
                    this.targetObject.send(action, item);
                }
            },

            setInitialContent: function(initialContent) {
                var valueKey = get(this, 'valueKey') || get(this, 'valueKeyDefault');

                console.log('setInitialContent', valueKey);
                if(initialContent) {
                    if(valueKey) {
                        set(this, 'loadingInitialContent', 'true');
                    }

                    //for the "if(valueKey)" case, check extractItems method
                    if(!valueKey) {
                        if( typeof initialContent === 'string') {
                            //we assume here that there is only one value
                            set(this, 'selectionUnprepared', [{ 'name': initialContent}]);
                        } else {
                            if( typeof initialContent === 'object' && initialContent !== null) {
                                var buffer = [];

                                for (var key in initialContent) {
                                    if (initialContent.hasOwnProperty(key)) {
                                        var additionnalData = this.deserializeAdditionnalData(get(initialContent, key));
                                        additionnalData.name = key;
                                        buffer.pushObject(additionnalData);
                                    }
                                }

                                set(this, 'selectionUnprepared', buffer);
                            } else {
                                set(this, 'selectionUnprepared', []);
                            }

                        }
                    }
                }
            },

            /*
             * Compute a structure with classified item each time the 'items' property changed
             */
            classifiedItems : function() {
                var items = get(this, 'items');
                var valueKey = get(this, 'valueKey') || get(this, 'valueKeyDefault');
                var nameKey = get(this, 'nameKey') || get(this, 'nameKeyDefault');

                console.log('recompute classifiedItems', get(this, 'items'), valueKey);

                var res = Ember.Object.create({
                    all: Ember.A()
                });

                for (var i = 0, l = items.length; i < l; i++) {
                    var currentItem = items[i];
                    var objDict = { name: currentItem.get(nameKey) };
                    if(valueKey) {
                        console.log('add valueKey', currentItem.get(valueKey));
                        objDict.value = currentItem.get(valueKey);
                        console.log('objDict value', currentItem, currentItem.get(valueKey));
                    }

                    this.serializeAdditionnalData(currentItem, objDict);

                    res.all.pushObject(Ember.Object.create(objDict));
                }

                return res;
            }.property('items', 'items.@each'),

            selectionChanged: function(){
                console.log('selectionChanged');
                var selectionUnprepared = get(this, 'selectionUnprepared');
                var res;

                //witch value is used as data reference in the selection.
                var valueKey = get(this, 'valueKey');

                if(get(this, 'multiselect')) {
                    res = Ember.A();
                    console.log('Push', selectionUnprepared[0]);

                    if(valueKey) {
                        for (var i = 0, li = selectionUnprepared.length; i < li; i++) {
                            res.pushObject(selectionUnprepared[i]);
                        }
                    } else {
                        for (var j = 0, lj = selectionUnprepared.length; j < lj; j++) {
                            res.pushObject(selectionUnprepared[j]);
                        }
                    }
                } else {
                    if(Ember.isArray(selectionUnprepared)) {
                        if(valueKey) {
                            res = selectionUnprepared[0];
                        } else {
                            res = selectionUnprepared[0];
                        }
                    }
                }

                set(this, 'selection', res);
            }.observes('selectionUnprepared', 'selectionUnprepared.@each'),

            onDataChange: function() {
                console.error('refreshContent');
                this.refreshContent();
            }.observes('data.@each'),

            onModelFilterChange: function() {
                this.set('currentPage', 1);
                this.refreshContent();
            }.observes('modelfilter'),

            refreshContent: function() {
                this._super(arguments);

                this.findItems();

                console.log(this.get('widgetDataMetas'));
            },

            /*
             * Fetch items as crecords, for performance reasons (userviews slowed down the component a lot because of embedded records for instance)
             */
            findItems: function() {
                var me = this;

                var crecordtype = get(this, 'crecordtype');

                var store = this.get('store_' + get(this, 'elementId'));

                var query = {
                    start: 0,
                    limit: 10000,
                    filter: get(this, 'modelfilter')
                };

                console.log('findItems', this.get('crecordtype'), query);

                if(crecordtype === 'view')
                    crecordtype = 'userview';

                store.findQuery(crecordtype, query).then(function(result) {
                    me.set('widgetDataMetas', result.meta);
                    var items = result.get('content');
                    me.set('items', items);

                    Ember.run.scheduleOnce('afterRender', {}, function() { me.rerender(); });
                    me.extractItems(items);
                });
            },

            extractItems: function(items) {
                var valueKey = get(this, 'valueKey');
                var idKey = get(this, 'idKey')  || get(this, 'idKeyDefault');
                var nameKey = get(this, 'nameKey') || get(this, 'nameKeyDefault');

                var initialContent = get(this, 'content');

                console.log('extractItems', items, initialContent);
                if(valueKey) {
                    //Fetch values with ajax request content
                    var correspondingExtractedItem;

                    if(typeof initialContent === 'string') {
                        console.log('extractItems with valueKey', arguments, Ember.inspect(initialContent));

                        correspondingExtractedItem = items.findBy(idKey, initialContent);

                        console.log('correspondingExtractedItem', correspondingExtractedItem);
                        if(correspondingExtractedItem !== undefined) {

                            var selectionUnprepared = [{ name: get(correspondingExtractedItem, nameKey), value: get(correspondingExtractedItem, idKey)}];
                            set(this, 'selectionUnprepared', selectionUnprepared);
                        }
                    } else if( typeof initialContent === 'object' && initialContent !== null) {
                        var buffer = [];
                        var keys = Ember.keys(initialContent);

                        for (var i = 0, l = keys.length ; i < l ; i++) {
                            var key = keys[i];

                            if (initialContent.hasOwnProperty(key)) {
                                var prop = get(initialContent, key);
                                console.log('findBy', idKey, key, prop);
                                correspondingExtractedItem = items.findBy(idKey, prop);

                                var selectionObject = {
                                    name: get(correspondingExtractedItem, nameKey),
                                    value: get(correspondingExtractedItem, idKey)
                                };

                                buffer.pushObject(selectionObject);
                            }
                        }

                        set(this, 'selectionUnprepared', buffer);
                    }

                    set(this, 'loadingInitialContent', false);
                }
            },

            deserializeAdditionnalData: function(additionnalData) {
                void(additionnalData);

                console.log('deserializeAdditionnalData', arguments);
            },

            serializeAdditionnalData: function(additionnalData) {
                void(additionnalData);

                console.log('serializeAdditionnalData', arguments);
            }
        });

        application.register('component:component-classifiedcrecordselector', component);
    }
});
