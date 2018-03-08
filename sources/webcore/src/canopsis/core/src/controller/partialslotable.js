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
    name: 'PartialslotAbleController',
    after: 'MixinsRegistry',
    initialize: function(container, application) {
        var mixinsRegistry = container.lookupFactory('registry:mixins');

        var get = Ember.get,
            set = Ember.set;

        /**
         * @function removeMixinsPartials
         * @private
         * @param widget
         * @param mixinName
         */
        function removeMixinsPartials(widget, mixinName) {
            console.log('removing mixin partials', arguments);

            var partials = get(widget, '_partials');

            var partialsToRemove = mixinsRegistry.getByName(mixinName.decamelize()).get('EmberClass.mixins')[0].properties.partials;

            for (var k in partialsToRemove) {
                if (partialsToRemove.hasOwnProperty(k)) {
                    for (var i = 0, l = partialsToRemove[k].length; i < l; i++) {
                        get(widget, '_partials')[k].removeObject(partialsToRemove[k][i]);
                    }
                }
            }

            set(widget, '_partials', partials);
        }

        //TODO put this in arrayutils
        function union_arrays (x, y) {
            if(x === undefined) {
                x = [];
            }
            if(y === undefined) {
                y = [];
            }

            var obj = {};
            for (var i = x.length-1; i >= 0; -- i)
                obj[x[i]] = x[i];
            for (var j = y.length-1; j >= 0; -- j)
                obj[y[j]] = y[j];
            var res = [];
            for (var k in obj) {
                if (obj.hasOwnProperty(k))  // <-- optional
                    res.push(obj[k]);
            }
            return res;
        }

        /**
         * @class PartialslotAbleController
         * @extends Ember.ObjectController
         * @constructor
         * @description Abstract class to manage partial slots and mixin with models, as used in Canopsis
         */
        var controller = Ember.ObjectController.extend({

            mergedProperties: ['partials'],

            /**
             * @property _partials
             * @type Object
             */
            _partials: {},

            /**
             * @property mixinOptions
             * @type Object
             */
            mixinOptions : {},


            /**
             * @method refreshPartialsList
             * @description recompute the "_partials" object. See "PartialslotAbleController#_partials" property"
             */
            refreshPartialsList: function() {
                console.log('refreshPartialsList', get(this, 'partials'));
                var partials = get(this, 'partials');
                set(this, '_partials', partials);
                var mixins = get(this, 'content.mixins');

                if(Ember.isArray(mixins)) {
                    for (var i = 0, l = mixins.length; i < l; i++) {
                        var currentMixin = mixins[i];

                        //DEPRECATE handle progressive deprecation of mixins as strings
                        if(typeof currentMixin === 'string') {
                            Ember.deprecate('Defining mixins as strings is deprecated. The new format is : \'{ name: "mixinName" }\'. This is required by the mixin options system.');
                        } else {
                            currentMixin = currentMixin.name.camelize();
                        }

                        partials = this.mergeMixinPartials(currentMixin, partials);
                    }
                }

                console.log('set partials for ', this, ' --> ', partials);
                set(this, '_partials', partials);
            },

            /**
             * @method mergeMixinPartials
             */
            mergeMixinPartials: function(Mixin, partials) {
                var me = this;
                var mixinName = Mixin.decamelize();
                var mixinRegistryEntry = mixinsRegistry.getByName(mixinName);

                console.log("mergeMixinPartials mixin:", Mixin);
                if(mixinRegistryEntry) {
                    var partialsToAdd = mixinRegistryEntry.EmberClass.mixins[0].properties.partials;

                    for (var k in partialsToAdd) {
                        if (partialsToAdd.hasOwnProperty(k)) {
                            var partialsArray = partialsToAdd[k];

                            var partialKey = '_partials.' + k;
                            set(this, partialKey, union_arrays(get(this, partialKey), partialsArray));
                        }
                    }
                    return partials;
                }
            }
        });
        application.register('controller:partialslot-able', controller);
    }
});
