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

define(['ember-lib', 'ember-data-lib'], function () {

    /*
     * Here is the canopsis UI main configuration file.
     * It is possible to add properies and values that are reachable
     * from the whole application through the namespace Canopsis.conf.PROPERTY
     */
    var canopsisConfiguration = {
        DEBUG: false,
        VERBOSE: 1,
        showPartialslots: false,
        DISPLAY_SCHEMA_MANAGER: true,
        REFRESH_ALL_WIDGETS: true,
        TRANSLATE: true,
        SHOW_TRANSLATIONS: false,
        TITLE: 'Canopsis Sakura',
        SHOWMODULES: false,
        getUserLanguage: function(){
            var language = 'en';

            $.ajax({
                url: '/account/me',
                success: function(data) {
                    if (data.success && data.data && data.data.length && data.data[0].ui_language) {
                        language = data.data[0].ui_language;
                        console.log('Lang initialization succeed, default language for application is set to ' + language.toUpperCase());
                    } else {
                        $.ajax({
                            url: '/rest/object/frontend/cservice.frontend',
                            success: function(data) {
                                if (data.success && data.data && data.data.length && data.data[0].ui_language) {
                                    language = data.data[0].ui_language;
                                } else {
                                    //FIXME Turning off warning for test environment. This removal is not supposed to be here.
                                    if(environment !== 'test') {
                                        console.warn('Lang data fetch failed, default language for application is set to EN', data);
                                    }
                                    language = 'en';
                                }

                            },
                            async: false
                        });
                    }
                },
                async: false
            }).fail(function () {
                console.error('Lang initialization failed, default language for application is set to EN');
                i18n.uploadDefinitions();
            });

            return language;
        },
        getEnabledModules: function (callback) {
            $.get('enabled-bricks.json', function (data) {
                if (data) {
                    if (typeof data === "string") {
                        data = JSON.parse(data);
                    }

                    callback(data);
                } else {
                    console.error('Could not load module information.');
                }

            });
        }
    };

    window.canopsisConfiguration = canopsisConfiguration;

    return canopsisConfiguration;
});
