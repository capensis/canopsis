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

define([
    'canopsis/canopsisConfiguration'
], function(conf) {

    var i18n = {

        name: 'i18n',

        todo: [],
        translations: {},
        newTranslations: true,
        _: function(word, noDeprecation) {

            if(window.Ember && noDeprecation === undefined) {
                Ember.deprecate('You should not use i18n tools directly when ember is loaded. Please consider using Ember.String.loc instead. ', !conf.EmberIsLoaded);
            }

            if (typeof word !== 'string') {
                //This is not an interesting data type
                return word;
            } else if (!isNaN(parseInt(word))) {
                //This is just a number, it is useless to translate it.
                return word;
            } else {
                translated = i18n.translations[i18n.lang][word];

                if (translated) {
                    return i18n.showTranslation(translated);
                } else {
                    var isTranslated = true;
                    //adding translation to todo list
                    if (i18n.todo.indexOf(word) === -1) {

                        i18n.todo.push(word);
                        i18n.newTranslations = true;
                        isTranslated = false;

                    }
                    //returns original not translated string
                    return i18n.showTranslation(word, isTranslated);
                }
            }
        },
        showTranslation: function (word, isTranslated) {
            if (conf.SHOW_TRANSLATIONS) {
                if(isTranslated) {
                    circleColor = 'text-success';
                } else {
                    circleColor = 'text-danger';
                }
                return word + '<span class="fa-stack superscript"><i class="fa fa-circle fa-stack-2x ' + circleColor + '"></i><i class="fa fa-flag fa-stack-1x fa-inverse"></i></span>';
            } else {
                return word;
            }
        },

        uploadDefinitions: function () {}
    };

    window.__ = i18n._;
    window.i18n = i18n;

    i18n.lang = conf.getUserLanguage();

    if (conf.DEBUG && conf.TRANSLATE) {
        Ember.run(function(){
            setInterval(function () {
                if (i18n.newTranslations) {
                    console.log('Uploading new translations');
                    i18n.newTranslations = false;
                    i18n.uploadDefinitions();
                }

            }, 10000);
        });
    }

    Ember.Application.initializer({
        name:"I18nUtils",
        initialize: function(container, application) {
            application.register('utility:i18n', i18n);
        }
    });

    return i18n;
});
