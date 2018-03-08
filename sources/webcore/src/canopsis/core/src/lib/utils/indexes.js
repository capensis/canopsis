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
    name:"IndexesUtils",
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');
        var conf = window.canopsisConfiguration;

        var indexes = Utility.create({

            name: 'indexes',

            cache: {},
            load: function (collection) {
                $.ajax({
                    url: '/rest/indexes/' + collection,
                    success: function(data) {
                        indexes.cache[collection] = data.indexes;
                    },
                    async: false
                });
                return this.cache[collection];
            },
            parse: function (indexes) {

            },
            get: function(collection) {
                console.log(this);
                if (Ember.isNone(this.cache[collection])) {
                    this.load(collection);
                }
                return this.cache[collection];
            },
            getAsList: function (collection) {
                var indexes = this.get(collection);
                indexesList = [];

                for (var key in indexes) {
                    indexesList.push(indexes[key].key);
                }

                return indexesList;
            },
            getIndexWithField: function (collection, fields) {
                //Retruns all indexes from given collection that match one element in the fields array
                var indexes = this.getAsList(collection);
                var indexSelection = [];

                for (var x = 0, l = indexes.length; x < l; x++){
                    console.log('searching ' + indexes[x][0][0] + ' in' , fields);
                    if ($.inArray(indexes[x][0][0], fields) !== -1) {
                        indexSelection.push(indexes[x]);
                    }
                }

                return indexSelection;
            }
        });

        application.register('utility:indexes', indexes);
    }
});
