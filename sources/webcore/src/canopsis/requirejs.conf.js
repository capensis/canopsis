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

window.schemasToLoad = [];

require.config({
    waitSeconds: 30,
    baseUrl: '/static/',
    paths: {
        'text': 'canopsis/brick-loader/externals/requirejs-text/text',
        'link': 'canopsis/brick-loader/externals/requirejs-link/link',

        'jquery': 'canopsis/brick-loader/externals/jquery/dist/jquery',
        'handlebars': 'canopsis/brick-loader/externals/handlebars/handlebars',
        "ehbs" : 'canopsis/brick-loader/externals/requirejs-ember-handlebars/ehbs',
        'ember-template-compiler': 'canopsis/brick-loader/externals/ember-template-compiler',
        'ember-lib': 'canopsis/brick-loader/externals/ember.debug',
        'ember-data-lib': 'canopsis/brick-loader/externals/ember-data',
        'schemasregistry': 'canopsis/brick-loader/schemasregistry'
    },
    shim: {
        'ember-lib': {
            deps: ['jquery', 'ember-template-compiler', 'handlebars'],
            exports: 'Ember'
        },
        'ember-data-lib': {
            deps: ['ember-lib']
        },
        'schemasregistry': {
            deps: ['ember-lib', 'ember-data-lib']
        }
    },
    ehbs : {
        ember : 'ember-lib'
    }
});
require(['text!canopsis/brick-loader/bower.json'], function(loaderManifest) {
    loaderManifest = JSON.parse(loaderManifest);

    if(loaderManifest.envMode === 'production' && window.environment !== 'debug' && window.environment !== 'test') {
        require.config({
            waitSeconds: 30,
            baseUrl: '/static/',
            paths: {
                'jquery': 'canopsis/brick-loader/externals/jquery/dist/jquery.min',
                'handlebars': 'canopsis/brick-loader/externals/handlebars/handlebars.min',
                'ember-template-compiler': 'canopsis/brick-loader/externals/ember-template-compiler',
                'ember-lib': 'canopsis/brick-loader/externals/ember.min',
                'ember-data-lib': 'canopsis/brick-loader/externals/ember-data.min',
                'schemasregistry': 'canopsis/brick-loader/schemasregistry'
            },
            shim: {
                'ember-lib': {
                    deps: ['jquery', 'ember-template-compiler', 'handlebars'],
                    exports: 'Ember'
                },
                'ember-data-lib': {
                    deps: ['ember-lib']
                },
                'schemasregistry': {
                    deps: ['ember-lib', 'ember-data-lib']
                }
            },
            ehbs : {
                ember : 'ember-lib'
            }
        });
    }

    window.isIE = navigator.appName.indexOf('Internet Explorer') !== -1;

    if (isIE) {
        //this force console to use log method for early loaded
        //modules that could use other console methods.
        console.group = function () {};
        console.groupEnd = function() {};
        console.debug = console.log;
        console.warning = console.log;
        console.error = console.log;
        console.tags = {
            add: function() {},
            remove: function () {}
        };

        console.settings = {
            save: function() {}
        };
    }

    var setLoadingInfo = function(text, icon) {
        if(window.__) {
            text = window.__(text, true);
        }

        $('#loadingInfo').html(text);
        
        //Hack for let the user open several modals
        $(document).ready(function () {

            $('#openBtn').click(function () {
                $('#myModal').modal({
                    show: true
                })
            });

            $(document).on('show.bs.modal', '.modal', function (event) {
                var zIndex = 1040 + (10 * $('.modal:visible').length);
                $(this).css('z-index', zIndex);
                $(this).css('overflow', 'scroll');
                setTimeout(function () {
                    $('.modal-backdrop').not('.modal-stack').css('z-index', zIndex - 1).addClass('modal-stack');
                }, 0);
            });


        });
        
        if(icon) {
            $('#loading').append('<i class="fa '+ icon +'"></i>');
        }
    };

    setModuleInfo = function (modules, showmodules) {
        if (showmodules) {
            var title = '<h5>Enabled modules :</h5>';
            $('#moduleList').append(title + modules.join('<br />'));
        }
    };
    require(['ember-lib'], function() {

        if(window.environment === 'test') {
            Ember.deprecate = function() {};
            Ember.warn = function() {};
        }

        require(['canopsis/canopsisConfiguration',
                'canopsis/brick-loader/i18n',
                'canopsis/brick-loader/loader',
                'ember-data-lib',
                'schemasregistry'], function(canopsisConfiguration, i18n) {


            require([
                'text!canopsis/brick-loader/i18n/' + i18n.lang + '.json'
            ], function (langFile) {
                var langFile = JSON.parse(langFile);
                var langKeys = Em.keys(langFile);

                i18n.translations[i18n.lang] = {};
                for (var i = 0; i < langKeys.length; i++) {
                    i18n.translations[i18n.lang][langKeys[i]] = langFile[langKeys[i]];
                }

                window.canopsisConfiguration = canopsisConfiguration;

                var get = Ember.get;

                canopsisConfiguration.EmberIsLoaded = true;

                canopsisConfiguration.getEnabledModules(function (enabledPlugins) {

                    if (enabledPlugins.length === 0) {
                        alert('No module loaded in Canopsis UI. Cannot go beyond');
                    }

                    setLoadingInfo('Fetching frontend bricks', 'fa-cubes');
                    setModuleInfo(enabledPlugins, canopsisConfiguration.SHOWMODULES);
                    var language = i18n.lang;

                    if(!language) {
                        language = 'en';
                    }

                    var loc = Ember.String.loc;
                    Ember.String.loc = function (fieldToTranslate) {
                        i18n._(fieldToTranslate, true);
                        return loc(fieldToTranslate);
                    };

                    Ember.STRINGS = i18n.translations[language] || {};

                    var deps = [];

                    for (var i = 0; i < enabledPlugins.length; i++) {
                        var currentPlugin = enabledPlugins[i];

                        if(currentPlugin !== 'core') {
                            deps.push('text!canopsis/'+ currentPlugin +'/bower.json');
                        }
                    }
                    deps.push('text!canopsis/core/bower.json');

                    if(window.environment) {
                        deps.push('canopsis/environment.' + window.environment);
                    } else {
                        deps.push('canopsis/environment.production');
                    }

                    deps.push('canopsis/brick-loader/extend');
                    deps.push('link');

                    require(deps, function() {
                        var initFiles = [];
                        var schemasInitFiles = [];
                        window.bricks = {};

                        for (var i = 0, l = enabledPlugins.length; i < l; i++) {
                            var currentPlugin = enabledPlugins[i];
                            var brickManifest = JSON.parse(arguments[i]);

                            window.bricks[brickManifest.name] = brickManifest;

                            if(window.environment === 'debug') {
                                brickMainModule = 'canopsis/' + currentPlugin + '/' + 'init.dev.js';
                                brickManifest.envMode = 'development';
                            } else {
                                brickMainModule = 'canopsis/' + currentPlugin + '/' + 'init.js';
                            }

                            schemasInitFiles.push('canopsis/' + currentPlugin + '/' + 'init.schemas');
                            //remove the .js extension
                            brickMainModule = brickMainModule.slice(0, -3);

                            initFiles.push(brickMainModule);
                        }

                        require(schemasInitFiles, function() {
                            require(['canopsis/brick-loader/schemasloader'], function() {
                                require(initFiles, function() {
                                    var testDeps = [];
                                    if(window.environment === 'test') {
                                        for (var i = 0, l = enabledPlugins.length; i < l; i++) {
                                            var currentPlugin = enabledPlugins[i];
                                            testDeps.push('canopsis/' + currentPlugin + '/' + 'init.test');
                                        }
                                    }

                                    require(testDeps, function() {
                                        //This flag allow to prevent too early application requirement. @see "app/application" module
                                        window.appShouldNowBeLoaded = true;

                                        setLoadingInfo('Fetching application starting point', 'fa-plug');
                                        require(['canopsis/brick-loader/application'], function(Application) {
                                            setLoadingInfo('Initializing user interface', 'fa-desktop');

                                        });
                                    });
                                });
                            });
                        });
                    });
                });
            });
        });
    });
});
