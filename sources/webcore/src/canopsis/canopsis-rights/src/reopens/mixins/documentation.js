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
    name: 'CanopsisRightsDocumentationMixinReopen',
    after: ['DocumentationMixin'],
    before: ['ApplicationController', 'LoginController'],
    initialize: function(container, application) {
        var DocumentationMixin = container.lookupFactory('mixin:documentation');

        var get = Ember.get,
            __ = Ember.String.loc;


        DocumentationMixin.reopen({
            needs: ['login'],
            loggedaccountId: Ember.computed.alias('controllers.login.record._id'),
            loggedaccountRights: Ember.computed.alias('controllers.login.record.rights'),

            showDocumentationButton: function() {
                var rights = get(this, 'loggedaccountRights');

                if (get(rights, 'menu_documentation.checksum') || get(this, 'loggedaccountId') === 'root') {
                    return true;
                }

                return false;
            }.property('loggedaccountId', 'loggedaccountRights')
        });
    }
});
