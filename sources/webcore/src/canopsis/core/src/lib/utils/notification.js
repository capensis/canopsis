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
    name: 'NotificationUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        targetcontroller = {
            temp_buffer: [],
            createNotification: function(notificationType, notificationMessage) {
                this.temp_buffer.push({
                    notificationType: notificationType,
                    notificationMessage: notificationMessage
                });
            }
        };



        var notification = Utility.create({

            name: 'notification',
            /**
             * Initialize the notification controller
             * when the controller is not set up, it stores all the messages in a buffer stack.
             */
            setController: function(controller) {
            },

            //will be defined when notification controller is called.
            info: function (message) {
                //TODO doing it clean
            },
            warning: function (message) {
                //TODO doing it clean
            },
            error: function (message) {
                //TODO doing it clean
            },
            help: function () {
                console.log("usage is: utils.notification.notificate('info'|'warning'|'error', 'my message');");
            }
        });

        application.register('utility:notification', notification);
    }
});

