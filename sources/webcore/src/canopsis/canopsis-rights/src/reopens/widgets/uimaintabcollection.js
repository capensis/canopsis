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
    name: 'CanopsisRightsUimaintabcollectionWidgetReopen',
    after: ['RightsflagsUtils', 'UimaintabcollectionWidget', 'WidgetsRegistry'],
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        var widgetsRegistry = container.lookupFactory('registry:widgets');
        var UimaintabcollectionWidget = container.lookupFactory('widgetbase:uimaintabcollection');
        var rightsflagsUtils = container.lookupFactory('utility:rightsflags');

        UimaintabcollectionWidget.reopen({

            loggedaccountId: Ember.computed.alias('controllers.login.record._id'),
            loggedaccountRights: Ember.computed.alias('controllers.login.record.rights'),

            isViewDisplayable: function(viewId) {
                var user = get(this, 'loggedaccountId'),
                    rights = get(this, 'loggedaccountRights');

                if (user === 'root') {
                    return true;
                }

                return viewId && rightsflagsUtils.canRead(get(rights, viewId + '.checksum'));
            },

            userCanShowEditionMenu: function() {
                if(get(this, 'loggedaccountId') === "root") {
                    return true;
                }

                var rights = get(this, 'loggedaccountRights');

                if (get(rights, 'tabs_showeditionmenu.checksum')) {
                    return true;
                }

                return false;
            }.property(),

            userCanEditView: function() {
                if(get(this, 'loggedaccountId') === "root") {
                    return true;
                }

                var rights = get(this, 'loggedaccountRights'),
                    viewId = get(this, 'currentViewId');

                viewId = viewId.replace('.', '_');

                if (rightsflagsUtils.canWrite(get(rights, viewId + '.checksum'))) {
                    return true;
                }

                return false;
            }.property('currentViewId'),

            userCanCreateView: function() {
                if(get(this, 'loggedaccountId') === "root") {
                    return true;
                }

                var rights = get(this, 'loggedaccountRights');

                if (get(rights, 'userview_create.checksum')) {
                    return true;
                }

                return false;
            }.property()
        });
        
        application.register('widget:uimaintabcollection', UimaintabcollectionWidget);

    }
});
