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
    name: 'component-actionbutton',
    after: 'ActionsUtils',
    initialize: function(container, application) {

        var actionsUtils = container.lookupFactory('utility:actions');

        /**
         * @component actionbutton
         * @description shows a button that triggers an action. This component requires a nested "yield" template.
         */
        var component = Ember.Component.extend({
            /**
             * @property action
             * @type string
             * @description the action name. It must be an action handled by the "action" utility
             */
            action: undefined,

            actions: {
                /**
                 * @method actions_doAction
                 * @argument action
                 * @argument {array} params
                 */
                doAction: function (actionName, param) {
                    actionsUtils.doAction(actionName, param);
                }
            }
        });

        application.register('component:component-actionbutton', component);
    }
});
