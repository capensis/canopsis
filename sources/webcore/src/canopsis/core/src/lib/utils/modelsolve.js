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
    name: 'ModelsolveUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var modelsolve = Utility.create({

            name: 'modelsolve',

            gen_resolve: function(callback) {
                return function(data) {
                    for (var i = 0; i < data.data.length; i++) {
                        //data.data[i].id = data.data[i]._id;
                        //delete data.data[i]._id;
                    }

                    if(data.success === false && data.data.msg) {
                        console.error('Server Error', data.data.msg, data.data.type);
                        console.error(data.data.traceback);
                    }

                    Ember.run(null, callback, data);
                };
            },

            gen_reject: function(callback) {
                return function(xhr) {
                    xhr.then = null;
                    Ember.run(null, callback, xhr);
                };
            },
        });

        application.register('utility:modelsolve', modelsolve);
    }
});
