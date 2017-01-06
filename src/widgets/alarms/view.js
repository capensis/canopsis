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
    name: 'AlarmsView',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This mixin is the widget's view.
         *
         * @class viewMixin
         * @memberOf canopsis.frontend.brick-calendar
         */
        var view = Ember.Mixin.create({


			/**
			 * disable loading css when the controller get data
 			 * @method disableLoading
 			 */
			/*disableLoading: function () {
				if (this.$('overlay') !== undefined) {
					this.$('.overlay').remove();
					this.$('.fa-spin').remove();
				}
			},*/

            /**
             * Create the fullcalendar at the beginning and catch every changed view
             * @method didInsertElement
             */
            didInsertElement: function () {
				/*this._super.apply(this, arguments);
                var globalView = this;
                var controller = get(globalView, 'controller');
				this.$('.box').append('<div class="overlay"><i class="fa fa-refresh fa-spin"></i></div>');

				controller.on('disableLoading',this, this.disableLoading);*/
			},

            /**
             * Disable all triggered actions and destroy all view's objects
             * @method willDestroyElement
             */
            willDestroyElement: function() {
            }
        });
        application.register('view:alarms', view);
    }
});
