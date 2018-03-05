/**
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 * This file is part of Canopsis.
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'component-statemapping',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component statemapping
         */
        var component = Ember.Component.extend({
            /**
             * @property placeholder
             * @default __('write template')
             */
            placeholder: __('write template'),

            /**
             * Load data into the component when data are received on init
             */

            didInsertElement: function () {
                var content = get(this, 'content');
                var initialized = get(this, 'initialized');

                if (isNone(initialized) && !isNone(content)) {
                    Ember.setProperties(this, {
                        info: content.info,
                        minor: content.minor,
                        major: content.major,
                        critical: content.critical,
                        initialized: true
                    });
                }
            }.observes('content'),

            /**
             * Generate the content value from component inputs
             */

            updateContent: function () {
                var content = {
                    info: get(this, 'info'),
                    minor: get(this, 'minor'),
                    major: get(this, 'major'),
                    critical: get(this, 'critical')
                };
                console.log('statemapping', content);
                set(this, 'content', content);
            }.observes('critical', 'major', 'minor', 'info')

        });

        application.register('component:component-statemapping', component);
    }
});
