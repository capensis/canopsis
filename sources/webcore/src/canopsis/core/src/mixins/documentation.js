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
    name: 'DocumentationMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        /**
         * Put a "view documentation button on the status bar"
         *
         * @class DocumentationMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('Documentation', {
            /**
             * @property showDocumentationButton
             * @type boolean
             * @description whether to show or not the documentation button. The property value must be assigned before object constructor call, or else it will be ignored
             */
            showDocumentationButton: true,

            init: function() {
                this.partials.statusbar.pushObject('documentation');

                this._super();
            }
        });

        application.register('mixin:documentation', mixin);
    }
});
