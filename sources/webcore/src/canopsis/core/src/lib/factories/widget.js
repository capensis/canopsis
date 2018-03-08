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
    name: 'WidgetFactory',
    after: ['WidgetsRegistry', 'WidgetController', 'NotificationUtils', 'SchemasLoader'],
    initialize: function(container, application) {
        var WidgetsRegistry = container.lookupFactory('registry:widgets'),
            schemasregistry = window.schemasRegistry,
            WidgetController = container.lookupFactory('controller:widget'),
            notificationUtils = container.lookupFactory('utility:notification'),

            get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * Widget factory. Creates a controller, stores it in Application
         * @param widgetName {string} the name of the new widget. lowercase
         * @param classdict {dict} the controller dict
         * @param options {dict} options :
         *            - subclass: to handle widget's controller inheritance: default is WidgetController
         *            - templateName: to use another template in the editor
         *
         * @author Gwenael Pluchon <info@gwenp.fr>
         */
        function Widget(widgetName, classdict, options) {

            var extendArguments = [];

            if (options === undefined) {
                options = {};
            }

            //This option allows to manually define inheritance for widgets
            if (options.subclass === undefined) {
                options.subclass = WidgetController;
            }

            //TODO check if this is still needed, as mixins are in configuration now
            if (options.mixins !== undefined) {
                for (var i = 0, l = options.mixins.length; i < l; i++) {
                    extendArguments.push(options.mixins[i]);
                }
            }

            extendArguments.push(classdict);

            var controllerClass = WidgetController.extend.apply(WidgetController, extendArguments);

            return controllerClass;
        }

        application.register('factory:widget', Widget);
    }
});
