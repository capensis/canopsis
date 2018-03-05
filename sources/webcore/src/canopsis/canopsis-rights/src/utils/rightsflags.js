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
    name: 'RightsflagsUtils',
    initialize: function(container, application) {
        var rightsflagsUtils = {
            canRead: function(checksum) {
                return (checksum >> 2) % 2 === 1;
            },
            canWrite: function(checksum) {
                return (checksum >> 1) % 2 === 1;
            },

            canCreate: function(checksum) {
                return (checksum >> 3) % 2 === 1;
            },
            canUpdate: function(checksum) {
                return this.canWrite(checksum);
            },
            canDelete: function(checksum) {
                return checksum % 2 === 1;
            }
        };

        application.register('utility:rightsflags', rightsflagsUtils);
    }
});
