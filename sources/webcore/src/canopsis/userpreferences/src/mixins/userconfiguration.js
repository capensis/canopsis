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
    name:'UserconfigurationMixin',
    after: ['MixinFactory', 'HashUtils', 'DataUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');
        var hashUtils = container.lookupFactory('utility:hash');
        var dataUtils = container.lookupFactory('utility:data');

        var get = Ember.get,
            set = Ember.set;

        // /**
        //  * DS.Store hack to make transparent userpreferences persistence (when saving and retreiving models)
        //  */
        DS.Store.reopen({
            serialize: function(record, options) {

                var preferences = {};
                var loginController = dataUtils.getLoggedUserController();

                if(record.userPreferencesModel.attributes.list && record.userPreferencesModel.attributes.list.length > 0) {
                    console.group('userpreferences for widget', record.get('title'), record);
                    var userpreferenceAttributes = record.userPreferencesModel.attributes.list;
                    var preference_id = get(record, 'preference_id'),
                        user = get(loginController,'record._id');

                    if (preference_id === undefined) {
                        preference_id = hashUtils.generate_GUID();
                        set(record, 'preference_id', preference_id);
                    }

                    for (var i = 0, l = userpreferenceAttributes.length; i < l; i++) {
                        console.log('userpreferenceAttributes', userpreferenceAttributes[i]);
                        preferences[userpreferenceAttributes[i].name] = get(record, userpreferenceAttributes[i].name);
                        console.log('key', userpreferenceAttributes[i].name,'value', get(record, userpreferenceAttributes[i].name));
                    }

                    var userConfiguration = {
                        widget_preferences: preferences,
                        crecord_name: user,
                        widget_id: get(record, 'id'),
                        widgetXtype: get(record, 'xtype'),
                        title: get(record, 'title'),
                        viewId: get(record, 'viewId'),
                        id: get(record, 'id') + "_" + user,
                        _id: get(record, 'id') + "_" + user,
                        crecord_type: 'userpreferences'
                    };

                    console.log('push UP', userConfiguration);

                    $.ajax({
                        url: '/rest/userpreferences/userpreferences',
                        type: 'POST',
                        data: JSON.stringify(userConfiguration)
                    });

                    console.groupEnd('userpreferences for widget', record.get('title'));
                } else {
                    console.log('no userpreferences to save for widget', record.get('title'));
                }

                return this._super.apply(this, arguments);
            },

            push: function(type, data) {
                var record = this._super.apply(this, arguments);
                var loginController = dataUtils.getLoggedUserController();

                var userpreferenceAttributes = record.userPreferencesModel.attributes.list;

                if(userpreferenceAttributes.length > 0) {
                    var user = get(loginController, 'record._id');

                    $.ajax({
                        url: '/rest/userpreferences/userpreferences',
                        async: false,
                        data: {
                            limit: 1,
                            filter: JSON.stringify({
                                crecord_name: user,
                                widget_id: get(record, 'id'),
                                _id: get(record, 'id') + '_' + user
                            })
                        },
                        success: function(data) {
                            if (data.success && data.data.length && data.data[0].widget_preferences !== undefined) {
                                console.log('User configuration load for widget complete', JSON.stringify(data));
                                var preferences = data.data[0].widget_preferences;

                                set(record, get(record, 'id') + "_" + user);
                                set(record, 'userPreferences', preferences);

                                for (var key in preferences) {
                                    console.log('User preferences: will set key', key, 'in widget', get(record, 'title'), preferences[key]);
                                    record.set(key, preferences[key]);
                                }

                            } else {
                                console.log('No user preference exists for widget' + get(record, 'title'));
                            }
                        }
                    }).fail(
                        function (error) {
                            void (error);
                            console.log('No user s preference found for this widget');
                        }
                    );
                }

                return record;
            }
        });

        /**
         * Provides userconfiguration handling for any kind of objects. It is usually applied on Canopsis widgets
         *
         * @mixin Userconfiguration
         */
        var mixin = Mixin('userconfiguration', {

            needs: ['login'],

            /**
             * Persists the user preferences into the backend
             * @method saveUserConfiguration
             */
             saveUserConfiguration: function () {
                var record = get(this, 'model');

                console.log('saveUserConfiguration', record);

                var preferences = {};

                if(record.userPreferencesModel.attributes.list && record.userPreferencesModel.attributes.list.length > 0) {
                    console.group('userpreferences for widget', record.get('title'), record);
                    var userpreferenceAttributes = record.userPreferencesModel.attributes.list;
                    var loginController = dataUtils.getLoggedUserController();
                    var preference_id = get(record, 'preference_id'),
                        user = get(loginController, 'record._id');

                    if (preference_id === undefined) {
                        record.save();
                    } else {
                        for (var i = 0, l = userpreferenceAttributes.length; i < l; i++) {
                            console.log('userpreferenceAttributes', userpreferenceAttributes[i]);
                            preferences[userpreferenceAttributes[i].name] = get(record, userpreferenceAttributes[i].name);
                            console.log('key', userpreferenceAttributes[i].name,'value', get(record, userpreferenceAttributes[i].name));
                        }

                        var userConfiguration = {
                            widget_preferences: preferences,
                            crecord_name: user,
                            widget_id: get(record, 'id'),
                            widgetXtype: get(record, 'xtype'),
                            title: get(record, 'title'),
                            viewId: get(record, 'viewId'),
                            id: preference_id + '_' + user,
                            _id: get(record, 'id') + "_" + user,
                            crecord_type: 'userpreferences'
                        };

                        $.ajax({
                            url: '/rest/userpreferences/userpreferences',
                            type: 'POST',
                            data: JSON.stringify(userConfiguration)
                        });
                    }
                    console.groupEnd('userpreferences for widget', record.get('title'));
                } else {
                    console.log('no userpreferences to save for widget', record.get('title'));
                }
            }
        });

        application.register('mixin:userconfiguration', mixin);
    }
});
