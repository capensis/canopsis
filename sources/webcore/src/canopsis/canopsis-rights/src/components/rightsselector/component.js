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
    name: 'component-rightsselector',
    after: ['RightsRegistry', 'component-dictclassifiedcrecordselector'],
    initialize: function(container, application) {
        var rightsRegistry = container.lookupFactory('registry:rights');
        var DictclassifiedcrecordselectorComponent = container.lookupFactory('component:component-dictclassifiedcrecordselector');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component rightselector
         * @description Allows to assign rights to a record, and to manage their checksums and options
         *
         * ![Component preview](../screenshots/component-rightselector.png)
         *
         * @augments dictclassifiedcrecordselector
         */
        var component = DictclassifiedcrecordselectorComponent.extend({
            /**
             * @property nameKey
             * @description the property of the right to use as the name key on the component
             */
            nameKey: '_id',
            /**
             * @property idKey
             * @description the property of the right to use as the id key on the component
             */
            idKey: '_id',

            actions: {
                /**
                 * @method actions-selectItem
                 * @description recomputes the generated rights dictionnary when an item is selected
                 */
                selectItem: function() {
                    this.recomputeValue();
                },

                /**
                 * @method actions-selectItem
                 * @description recomputes the generated rights dictionnary when an item is unselected
                 */
                unselectItem: function() {
                    this.recomputeValue();
                }
            },

            /**
             * @property classifiedItems
             * @description Computed property dependant on "items" and "items.@each". Compute a structure with classified item each time the 'items' property changed
             */
            classifiedItems : function(){
                var items = get(this, 'items');
                var valueKey = get(this, 'valueKey') || get(this, 'valueKeyDefault');
                var nameKey = get(this, 'nameKey') || get(this, 'nameKeyDefault');

                console.log("recompute classifiedItems", get(this, 'items'), valueKey);

                var res = Ember.Object.create({
                    all: Ember.A(),
                    byClass: {}
                });

                for (var i = 0, l = items.length; i < l; i++) {
                    var currentItem = items[i];

                    var objDict = {
                        name: currentItem.get(nameKey)
                    };

                    if(valueKey) {
                        console.log('add valueKey', currentItem.get(valueKey));
                        objDict.value = currentItem.get(valueKey);
                        console.log('objDict value', objDict);
                    }

                    this.serializeAdditionnalData(currentItem, objDict);

                    res.all.pushObject(Ember.Object.create(objDict));

                    possibleClassSplit = objDict.name.split("_");
                    if(possibleClassSplit.length > 1) {
                        var className = possibleClassSplit[0];

                        if(isNone(res.byClass[className])) {
                            res.byClass[className] = [];
                        }

                        res.byClass[className].pushObject(objDict);
                    }
                }

                console.log('recompute classifiedItems done', res);
                return res;
            }.property('items', 'items.@each'),

            /**
             * @method recomputeValue
             * @description Observer on "selectionUnprepared", "selectionUnprepared.@each", "selectionUnprepared.@each.checksum". Recomputes the generated rights dictionnary
             */
            recomputeValue: function(){
                console.group('recomputeValue', get(this, 'selectionUnprepared'));

                var selection = get(this, 'selectionUnprepared');

                var buffer = {};
                if(selection && selection.length) {
                    for (var i = 0, l = selection.length; i < l; i++) {
                        var currentItem = selection[i];
                        console.log('iteration', currentItem, get(currentItem, 'checksum'));
                        set(buffer, currentItem.name, {
                            checksum: get(currentItem, 'checksum') || 19
                        });
                    }
                }

                console.log('buffer', buffer);

                set(this, 'content', buffer);
                console.groupEnd();
            }.observes('selectionUnprepared', 'selectionUnprepared.@each', 'selectionUnprepared.@each.checksum'),

            /**
             * @method findItems
             * @description Retreives the rights list to display
             */
            findItems: function() {
                var me = this;

                var store = this.get('store_' + get(this, 'elementId'));

                var query = {
                    start: 0,
                    limit: 10000
                };

                query.filter = JSON.stringify({'crecord_type': this.get('crecordtype')});
                console.log('findItems', this.get('crecordtype'), query);

                store.findQuery('action', query).then(function(result) {
                    me.set('widgetDataMetas', result.meta);
                    var items = result.get('content');

                    me.set('items', items);

                    Ember.run.scheduleOnce('afterRender', {}, function() {
                        me.rerender();
                    });

                    me.extractItems(items);
                });
            },

            /**
             * @method deserializeAdditionnalData
             * @param item
             */
            deserializeAdditionnalData: function(item) {
                console.log('deserializeAdditionnalData', arguments);
                return {checksum: item.checksum};
            },

            /**
             * @method serializeAdditionnalData
             * @param item
             */
            serializeAdditionnalData: function(additionnalData) {
                console.log('serializeAdditionnalData', arguments);
            }
        });
        application.register('component:component-rightsselector', component);
    }
});
