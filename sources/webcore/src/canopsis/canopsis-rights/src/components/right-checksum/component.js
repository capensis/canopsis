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
    name:"component-right-checksum",
    after: 'RightsRegistry',
    initialize: function(container, application) {
        var rightsRegistry = container.lookupFactory('registry:rights');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component right-checksum
         * @description display buttons to edit rights checksums ("read/write", "chmod", etc...)
         */
        var component = Ember.Component.extend({
            /**
             * @property right
             * @description the right on which the checksum is edited
             * @type object
             */
            right: undefined,

            /**
             * @property checksum1flag
             * @description The first binary flag (at the right side) of the checksum
             * @type boolean
             */
            checksum1flag: undefined,

            /**
             * @property checksum2flag
             * @description The second binary flag of the checksum
             * @type boolean
             */
            checksum2flag: undefined,

            /**
             * @property checksum4flag
             * @description The third binary flag of the checksum
             * @type boolean
             */
            checksum4flag: undefined,

            /**
             * @property checksum8flag
             * @description The fourth binary flag of the checksum
             * @type boolean
             */
            checksum8flag: undefined,

            //TODO not used anymore? check and delete this property if possible
            /**
             * @property computedNumericChecksum
             */
            computedNumericChecksum: undefined,

            /**
             * @method init
             */
            init: function() {
                var right = get(this, 'right');

                if(isNone(get(right, 'data'))) {
                    set(right, 'data', Ember.Object.create());
                }

                if(!get(right, 'checksum')) {
                    set(right, 'checksum', 1);
                }

                var checksum = get(right, 'checksum');

                if(checksum >= 8) {
                    checksum -= 8;
                    set(this, 'checksum8flag', true);
                }

                if(checksum >= 4) {
                    checksum -= 4;
                    set(this, 'checksum4flag', true);
                }

                if(checksum >= 2) {
                    checksum -= 2;
                    set(this, 'checksum2flag', true);
                }

                if(checksum >= 1) {
                    checksum -= 1;
                    set(this, 'checksum1flag', true);
                }

                this._super();

                this.recomputeNumericChecksum();
            },

            /**
             * @property checksumType
             * @description computed property, dependant on "right.name". Retreives the checksum type, based on the right. Values can be either "RW", "CRUD", or by default they are considered as boolean (no flags used)
             * @type string
             */
            checksumType: function() {
                var value = get(this, 'right.name');
                var action = rightsRegistry.getByName(value);
                //FIXME don't use "_data"!
                if(action && action._data) {
                    return action._data.type;
                }
            }.property('right.name'),

            /**
             * @property checksumIsRW
             * @description computed property, dependant on "checksumType". True if the checksum type is "RW"
             * @type boolean
             */
            checksumIsRW: function() {
                return get(this, 'checksumType') === 'RW';
            }.property('checksumType'),

            /**
             * @property checksumIsCRUD
             * @description computed property, dependant on "checksumType". True if the checksum type is "CRUD"
             * @type boolean
             */
            checksumIsCRUD: function() {
                return get(this, 'checksumType') === 'CRUD';
            }.property('checksumType'),

            actions: {
                /**
                 * @method actions_toggleRightChecksum
                 * @description Action handling checksum edition
                 * @param flagNumber {integer} the flag to toggle
                 */
                toggleRightChecksum: function(flagNumber) {
                    var right = get(this, 'right');

                    console.info('toggleRightChecksum action', arguments);

                    var checksumFlagValue = get(this, 'checksum' + flagNumber + 'flag');

                    if(checksumFlagValue) {
                        set(this, 'checksum' + flagNumber + 'flag', false);
                    } else {
                        set(this, 'checksum' + flagNumber + 'flag', true);
                    }

                    var onChecksumChange = get(this, 'onChecksumChange');
                    var onChecksumChangeTarget = get(this, 'onChecksumChangeTarget');
                    if(onChecksumChange && onChecksumChangeTarget) {
                        onChecksumChangeTarget[onChecksumChange](right);
                    }
                }
            },

            /**
             * @property checksum8Class
             * @description Computed property, dependant on "checksum8flag". Css class for the fourth checksum
             */
            checksum8Class: function() {
                if(get(this, 'checksum8flag')) {
                    return 'btn btn-xs btn-success active';
                } else {
                    return 'btn btn-xs btn-danger';
                }
            }.property('checksum8flag'),

            /**
             * @property checksum4Class
             * @description Computed property, dependant on "checksum4flag". Css class for the third checksum
             */
            checksum4Class: function() {
                if(get(this, 'checksum4flag')) {
                    return 'btn btn-xs btn-success active';
                } else {
                    return 'btn btn-xs btn-danger';
                }
            }.property('checksum4flag'),

            /**
             * @property checksum2Class
             * @description Computed property, dependant on "checksum2flag". Css class for the second checksum
             */
            checksum2Class: function() {
                if(get(this, 'checksum2flag')) {
                    return 'btn btn-xs btn-success active';
                } else {
                    return 'btn btn-xs btn-danger';
                }
            }.property('checksum2flag'),

            /**
             * @property checksum1Class
             * @description Computed property, dependant on "checksum1flag". Css class for the first checksum
             */
            checksum1Class: function() {
                if(get(this, 'checksum1flag')) {
                    return 'btn btn-xs btn-success active';
                } else {
                    return 'btn btn-xs btn-danger';
                }
            }.property('checksum1flag'),

            /**
             * @method recomputeNumericChecksum
             * @description Observer, dependant on "checksum8flag", "checksum4flag", "checksum2flag", "checksum1flag". assign the computed checksum into the "right" object
             */
            recomputeNumericChecksum: function() {
                var checksum8flag = get(this, 'checksum8flag'),
                    checksum4flag = get(this, 'checksum4flag'),
                    checksum2flag = get(this, 'checksum2flag'),
                    checksum1flag = get(this, 'checksum1flag'),
                    numericChecksum = 0;

                if(checksum8flag) {
                    numericChecksum += 8;
                }

                if(checksum4flag) {
                    numericChecksum += 4;
                }

                if(checksum2flag) {
                    numericChecksum += 2;
                }

                if(checksum1flag) {
                    numericChecksum += 1;
                }

                set(this, 'computedNumericChecksum', numericChecksum);
                set(this, 'right.checksum', numericChecksum);
            }.observes('checksum8flag', 'checksum4flag', 'checksum2flag', 'checksum1flag')
        });
        application.register('component:component-right-checksum', component);
    }
});
