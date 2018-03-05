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
    name: 'CanopsisRightsShowviewbuttonMixinReopen',
    after: ['ShowviewbuttonMixin', 'FormsUtils'],
    initialize: function(container, application) {
        var formsUtils = container.lookupFactory('utility:forms');
        var ShowviewbuttonMixin = container.lookupFactory('mixin:showviewbutton');

        var get = Ember.get,
            __ = Ember.String.loc;


        ShowviewbuttonMixin.reopen({
            init: function () {
                this.get('partials.itemactionbuttons').pushObject('actionbutton-viewrights');
                return this._super();
            },

            actions: {
                editUserviewRights: function(view) {
                    var viewName = get(view, 'crecord_name');
                    console.log('editUserviewRights view', view, viewName);
                    var formTitle = __('Edit rights for view : ') +  '"' + viewName + '"';

                    formsUtils.showNew('viewrightsform', view, { title: formTitle});
                }
            }
        });
    }
});
