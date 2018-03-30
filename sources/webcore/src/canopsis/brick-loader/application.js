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

define([], function() {
    var get = Ember.get;

    if(window.appShouldNowBeLoaded !== true) {
        console.error('Application module is required too early, and it is probably leading to bad application behaviour and errors. Please do NOT require "app/application" in your modules.');
    }

    var Application = Ember.Application.create({
        LOG_ACTIVE_GENERATION: false,
        LOG_TRANSITIONS: false,
        LOG_TRANSITIONS_INTERNAL: false,
        LOG_VIEW_LOOKUPS: false,
        LOG_BINDINGS: true,
        rootElement: '#applicationdiv'
    });

    Application.deferReadiness();


    Ember.Object.reopen({
        toJson: function() {
            return JSON.parse(JSON.stringify(this));
        },
        json: function() {
            return JSON.parse(JSON.stringify(this));
        }.property()
    });

    var controllerDict = {
        init: function() {
            if(get(this, 'isGenerated')) {
                console.error('Ember is Instantiating a generated controller for "' + get(this, '_debugContainerKey') + '". This practice is not encouraged, as it might also be an underlying requireJS problem.', this);
            }
            this._super.apply(this, arguments);
        }
    };

    Ember.Controller.reopen(controllerDict);
    Ember.ArrayController.reopen(controllerDict);
    Ember.ObjectController.reopen(controllerDict);

    DS.ArrayTransform = DS.Transform.extend({
        deserialize: function(serialized) {
            if (Ember.typeOf(serialized) === 'array') {
                return serialized;
            }

            return [];
        },

        serialize: function(deserialized) {
            var type = Ember.typeOf(deserialized);

            if (type === 'array') {
                return deserialized;
            } else if (type === 'string') {
                return deserialized.split(',').map(function(item) {
                    return jQuery.trim(item);
                });
            }

            return [];
        }
    });

    DS.IntegerTransform = DS.Transform.extend({
        deserialize: function(serialized) {
            if (typeof serialized === "number") {
                return serialized;
            } else {
                // console.warn("deserialized value is not a number as it is supposed to be", arguments);
                return 0;
            }
        },

        serialize: function(deserialized) {
            return Ember.isEmpty(deserialized) ? null : Number(deserialized);
        }
    });

    DS.ObjectTransform = DS.Transform.extend({
        deserialize: function(serialized) {
            if (Ember.typeOf(serialized) === 'object') {
                return Ember.Object.create(serialized);
            }

            return Ember.Object.create({});
        },

        serialize: function(deserialized) {
            var type = Ember.typeOf(deserialized);

            if (type === 'object' || type === 'instance') {
                return Ember.Object.create(deserialized);
            } else {
                console.warn("bad format", type, deserialized);
            }

            return null;
        }
    });

    Application.Router.map(function() {
        this.resource('userview', { path: '/userview/:userview_id' });
    });

    loader.setApplication(Application);

    Ember.Application.initializer({
        name: 'AppLaunch',
        after: ['DataUtils', 'InflexionsRegistry'],
        initialize: function(container, application) {
            var dataUtils = container.lookupFactory('utility:data');
            var inflectionsManager = container.lookupFactory('registry:inflexions');

            dataUtils.setEmberApplicationSingleton(Application);
            inflectionsManager.loadInflections();
        }
    });

    Ember.Application.initializer({
        name:"RESTAdaptertransforms",
        after: "transforms",
        initialize: function(container, application) {
            void (container);
            application.register('transform:array', DS.ArrayTransform);
            application.register('transform:integer', DS.IntegerTransform);
            application.register('transform:object', DS.ObjectTransform);
        }
    });
    Application.advanceReadiness();
    window.$A = Application;

    //TODO create a more generic hook to provide a signal when app is ready, and put test handling in test brick
    //This is the test environment entry point.
    if(window.environment === 'test') {
        window.startCanopsisTests(Application);
    }

    if(Ember.TEMPLATES.application === undefined) {
        var tpl = Ember.Handlebars.compile('{{outlet "recordinfopopup"}}{{outlet "formwrapper"}}{{outlet}}');
        Ember.TEMPLATES['application'] = tpl;
    }

    return Application;
});
