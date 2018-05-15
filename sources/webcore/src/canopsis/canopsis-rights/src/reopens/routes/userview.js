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
    name: 'CanopsisRightsUserviewRouteReopen',
    after: ['UserviewRoute', 'RightsflagsUtils'],
    initialize: function(container, application) {
        var UserviewRoute = container.lookupFactory('route:userview');
        var rightsflagsUtils = container.lookupFactory('utility:rightsflags');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @class UserviewRoute
         * @extends AuthenticatedRoute
         * @constructor
         * @description UserviewRoute reopen
         */
        UserviewRoute.reopen({
            /**
             * @method beforeModel
             * @param {Transition} transition
             * @return {Promise}
             *
             * Ensure the target view can be displayed.
             * Otherwise, put a "hasToBeRedirected" flag into the transition, in order to handle the redirection in the "afterModel" method.
             */
            beforeModel: function(transition) {
                var applicationController = this.controllerFor('application'),
                    loginController = this.controllerFor('login');

                var viewId = get(transition, 'params.userview.userview_id');
                viewId = viewId.replace('.', '_');

                var checksum = get(loginController, 'record.rights.' + viewId + '.checksum');
                var userId = get(loginController, 'record._id');

                if(!(rightsflagsUtils.canRead(checksum) || viewId === 'view_404' || viewId === 'view_401' || userId === 'root')) {
                    set(transition, 'hasToBeRedirected', true);
                }
                return this._super(transition);
            },

            /**
             * @method afterModel
             * @param {Userview} view The resolved model instance
             * @param {Transition} transition
             * @return {Promise}
             *
             * If a "hasToBeRedirected" flag is present into the transition, handle the redirection.
             */
            afterModel: function(view, transition) {
                var hasToBeRedirected = get(transition, 'hasToBeRedirected');

                if(hasToBeRedirected) {
                    this.transitionTo('/userview/view.404');
                }

                return this._super(view, transition);
            },

            actions: {
                /**
                 * @event toggleEditMode
                 * Handle rights management when toggling edit mode.
                 */
                toggleEditMode: function () {
                    var loginController = this.controllerFor('login');
                    var viewId = get(this, 'controller.model.id');
                    viewId = viewId.replace('.', '_');

                    var userId = get(loginController, 'record._id');

                    var checksum = get(loginController, 'record.rights.' + viewId + '.checksum');

                    if(rightsflagsUtils.canWrite(checksum) || userId === 'root') {
                        //call the regular "toggleEditMode" action
                        this._super();
                    }
                }
            }
        });
    }
});
