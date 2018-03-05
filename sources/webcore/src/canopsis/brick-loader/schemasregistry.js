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

(function () {
    var AbstractClassRegistry = Ember.Object.extend({

        /**
         * The name of the registry
         * @property name
         * @type {string}
         */
        name: 'unnamed registry',

        all: [],
        byClass: {},

        /**
         * Aims to provide a way to inspect and display items
         *
         * @property tableColumns
         */
        tableColumns: [{title: 'name', name: 'name'}],

        /**
         * Appends the item into the "all" array, and into the corresponding class arrays in the "byClass" dict
         *
         * @method add
         * @param {object} item the item to add
         * @param {string} name the name of the item to add
         * @param {array} classes classes of the item
         */
        add: function(item, name, classes) {
            if(isNone(name)) {
                name = get(item, 'name');
            } else {
                set(item, 'name', name);
            }

            if(isNone(classes)) {
                classes = get(item, 'classes');
            } else {
                set(item, 'classes', classes);
            }

            console.log('registering item', get(item, 'name'), 'into registry', name, 'with classes', classes);
            this.all.pushObject(item);

            if(isArray(classes)) {
                for (var i = 0, l = classes.length; i < l; i++) {
                    if(isNone(this.byClass[classes[i]])) {
                        this.byClass[classes[i]] = Ember.A();
                    }

                    this.byClass[classes[i]].pushObject(item);
                }
            }
        },

        /**
         * Get an item by its name. Implemented because all must be migrated from an array to a dict
         *
         * @method getByName
         * @param {string} name the name of the item to get
         * @return {object} the object with the specified name
         */
        getByName: function(name) {
            for (var i = 0, l = this.all.length; i < l; i++) {
                if(get(this.all[i], 'name') === name) {
                    return this.all[i];
                }
            }
        },

        /**
         * Get a list of item that are registered in the specified class
         *
         * @method getByClassName
         * @param {string} name the name of the class
         * @return {array} the array of items that are defined with the specified class name
         */
        getByClassName: function(name) {
            return get(this.byClass, name);
        }
    });

    /**
     * Schemas Registry
     *
     * @class SchemasRegistry
     * @memberOf canopsis.frontend.core
     * @extends Abstractclassregistry
     * @static
     */
    window.schemasRegistry = AbstractClassRegistry.create({
        name: 'schemas',

        all: {},

        /**
         * Appends the item into the "all" array, and into the corresponding class arrays in the "byClass" dict
         */
        add: function(item, name) {
            this.all[name] = item;
        },

        /**
         * Get an item by its name. Implemented because all must be migrated from an array to a dict
         */
        getByName: function(name) {
            return this.all[name];
        },

        /**
         * Aims to provide a way to inspect and display items
         * Strictly typed object, at term, will not need this anymore
         */
        tableColumns: [{title: 'icon', name: 'icon'}, {title: 'name', name: 'name'}]
    });
})();


