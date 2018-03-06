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
    name: 'CanopsisRightsCrudMixinReopen',
    after: ['CrudMixin', 'RightsflagsUtils'],
    initialize: function(container, application) {

        var CrudMixin = container.lookupFactory('mixin:crud');
        var rightsflagsUtils = container.lookupFactory('utility:rightsflags');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;


        CrudMixin.reopen({
            userCanReadRecord: function() {
                if(get(this, 'user') === "root") {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');
                var checksum = get(this, 'rights.models_' + crecord_type + '.checksum');

                return rightsflagsUtils.canRead(checksum);
            }.property('config.listed_crecord_type'),

            userCanCreateRecord: function() {
                if(get(this, 'user') === "root") {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');
                var checksum = get(this, 'rights.models_' + crecord_type + '.checksum');

                return rightsflagsUtils.canCreate(checksum);
            }.property('config.listed_crecord_type'),

            userCanUpdateRecord: function() {
                if(get(this, 'user') === "root") {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');
                var checksum = get(this, 'rights.models_' + crecord_type + '.checksum');

                return rightsflagsUtils.canUpdate(checksum);
            }.property('config.listed_crecord_type'),

            userCanDeleteRecord: function() {
                if(get(this, 'user') === "root") {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');
                var checksum = get(this, 'rights.models_' + crecord_type + '.checksum');

                return rightsflagsUtils.canDelete(checksum);
            }.property('config.listed_crecord_type')
        });
    }
});
