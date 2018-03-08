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
 *
 * @module canopsis-frontend-core
 */


var deps = [];
if (!isIE && window.environment !== 'test') {
    deps.push('canopsis/core/externals/console.js/console');
}

define(deps, function() {
    conf = canopsisConfiguration;

    delete console.init;

    if(!isIE) {
        console.group('init');
        console.tags.add('init');
    } else {
        if(conf.DEBUG) {
            var baseconsole = console;
            console = {
                log: function (){
                    var fileinfo = '';
                    try {
                        var arrayargs = [];
                        for (var i=0; i<arguments.length; i++) {
                            arrayargs.push(arguments[i]);
                        }
                        var args = JSON.stringify(arrayargs);
                        baseconsole.log.apply(baseconsole, [fileinfo, args]);
                    } catch (e) {
                        Array.prototype.unshift.call(arguments, [' > unable to serialize next message']);
                        Array.prototype.unshift.call(arguments, fileinfo);
                    }
                },
                tags: {
                    add: function (){},
                    remove: function (){}
                },
                group: function () {},
                groupEnd: function () {},
                warn: function (){},
                debug: function (){},
                info: function (){},
                error: function (){},
            };
        }
    }

    // console.log = function(){};
    // console.info = function(){};
    // console.error = function(){};
    // console.group = function(){};
    // console.groupEnd = function(){};
    // console.warn = function(){};
    // Ember.warn = function(){};
    // Ember.deprecate = function(){};
    console.debug = console.log;

    return console;
});
